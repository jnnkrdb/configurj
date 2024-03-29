package v1alpha1

import (
	"context"

	"github.com/jnnkrdb/k8s/operator"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type GlobalConfigSpec struct {
	Immutable  bool   `json:"immutable"`
	Name       string `json:"name"`
	Namespaces struct {
		AvoidRegex []string `json:"avoidregex"`
		MatchRegex []string `json:"matchregex"`
	} `json:"namespaces"`
	Data map[string]string `json:"data"`
}

// deepcopy
func (in *GlobalConfig) DeepCopyInto(out *GlobalConfig) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = GlobalConfigSpec{
		Immutable:  in.Spec.Immutable,
		Name:       in.Spec.Name,
		Namespaces: in.Spec.Namespaces,
		Data:       in.Spec.Data,
	}
}

// ----------------------------------------------------
// kubernetes dependencies

type GlobalConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GlobalConfigSpec `json:"spec"`
}

type GlobalConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalConfig `json:"items"`
}

func (in *GlobalConfig) DeepCopyObject() runtime.Object {
	out := GlobalConfig{}
	in.DeepCopyInto(&out)
	return &out
}

func (in *GlobalConfigList) DeepCopyObject() runtime.Object {
	out := GlobalConfigList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]GlobalConfig, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return &out
}

// ----------------------------------------------------
// helper functions

const _GC_RESOURCE string = "globalconfigs"

// requests all deployed GlobalConfigs and returns them as a GlobalConfigList
func GetGlobalConfigList() (gcl GlobalConfigList, err error) {

	err = operator.CRD().Get().Resource(_GC_RESOURCE).Do(context.TODO()).Into(&gcl)

	return
}
