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
			_matchfunc := func() {

				for _, match := range matches {

					if regexMatch, err := regexp.MatchString(match, ns.Name); err != nil {

						prtcl.Log.Println("error while comparing [match] regexp with namespace-name:", err)

						prtcl.PrintObject(match, ns, err)

					} else {

						if regexMatch {

							_match = append(_match, ns.Name)
						}
					}
				}
			}

			if len(avoids) > 0 {

				for _, avoid := range avoids {

					if regexAvoid, err := regexp.MatchString(avoid, ns.Name); err != nil {

						prtcl.Log.Println("error while comparing [avoid] regexp with namespace-name:", err)

						prtcl.PrintObject(avoid, ns, err)

					} else {

						if !regexAvoid {

							_matchfunc()
						}
					}
				}

			} else {

				_matchfunc()
			}
		}
	}

	return
}
