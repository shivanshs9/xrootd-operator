// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Xrootd) DeepCopyInto(out *Xrootd) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Xrootd.
func (in *Xrootd) DeepCopy() *Xrootd {
	if in == nil {
		return nil
	}
	out := new(Xrootd)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Xrootd) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XrootdConfigSpec) DeepCopyInto(out *XrootdConfigSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XrootdConfigSpec.
func (in *XrootdConfigSpec) DeepCopy() *XrootdConfigSpec {
	if in == nil {
		return nil
	}
	out := new(XrootdConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XrootdList) DeepCopyInto(out *XrootdList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Xrootd, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XrootdList.
func (in *XrootdList) DeepCopy() *XrootdList {
	if in == nil {
		return nil
	}
	out := new(XrootdList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *XrootdList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XrootdRedirectorSpec) DeepCopyInto(out *XrootdRedirectorSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XrootdRedirectorSpec.
func (in *XrootdRedirectorSpec) DeepCopy() *XrootdRedirectorSpec {
	if in == nil {
		return nil
	}
	out := new(XrootdRedirectorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XrootdSpec) DeepCopyInto(out *XrootdSpec) {
	*out = *in
	out.Storage = in.Storage
	out.Worker = in.Worker
	out.Redirector = in.Redirector
	out.Config = in.Config
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XrootdSpec.
func (in *XrootdSpec) DeepCopy() *XrootdSpec {
	if in == nil {
		return nil
	}
	out := new(XrootdSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XrootdStatus) DeepCopyInto(out *XrootdStatus) {
	*out = *in
	in.EnforcingReconcileStatus.DeepCopyInto(&out.EnforcingReconcileStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XrootdStatus.
func (in *XrootdStatus) DeepCopy() *XrootdStatus {
	if in == nil {
		return nil
	}
	out := new(XrootdStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XrootdStorageSpec) DeepCopyInto(out *XrootdStorageSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XrootdStorageSpec.
func (in *XrootdStorageSpec) DeepCopy() *XrootdStorageSpec {
	if in == nil {
		return nil
	}
	out := new(XrootdStorageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *XrootdWorkerSpec) DeepCopyInto(out *XrootdWorkerSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new XrootdWorkerSpec.
func (in *XrootdWorkerSpec) DeepCopy() *XrootdWorkerSpec {
	if in == nil {
		return nil
	}
	out := new(XrootdWorkerSpec)
	in.DeepCopyInto(out)
	return out
}
