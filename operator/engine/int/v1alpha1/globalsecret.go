package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type GlobalSecretSpec struct {
	Immutable  bool              `json:"immutable"`
	Name       string            `json:"name"`
	Namespaces []string          `json:"namespaces"`
	Type       string            `json:"type"`
	Data       map[string]string `json:"data"`
}

// deepcopy
func (in *GlobalSecret) DeepCopyInto(out *GlobalSecret) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = GlobalSecretSpec{
		Immutable:  in.Spec.Immutable,
		Name:       in.Spec.Name,
		Namespaces: in.Spec.Namespaces,
		Type:       in.Spec.Type,
		Data:       in.Spec.Data,
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