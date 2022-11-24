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

// function to create, update and delete secrets from a globalsecret
func CRUD_Secrets(gs v1alpha1.GlobalSecret) {

	prtcl.Log.Println("process globalsecret:", gs.Namespace+"/"+gs.Name)

	prtcl.PrintObject(gs)

	// request the namespace lists (all namespaces, avoided via regex, matched via regex)
	all_namespaces, matched_namespaces := GetNamespaceLists(gs.Spec.Namespaces.AvoidRegex, gs.Spec.Namespaces.MatchRegex)

	// build create/delete functions
	_create := func(namespace string) {

		prtcl.Log.Println("creating secret:", namespace+"/"+gs.Spec.Name)

		var new = v1.Secret{}
		new.Name = gs.Spec.Name
		new.Namespace = namespace
		new.Annotations = func() map[string]string {
			result := make(map[string]string)
			result[ANNOTATION_RESOURCEVERSION] = gs.ResourceVersion
			return result
		}()
		new.Immutable = &gs.Spec.Immutable
		new.StringData = func() map[string]string {
			result := make(map[string]string)
			for k, v := range gs.Spec.Data {
				result[k] = fnc.UnencodeB64(v)
			}
			return result
		}()
		new.Type = v1.SecretType(gs.Spec.Type)

		if res, err := core.K8SCLIENT.CoreV1().Secrets(new.Namespace).Create(context.TODO(), &new, metav1.CreateOptions{}); err != nil {

			prtcl.Log.Println("error while creating secret [", new.Namespace+"/"+new.Name+"]:", err)

			prtcl.PrintObject(new, res, err)

		} else {

			prtcl.Log.Println("secret created:", res.Namespace+"/"+res.Name)
		}
	}

	_delete := func(namespace string) (err error) {

		err = nil

		prtcl.Log.Println("deleting secret:", namespace+"/"+gs.Spec.Name)

		if err = core.K8SCLIENT.CoreV1().Secrets(namespace).Delete(context.TODO(), gs.Spec.Name, metav1.DeleteOptions{}); err != nil {

			prtcl.Log.Println("error deleting secret ["+namespace+"/"+gs.Spec.Name+"]:", err)

		} else {

			prtcl.Log.Println("deleted secret:", namespace+"/"+gs.Spec.Name)
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

		if scrt, err := core.K8SCLIENT.CoreV1().Secrets(matchednamespace).Get(context.TODO(), gs.Spec.Name, metav1.GetOptions{}); err != nil {

			_create(matchednamespace)

		} else {

			// UPDATE
			if scrt.Annotations[ANNOTATION_RESOURCEVERSION] != gs.ResourceVersion {

				prtcl.Log.Println("updating secret:", scrt.Namespace+"/"+scrt.Name)

				// delete the old secret
				if _delete(matchednamespace) == nil {

					_create(matchednamespace)
				}
			}
		}
	}
}
