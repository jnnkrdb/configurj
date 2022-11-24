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
func CRUD_Secrets(gs v1alpha1.GlobalSecret) {

	prtcl.Log.Println("process globalsecret:", gs.Namespace+"/"+gs.Name)

	// request the namespace lists (all namespaces, avoided via regex, matched via regex)
	all_namespaces, matched_namespaces := GetNamespaceLists(gs.Spec.Namespaces.AvoidRegex, gs.Spec.Namespaces.MatchRegex)

	// DELETE
	for _, clusternamespace := range all_namespaces {

		if !fnc.StringInList(clusternamespace, matched_namespaces) {

			prtcl.Log.Println("deleting secret:", clusternamespace+"/"+gs.Spec.Name)

			if err := core.K8SCLIENT.CoreV1().Secrets(clusternamespace).Delete(context.TODO(), gs.Spec.Name, metav1.DeleteOptions{}); err != nil {

				prtcl.Log.Println("error deleting secret [", clusternamespace+"/"+gs.Spec.Name+"]:", err)

				prtcl.PrintObject(gs, clusternamespace, err)

			} else {

				prtcl.Log.Println("deleted secret:", clusternamespace+"/"+gs.Spec.Name)
			}
		}
	}

	for _, matchednamespace := range matched_namespaces {

		if scrt, err := core.K8SCLIENT.CoreV1().Secrets(matchednamespace).Get(context.TODO(), gs.Spec.Name, metav1.GetOptions{}); err != nil {

			// CREATE
			prtcl.Log.Println("creating secret:", matchednamespace+"/"+gs.Spec.Name)

			var new = v1.Secret{}
			new.Name = gs.Spec.Name
			new.Namespace = matchednamespace
			new.Annotations[ANNOTATION_RESOURCEVERSION] = gs.ResourceVersion
			new.Immutable = &gs.Spec.Immutable
			new.Data = make(map[string][]byte) // gs.Spec.Data

			if res, err := core.K8SCLIENT.CoreV1().Secrets(matchednamespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

				prtcl.Log.Println("error while creating secret [", new.Namespace+"/"+new.Name+"]:", err)

				prtcl.PrintObject(gs, matched_namespaces, matchednamespace, scrt, res, err)

			} else {

				prtcl.Log.Println("secret created:", res.Namespace+"/"+res.Name)
			}

		} else {

			// UPDATE
			if scrt.Annotations[ANNOTATION_RESOURCEVERSION] != gs.ResourceVersion {

				prtcl.Log.Println("updating secret:", scrt.Namespace+"/"+scrt.Name)

				// delete the old secret
				if err := core.K8SCLIENT.CoreV1().Secrets(scrt.Namespace).Delete(context.TODO(), scrt.Name, metav1.DeleteOptions{}); err != nil {

					prtcl.Log.Println("error updating secret [", scrt.Namespace+"/"+scrt.Name, "]:", err)

					prtcl.PrintObject(gs, matched_namespaces, matchednamespace, scrt, err)

				} else {

					time.Sleep(1 * time.Minute)

					var new = v1.Secret{}
					new.Name = gs.Spec.Name
					new.Namespace = matchednamespace
					new.Annotations[ANNOTATION_RESOURCEVERSION] = gs.ResourceVersion
					new.Immutable = &gs.Spec.Immutable
					new.Data = func(resource map[string]string) map[string][]byte {
						result := make(map[string][]byte)
						for k, v := range resource {
							result[k] = []byte(v)
						}
						return result
					}(gs.Spec.Data)

					if res, err := core.K8SCLIENT.CoreV1().Secrets(new.Namespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

						prtcl.Log.Println("error updating secret [", new.Namespace+"/"+new.Name, "]:", err)

						prtcl.PrintObject(gs, matched_namespaces, matchednamespace, scrt, res, err)

					} else {

						prtcl.Log.Println("updated secret:", new.Namespace+"/"+new.Name)
					}
				}
			}
		}
	}
}
