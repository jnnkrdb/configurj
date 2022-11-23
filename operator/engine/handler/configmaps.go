package handler

import (
	"context"
	"time"

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

	// request the namespace lists (all namespaces, avoided via regex, matched via regex)
	all_namespaces, matched_namespaces := GetNamespaceLists(gc.Spec.Namespaces.AvoidRegex, gc.Spec.Namespaces.MatchRegex)

	// DELETE
	for _, clusternamespace := range all_namespaces {

		if !fnc.StringInList(clusternamespace, matched_namespaces) {

			prtcl.Log.Println("deleting configmap:", clusternamespace+"/"+gc.Spec.Name)

			if err := core.K8SCLIENT.CoreV1().ConfigMaps(clusternamespace).Delete(context.TODO(), gc.Spec.Name, metav1.DeleteOptions{}); err != nil {

				prtcl.Log.Println("error deleting configmap [", clusternamespace+"/"+gc.Spec.Name+"]:", err)

				prtcl.PrintObject(gc, clusternamespace, err)

			} else {

				prtcl.Log.Println("deleted configmap:", clusternamespace+"/"+gc.Spec.Name)
			}
		}
	}

	for _, matchednamespace := range matched_namespaces {

		if cm, err := core.K8SCLIENT.CoreV1().ConfigMaps(matchednamespace).Get(context.TODO(), gc.Spec.Name, metav1.GetOptions{}); err != nil {

			// CREATE
			prtcl.Log.Println("creating configmap", matchednamespace+"/"+gc.Spec.Name)

			var new = v1.ConfigMap{}
			new.Name = gc.Spec.Name
			new.Namespace = matchednamespace
			new.Annotations[ANNOTATION_RESOURCEVERSION] = gc.ResourceVersion
			new.Immutable = &gc.Spec.Immutable
			new.Data = gc.Spec.Data

			if res, err := core.K8SCLIENT.CoreV1().ConfigMaps(matchednamespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

				prtcl.Log.Println("error while creating configmap", new.Namespace+"/"+new.Name+":", err)

				prtcl.PrintObject(gc, matched_namespaces, matchednamespace, cm, res, err)

			} else {

				prtcl.Log.Println("configmap created:", res.Namespace+"/"+res.Name)
			}

		} else {

			// UPDATE
			if cm.Annotations[ANNOTATION_RESOURCEVERSION] != gc.ResourceVersion {

				prtcl.Log.Println("updating configmap", cm.Namespace+"/"+cm.Name)

				// delete the old configmap
				if err := core.K8SCLIENT.CoreV1().ConfigMaps(cm.Namespace).Delete(context.TODO(), cm.Name, metav1.DeleteOptions{}); err != nil {

					prtcl.Log.Println("error updating configmap [", cm.Namespace+"/"+cm.Name, "]:", err)

					prtcl.PrintObject(gc, matched_namespaces, matchednamespace, cm, err)

				} else {

					time.Sleep(1 * time.Minute)

					var new = v1.ConfigMap{}
					new.Name = gc.Spec.Name
					new.Namespace = matchednamespace
					new.Annotations[ANNOTATION_RESOURCEVERSION] = gc.ResourceVersion
					new.Immutable = &gc.Spec.Immutable
					new.Data = gc.Spec.Data

					if res, err := core.K8SCLIENT.CoreV1().ConfigMaps(new.Namespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

						prtcl.Log.Println("error updating configmap [", cm.Namespace+"/"+cm.Name, "]:", err)

						prtcl.PrintObject(gc, matched_namespaces, matchednamespace, cm, res, err)

					} else {

						prtcl.Log.Println("configmap updated:", new.Namespace+"/"+new.Name)
					}
				}
			}
		}
	}
}
