package handler

import (
	"context"

	"github.com/jnnkrdb/configurj-engine/core"
	"github.com/jnnkrdb/corerdb/fnc"
	"github.com/jnnkrdb/corerdb/prtcl"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func get_AllNamespaces() (*v1.NamespaceList, error) {

	// request all namespaces of the cluster
	if all_ns, err := core.K8SCLIENT.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{}); err != nil {

		prtcl.Log.Println("error while requesting the namespaces:", err)

		prtcl.PrintObject(all_ns, err)

		return nil, err

	} else {

		prtcl.Log.Println("successfully received all namespaces")

		return all_ns, nil
	}
}

// receive all namespaces and calculate the allowed ones
func calculateAllowedNamespaces(avoids []string) []string {

	prtcl.Log.Println("calculating the allowed namespaces")

	prtcl.PrintObject(avoids)

	allowedNamespaces := []string{}

	// request all namespaces of the cluster
	if all_ns, err := get_AllNamespaces(); err == nil {

		// calculate the allowed namespaces with the given avoids array
		for _, ns := range all_ns.Items {

			if fnc.StringInList(ns.Name, avoids) {

				allowedNamespaces = append(allowedNamespaces, ns.Name)
			}
		}
	}

	// set the allowed namespaces
	return allowedNamespaces
}
