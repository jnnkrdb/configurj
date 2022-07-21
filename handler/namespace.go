package handler

import (
	"context"
	"log"
	"strings"
	"time"

	"configurj/probes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// caching namespaces, grouped by entity (all, configmap, secret)
	__NS_ALL = &[]string{}
)

// returns the requested list of the namespace
func GetAllNamespaces() []string {
	return *__NS_ALL
}

// compares the given string with the given list, if the list contains the
// string, the returnvalue will be true
func ListContainsString(_list []string, _string string) bool {

	for _, ns := range _list {

		if ns != "" {

			if ns == _string {

				return true
			}
		}
	}

	return false
}

// get all required namespaces without the avoided ones
func AsyncWatchNamespaces(_log *log.Logger, _k8sclient *kubernetes.Clientset) {

	awaittimer := 60

	for {

		ns_all := &[]string{}

		if allnamespaces, err := _k8sclient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{}); err != nil {

			_log.Printf("%s | %s\n", "ERROR", err.Error())

			awaittimer = 2

			probes.LIVENESS = 500

		} else {

			for _, namespace := range allnamespaces.Items {

				*ns_all = append(*ns_all, namespace.Name)
			}

			awaittimer = 60

			probes.LIVENESS = 200
		}

		__NS_ALL = ns_all

		time.Sleep(time.Duration(int(awaittimer)) * time.Second)
	}
}

// get the namespaces where the item should be deployed to
func GetDistributeNamespaces(_globalavoidnamespaces []string, _annotations map[string]string) []string {

	currentAllNamespaces := GetAllNamespaces()

	// ###########################################################################################################################
	// get all avoided namespaces for the current secret
	// if the annotation <configurj.jnnkrdb.de/avoid> is empty, it will
	// be ignored, if the value is '*', all namespace will be avoided (just for completion, has no purpose)
	totalavoidednamespaces := _globalavoidnamespaces

	if _annotations[ANNOTATION_AVOID] != "" {

		if _annotations[ANNOTATION_AVOID] == "*" {

			totalavoidednamespaces = currentAllNamespaces

		} else {

			for _, namespace := range strings.Split(_annotations[ANNOTATION_AVOID], ";") {

				if !ListContainsString(_globalavoidnamespaces, namespace) {

					totalavoidednamespaces = append(totalavoidednamespaces, namespace)
				}
			}
		}
	}

	// ###########################################################################################################################
	// get all desired namespaces from the secret annotation
	// if the annotation <configurj.jnnkrdb.de/match> is empty, it will
	// be ignored, if the value is '*', all namespace will be matched
	totaldesirednamespaces := []string{}

	if _annotations[ANNOTATION_MATCH] != "" {

		if _annotations[ANNOTATION_MATCH] == "*" {

			totaldesirednamespaces = currentAllNamespaces

		} else {

			for _, namespace := range strings.Split(_annotations[ANNOTATION_MATCH], ";") {

				if ListContainsString(currentAllNamespaces, namespace) {

					totaldesirednamespaces = append(totaldesirednamespaces, namespace)
				}
			}
		}
	}

	// ###########################################################################################################################
	// compare all desired namespaces with the total avoided namespaces
	// to get all possible namespaces the secret should be distributed to
	distributenamespaces := []string{}

	for _, namespace := range totaldesirednamespaces {

		if !ListContainsString(totalavoidednamespaces, namespace) {

			distributenamespaces = append(distributenamespaces, namespace)
		}
	}

	return distributenamespaces
}
