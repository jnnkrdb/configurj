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

// function to create, update and delete secrets from a globalsecret
func crud_Secrets(allowedNamespaces []string, gs v1alpha1.GlobalSecret) {

	prtcl.Log.Println("process globalsecret:", gs.Namespace+"/"+gs.Name)

	// DELETE
	if allnamespaces, err := get_AllNamespaces(); err == nil {

		for _, clusternamespace := range allnamespaces.Items {

			if !fnc.StringInList(clusternamespace.Name, allowedNamespaces) || !fnc.StringInList(clusternamespace.Name, gs.Spec.Namespaces) {

				prtcl.Log.Println("deleting secret", gs.Spec.Name, "from namespace", clusternamespace.Name)

				if err := core.K8SCLIENT.CoreV1().Secrets(clusternamespace.Name).Delete(context.TODO(), gs.Spec.Name, metav1.DeleteOptions{}); err != nil {

					prtcl.Log.Println("error deleting secret", clusternamespace.Name+"/"+gs.Spec.Name+":", err)

					prtcl.PrintObject(gs, allnamespaces, clusternamespace, err)

				} else {

					prtcl.Log.Println("deleted secret", clusternamespace.Name+"/"+gs.Spec.Name)
				}
			}
		}
	}

	for _, gs_spec_namespace := range gs.Spec.Namespaces {

		if cm, err := core.K8SCLIENT.CoreV1().Secrets(gs_spec_namespace).Get(context.TODO(), gs.Spec.Name, metav1.GetOptions{}); err != nil {

			// CREATE
			if fnc.StringInList(gs_spec_namespace, allowedNamespaces) {

				prtcl.Log.Println("creating secret", gs_spec_namespace+"/"+gs.Spec.Name)

				var new = v1.Secret{}
				new.Name = gs.Spec.Name
				new.Namespace = gs_spec_namespace
				new.Annotations["configurj.jnnkrdb.de/version"] = gs.ResourceVersion
				new.Immutable = &gs.Spec.Immutable
				new.Data = func(resource map[string]string) map[string][]byte {
					result := make(map[string][]byte)
					for k, v := range resource {
						result[k] = []byte(v)
					}
					return result
				}(gs.Spec.Data)

				if res, err := core.K8SCLIENT.CoreV1().Secrets(gs_spec_namespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

					prtcl.Log.Println("error while creating secret", new.Namespace+"/"+new.Name+":", err)

					prtcl.PrintObject(gs, allowedNamespaces, gs_spec_namespace, cm, res, err)

				} else {

					prtcl.Log.Println("secret created:", res.Namespace+"/"+res.Name)
				}
			}

		} else {

			// UPDATE
			if cm.Annotations["configurj.jnnkrdb.de/version"] != gs.ResourceVersion {

				prtcl.Log.Println("updating secret", cm.Namespace+"/"+cm.Name)

				// delete the old secret
				if err := core.K8SCLIENT.CoreV1().Secrets(cm.Namespace).Delete(context.TODO(), cm.Name, metav1.DeleteOptions{}); err != nil {

					prtcl.Log.Println("updating secret", cm.Namespace+"/"+cm.Name, "failed:", err)

					prtcl.PrintObject(gs, allowedNamespaces, gs_spec_namespace, cm, err)

				} else {

					time.Sleep(2 * time.Minute)

					var new = v1.Secret{}
					new.Name = gs.Spec.Name
					new.Namespace = gs_spec_namespace
					new.Annotations["configurj.jnnkrdb.de/version"] = gs.ResourceVersion
					new.Immutable = &gs.Spec.Immutable
					new.Data = func(resource map[string]string) map[string][]byte {
						result := make(map[string][]byte)
						for k, v := range resource {
							result[k] = []byte(v)
						}
						return result
					}(gs.Spec.Data)

					if res, err := core.K8SCLIENT.CoreV1().Secrets(new.Namespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

						prtcl.Log.Println("updating secret", new.Namespace+"/"+new.Name, "failed:", err)

						prtcl.PrintObject(gs, allowedNamespaces, gs_spec_namespace, cm, res, err)

					} else {

						prtcl.Log.Println("secret", new.Namespace+"/"+new.Name, "updated")
					}
				}
			}
		}
	}
}
