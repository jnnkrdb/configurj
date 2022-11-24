package handler

import (
	"context"

	"github.com/jnnkrdb/configurj-engine/core"
	"github.com/jnnkrdb/configurj-engine/int/v1alpha1"

	"github.com/jnnkrdb/corerdb/fnc"
	"github.com/jnnkrdb/corerdb/prtcl"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// function to create, update and delete configmaps from a globalconfig
func CRUD_Configmaps(gc v1alpha1.GlobalConfig) {

	prtcl.Log.Println("process globalconfig:", gc.Namespace+"/"+gc.Name)

	prtcl.PrintObject(gc)

	// request the namespace lists (all namespaces, avoided via regex, matched via regex)
	all_namespaces, matched_namespaces := GetNamespaceLists(gc.Spec.Namespaces.AvoidRegex, gc.Spec.Namespaces.MatchRegex)

	// build create/delete functions
	_create := func(namespace string) {

		prtcl.Log.Println("creating configmap", namespace+"/"+gc.Spec.Name)

		var new = v1.ConfigMap{}
		new.Name = gc.Spec.Name
		new.Namespace = namespace
		new.Annotations = func() map[string]string {
			result := make(map[string]string)
			result[ANNOTATION_RESOURCEVERSION] = gc.ResourceVersion
			return result
		}()
		new.Immutable = &gc.Spec.Immutable
		new.Data = gc.Spec.Data

		if res, err := core.K8SCLIENT.CoreV1().ConfigMaps(new.Namespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

			prtcl.Log.Println("error while creating configmap", new.Namespace+"/"+new.Name+":", err)

			prtcl.PrintObject(new, res, err)

		} else {

			prtcl.Log.Println("configmap created:", res.Namespace+"/"+res.Name)

			prtcl.PrintObject(new)
		}
	}

	_delete := func(namespace string) (err error) {

		err = nil

		prtcl.Log.Println("deleting configmap:", namespace+"/"+gc.Spec.Name)

		if err := core.K8SCLIENT.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), gc.Spec.Name, metav1.DeleteOptions{}); err != nil {

			prtcl.Log.Println("error deleting configmap [", namespace+"/"+gc.Spec.Name+"]:", err)

		} else {

			prtcl.Log.Println("deleted configmap:", namespace+"/"+gc.Spec.Name)
		}

		return
	}

	// DELETE
	for _, clusternamespace := range all_namespaces {

		if !fnc.StringInList(clusternamespace, matched_namespaces) {

			_delete(clusternamespace)
		}
	}

	for _, matchednamespace := range matched_namespaces {

		if cm, err := core.K8SCLIENT.CoreV1().ConfigMaps(matchednamespace).Get(context.TODO(), gc.Spec.Name, metav1.GetOptions{}); err != nil {

			_create(matchednamespace)

		} else {

			// UPDATE
			if cm.Annotations[ANNOTATION_RESOURCEVERSION] != gc.ResourceVersion {

				prtcl.Log.Println("updating configmap", cm.Namespace+"/"+cm.Name)

				// delete the old configmap
				if _delete(matchednamespace) == nil {

					_create(matchednamespace)
				}
			}
		}
	}
}
