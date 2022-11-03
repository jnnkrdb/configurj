package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type GlobalSecretSpec struct {
	Namespaces []string `json:"namespaces"`
}

// deepcopy
func (in *GlobalSecret) DeepCopyInto(out *GlobalSecret) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = GlobalSecretSpec{
		Namespaces: in.Spec.Namespaces,
	}
}

// ----------------------------------------------------
// kubernetes dependencies
type GlobalSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GlobalSecretSpec `json:"spec"`
}

type GlobalSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalSecret `json:"items"`
}

func (in *GlobalSecret) DeepCopyObject() runtime.Object {
	out := GlobalSecret{}
	in.DeepCopyInto(&out)
	return &out
}

func (in *GlobalSecretList) DeepCopyObject() runtime.Object {
	out := GlobalSecretList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]GlobalSecret, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}
