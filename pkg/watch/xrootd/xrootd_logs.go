package xrootd

import (
	"context"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/msoap/byline"
	"github.com/pkg/errors"
	xrootdv1alpha1 "github.com/xrootd/xrootd-k8s-operator/pkg/apis/xrootd/v1alpha1"
	"github.com/xrootd/xrootd-k8s-operator/pkg/controller/reconciler"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/constant"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/k8sutil"
	"github.com/xrootd/xrootd-k8s-operator/pkg/utils/types"
	"github.com/xrootd/xrootd-k8s-operator/pkg/watch"
	corev1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type LogsWatcher struct {
	Component  types.ComponentName
	reconciler reconciler.Reconciler
}

var _ watch.Watcher = LogsWatcher{}

var log = logf.Log.WithName("XrootdLogsWatcher")

const waitMemberReadyDelay = 5 * time.Second

func (lw LogsWatcher) Watch(requests <-chan reconcile.Request) error {
	var reqLogger logr.Logger
	for request := range requests {
		reqLogger = log.WithValues("request", request, "component", lw.Component)

		instance := &xrootdv1alpha1.Xrootd{}
		if err := lw.reconciler.GetResourceInstance(request, instance); err != nil {
			return err
		}
		if reconciler.IsBeingDeleted(instance) {
			// Skip processing if requested instance is being deleted
			reqLogger.Info("Xrootd instance is being deleted...", "request", request)
			continue
		}
		if err := lw.monitorXrootdStatus(request); err != nil {
			reqLogger.Error(err, "Failed to monitor xrootd cluster...")
		}
	}
	return nil
}

func (lw LogsWatcher) monitorXrootdStatus(request reconcile.Request) error {
	reqLogger := log.WithValues("request", request, "component", lw.Component)
	reqLogger.Info("Started monitoring xrootd cluster...")

	clientset, err := kubernetes.NewForConfig(lw.reconciler.GetConfig())
	if err != nil {
		return errors.Wrap(err, "unable to get kubernetes clientset")
	}
	for {
		time.Sleep(waitMemberReadyDelay)
		instance := &xrootdv1alpha1.Xrootd{}
		if err := lw.reconciler.GetResourceInstance(request, instance); err != nil {
			return errors.Wrap(err, "failed to refresh xrootd instance")
		}
		var unreadyPods []string
		if lw.Component == constant.XrootdRedirector {
			unreadyPods = instance.Status.RedirectorStatus.Unready
		} else {
			unreadyPods = instance.Status.WorkerStatus.Unready
		}
		if len(unreadyPods) == 0 {
			break
		}
		if err := lw.updateInstanceStatus(instance, lw.obtainLogsOfAllPods(request, unreadyPods, clientset)); err != nil {
			return errors.Wrap(err, "failed updating xrootd status")
		}
	}
	return nil
}

func (lw LogsWatcher) updateInstanceStatus(instance *xrootdv1alpha1.Xrootd, resultChannel <-chan podStatus) error {
	unreadyPods := make([]string, 0)
	readyPods := make([]string, 0)
	for resultStatus := range resultChannel {
		if resultStatus.isReady {
			readyPods = append(readyPods, resultStatus.podName)
		} else {
			unreadyPods = append(unreadyPods, resultStatus.podName)
		}
	}
	status := utils.NewMemberStatus(readyPods, unreadyPods)
	if lw.Component == constant.XrootdWorker {
		instance.Status.WorkerStatus = status
	} else if lw.Component == constant.XrootdRedirector {
		instance.Status.RedirectorStatus = status
	}
	if err := lw.reconciler.GetClient().Status().Update(context.TODO(), instance); err != nil {
		return err
	}
	return nil
}

func (lw LogsWatcher) obtainLogsOfAllPods(request reconcile.Request, unreadyPods []string, clientset *kubernetes.Clientset) <-chan podStatus {
	totalPods := len(unreadyPods)
	opt := &corev1.PodLogOptions{
		Follow:    true,
		Container: string(constant.Cmsd),
	}
	resultChannel := make(chan podStatus, totalPods)
	for _, podName := range unreadyPods {
		pod := &corev1.Pod{}
		podNamespacedName := k8stypes.NamespacedName{
			Namespace: request.Namespace,
			Name:      podName,
		}
		if err := lw.reconciler.GetClient().Get(context.TODO(), podNamespacedName, pod); err != nil {
			resultChannel <- podStatus{
				podName: podName,
				isReady: false,
			}
		} else {
			go lw.processXrootdPodLogs(pod, opt, clientset, resultChannel)
		}
	}
	return resultChannel
}

func (lw LogsWatcher) processXrootdPodLogs(pod *corev1.Pod, opt *corev1.PodLogOptions, clientset *kubernetes.Clientset, resultChannel chan<- podStatus) {
	reqLogger := log.WithValues("pod", pod.Name, "component", lw.Component)

	unreadyStatus := podStatus{
		podName: pod.Name,
		isReady: false,
	}

	var err error
	var reader io.ReadCloser
	for {
		reader, err = k8sutil.GetPodLogStream(*pod, opt, clientset)
		if err != nil {
			if strings.Contains(err.Error(), "ContainerCreating") {
				reqLogger.V(1).Info("Container not started yet, retrying...", "error", err)
			} else {
				reqLogger.Error(err, "unable to get pod stream", "options", opt)
				resultChannel <- unreadyStatus
				return
			}
		} else {
			break
		}
	}
	defer reader.Close()

	lineReader := byline.NewReader(reader)

	var regex *regexp.Regexp
	if lw.Component == constant.XrootdRedirector {
		regex = regexp.MustCompile(`Protocol: redirector..+ logged in.$`)
	} else if lw.Component == constant.XrootdWorker {
		regex = regexp.MustCompile(`Protocol: Logged into .+$`)
	}

	reqLogger.Info("Grepping and reading...", "regex", regex)
	buffer := make([]byte, 50)
	read, err := lineReader.GrepByRegexp(regex).Read(buffer)
	reqLogger.V(1).Info("Read to buffer", "length", read, "buffer", buffer)

	result := read > 0

	status := corev1.ConditionFalse
	if result {
		status = corev1.ConditionTrue
	}
	pod.Status.Conditions = append(pod.Status.Conditions, corev1.PodCondition{
		Type:   constant.XrootdPodConnection,
		Status: status,
		Reason: "Cmsd logs confirmed logged-in status",
	})
	if err = lw.reconciler.GetClient().Status().Update(context.TODO(), pod); err != nil {
		reqLogger.Error(err, "failed updating pod status", "status", pod.Status)
		resultChannel <- unreadyStatus
	}

	unreadyStatus.isReady = result
	resultChannel <- unreadyStatus
}

func NewLogsWatcher(component types.ComponentName, reconciler reconciler.Reconciler) watch.Watcher {
	return watch.NewGroupedRequestWatcher(
		LogsWatcher{
			Component:  component,
			reconciler: reconciler,
		},
	)
}

type podStatus struct {
	podName string
	isReady bool
}
