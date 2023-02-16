package handler

import (
	"context"

	"github.com/jnnkrdb/configurj-engine/env"
	"github.com/jnnkrdb/k8s/operator"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// receive the namespace lists (all namespaces, avoided via regex, matched via regex)
func GetNamespaceLists(avoids, matches []string) ([]string, []string, error) {

	nslog := env.Log().WithFields(logrus.Fields{
		"namespaces.avoid": avoids,
		"namesapces.match": matches,
	})

	nslog.Debug("requesting namespaces")

	nslog.Trace("allocating storage")

	var (
		_all   []string
		_match []string

		all_ns *v1.NamespaceList
		err    error
	)

	nslog.Trace("requesting namespaces from cluster, parsing into a list")

	if all_ns, err = operator.K8S().CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{}); err != nil {

		nslog.WithError(err).Error("error while requesting the namespaces")

	} else {

		nslog.WithField("namespaces.all", all_ns.Items).Trace("received list of namespaces")

		nslog.Debug("received list of namespaces")

		for id, ns := range all_ns.Items {

			currnslog := nslog.WithFields(logrus.Fields{
				"list.id":                id,
				"current.namespace.name": ns.Name,
				"current.namespace.uid":  ns.UID,
			})

			currnslog.Trace("checking namespace")

			// add to all-list
			_all = append(_all, ns.Name)

			currnslog.WithField("namespaces._all", _all).Trace("added namespace to collection")

			if len(avoids) > 0 {

				currnslog.Trace("checking avoided namespaces")

				if res, err := FindStringInRegexpList(ns.Name, avoids); res && err == nil {

					currnslog.Trace("found namespace in avoids list -> continue")

					continue
				}

				currnslog.Trace("no namespaces to avoid")
			}

			currnslog.Trace("checking namespaces to match")

			if res, err := FindStringInRegexpList(ns.Name, matches); res && err == nil {

				_match = append(_match, ns.Name)

				currnslog.WithField("namespaces._match", _match).Trace("added namespace to matches")
			}
		}
	}

	nslog.WithFields(logrus.Fields{
		"namespaces._all":   _all,
		"namespaces._match": _match,
		"error":             err,
	}).Debug("finished namespaces calculation")

	return _all, _match, err
}
