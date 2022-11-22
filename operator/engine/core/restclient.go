package core

import (
	"github.com/jnnkrdb/configurj-engine/int/v1alpha1"

	"github.com/jnnkrdb/corerdb/prtcl"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// kubernetes client for defaut kubernetes objects
var K8SCLIENT *kubernetes.Clientset

// kubernetes rest client for crds
var CRDCLIENT *rest.RESTClient

// load the default rest client for the operator
func LoadRestClient() error {

	prtcl.Log.Println("initializing kubernetes api-connection")

	if config, err := rest.InClusterConfig(); err != nil {

		prtcl.Log.Println("error while initialization:", err)

		prtcl.PrintObject(config, err)

		return err

	} else {

		if cs, err := kubernetes.NewForConfig(config); err != nil {

			prtcl.Log.Println("clientset error:", err)

			prtcl.PrintObject(config, cs, err)

			return err

		} else {

			K8SCLIENT = cs
		}

		v1alpha1.AddToScheme(scheme.Scheme)

		// create the rest client for the crds
		crdConf := *config

		crdConf.ContentConfig.GroupVersion = &schema.GroupVersion{
			Group:   v1alpha1.GroupName,
			Version: v1alpha1.GroupVersion,
		}

		crdConf.APIPath = "/apis"

		crdConf.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)

		crdConf.UserAgent = rest.DefaultKubernetesUserAgent()

		if crdRestClient, err := rest.UnversionedRESTClientFor(&crdConf); err != nil {

			prtcl.Log.Println("failed loading restclient:", err)

			prtcl.PrintObject(config, crdConf, crdRestClient, err)

			return err

		} else {

			CRDCLIENT = crdRestClient

			return nil
		}
	}
}
