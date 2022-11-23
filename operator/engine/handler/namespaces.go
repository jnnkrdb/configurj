package handler

import (
	"context"
	"regexp"

	"github.com/jnnkrdb/configurj-engine/core"
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
			for _, a := range avoids {

				if regexAvoidMatched, err := regexp.MatchString(a, ns.Name); err != nil {

					prtcl.Log.Println("error while comparing [avoid] regexp with namespace-name:", err)

					prtcl.PrintObject(a, ns, err)

				} else {

					if !regexAvoidMatched {

						for _, m := range matches {

							if regexMatchMatched, err := regexp.MatchString(m, ns.Name); err != nil {

								prtcl.Log.Println("error while comparing [match] regexp with namespace-name:", err)

								prtcl.PrintObject(m, ns, err)

							} else {

								if regexMatchMatched {

									_match = append(_match, ns.Name)
								}
							}
						}
					}
				}
			}
		}
	}

	return
}
