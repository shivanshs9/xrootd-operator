module github.com/shivanshs9/xrootd-operator

go 1.13

require (
	github.com/RHsyseng/operator-utils v0.0.0-20200619180557-7c49e58877d7
	github.com/operator-framework/operator-sdk v0.17.0
	github.com/redhat-cop/operator-utils v0.2.4
	github.com/shivanshs9/ty v1.1.0
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.17.7
	k8s.io/apimachinery v0.17.7
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.5.6
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.17.7 // Required by prometheus-operator
)
