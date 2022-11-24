package handler

import (
	"context"

	"github.com/jnnkrdb/configurj-engine/core"
	"github.com/jnnkrdb/corerdb/fnc"
	"github.com/jnnkrdb/corerdb/prtcl"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// receive the namespace lists (all namespaces, avoided via regex, matched via regex)
func GetNamespaceLists(avoids, matches []string) (_all, _match []string) {

	_all, _match = []string{}, []string{}

	if all_ns, err := core.K8SCLIENT.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{}); err != nil {

		prtcl.Log.Println("error while requesting the namespaces:", err)

		prtcl.PrintObject(all_ns, err)

	} else {

		for _, ns := range all_ns.Items {

			// add to all-list
			_all = append(_all, ns.Name)

			// calculate match list
			if len(avoids) > 0 {

				if !fnc.FindStringInRegexpList(ns.Name, avoids) {

					if fnc.FindStringInRegexpList(ns.Name, matches) {

						_match = append(_match, ns.Name)
					}
				}

			} else {

				if fnc.FindStringInRegexpList(ns.Name, matches) {

					_match = append(_match, ns.Name)
				}
			}
		}
	}

	prtcl.PrintObject(_all, _match)

	return
}
