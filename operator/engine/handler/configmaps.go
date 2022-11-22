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
func crud_Configmaps(allowedNamespaces []string, gc v1alpha1.GlobalConfig) {

	prtcl.Log.Println("process globalconfig:", gc.Namespace+"/"+gc.Name)

	// DELETE
	if allnamespaces, err := get_AllNamespaces(); err == nil {

		for _, clusternamespace := range allnamespaces.Items {

			if !fnc.StringInList(clusternamespace.Name, allowedNamespaces) || !fnc.StringInList(clusternamespace.Name, gc.Spec.Namespaces) {

				prtcl.Log.Println("deleting configmap", gc.Spec.Name, "from namespace", clusternamespace.Name)

				if err := core.K8SCLIENT.CoreV1().ConfigMaps(clusternamespace.Name).Delete(context.TODO(), gc.Spec.Name, metav1.DeleteOptions{}); err != nil {

					prtcl.Log.Println("error deleting configmap", clusternamespace.Name+"/"+gc.Spec.Name+":", err)

					prtcl.PrintObject(gc, allnamespaces, clusternamespace, err)

				} else {

					prtcl.Log.Println("deleted configmap", clusternamespace.Name+"/"+gc.Spec.Name)
				}
			}
		}
	}

	for _, gc_spec_namespace := range gc.Spec.Namespaces {

		if cm, err := core.K8SCLIENT.CoreV1().ConfigMaps(gc_spec_namespace).Get(context.TODO(), gc.Spec.Name, metav1.GetOptions{}); err != nil {

			// CREATE
			if fnc.StringInList(gc_spec_namespace, allowedNamespaces) {

				prtcl.Log.Println("creating configmap", gc_spec_namespace+"/"+gc.Spec.Name)

				var new = v1.ConfigMap{}
				new.Name = gc.Spec.Name
				new.Namespace = gc_spec_namespace
				new.Annotations["configurj.jnnkrdb.de/version"] = gc.ResourceVersion
				new.Immutable = &gc.Spec.Immutable
				new.Data = gc.Spec.Data

				if res, err := core.K8SCLIENT.CoreV1().ConfigMaps(gc_spec_namespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

					prtcl.Log.Println("error while creating configmap", new.Namespace+"/"+new.Name+":", err)

					prtcl.PrintObject(gc, allowedNamespaces, gc_spec_namespace, cm, res, err)

				} else {

					prtcl.Log.Println("configmap created:", res.Namespace+"/"+res.Name)
				}
			}

		} else {

			// UPDATE
			if cm.Annotations["configurj.jnnkrdb.de/version"] != gc.ResourceVersion {

				prtcl.Log.Println("updating configmap", cm.Namespace+"/"+cm.Name)

				// delete the old configmap
				if err := core.K8SCLIENT.CoreV1().ConfigMaps(cm.Namespace).Delete(context.TODO(), cm.Name, metav1.DeleteOptions{}); err != nil {

					prtcl.Log.Println("updating configmap", cm.Namespace+"/"+cm.Name, "failed:", err)

					prtcl.PrintObject(gc, allowedNamespaces, gc_spec_namespace, cm, err)

				} else {

					time.Sleep(2 * time.Minute)

					var new = v1.ConfigMap{}
					new.Name = gc.Spec.Name
					new.Namespace = gc_spec_namespace
					new.Annotations["configurj.jnnkrdb.de/version"] = gc.ResourceVersion
					new.Immutable = &gc.Spec.Immutable
					new.Data = gc.Spec.Data

					if res, err := core.K8SCLIENT.CoreV1().ConfigMaps(new.Namespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

						prtcl.Log.Println("updating configmap", new.Namespace+"/"+new.Name, "failed:", err)

						prtcl.PrintObject(gc, allowedNamespaces, gc_spec_namespace, cm, res, err)

					} else {

						prtcl.Log.Println("configmap", new.Namespace+"/"+new.Name, "updated")
					}
				}
			}
		}
	}
}
