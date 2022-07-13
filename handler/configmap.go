package handler

import (
	"context"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// creates the configmap from a template (_configmap) in the namespace (_namespace)
func createConfigMap(_configmap v1.ConfigMap, _namespace string) error {

	var temp = v1.ConfigMap{}
	temp.Name = _configmap.Name
	temp.Namespace = _namespace
	temp.Annotations = GetAnnotations(_configmap.Namespace, _configmap.Name, _configmap.ResourceVersion)
	temp.Labels = _configmap.Labels
	temp.Labels[LABELS_K8S_INSTANCE] = _namespace + "--" + temp.Name + "--" + _configmap.ResourceVersion
	temp.Data = _configmap.Data
	temp.Immutable = &__IMMUTABLE
	temp.BinaryData = _configmap.BinaryData
	//

	if res, err := __K8SCLIENT.CoreV1().ConfigMaps(_namespace).Create(context.TODO(), &temp, metav1.CreateOptions{}); err != nil {

		__LOG.Printf("%s | %s\n", "ERROR", err.Error())

		return err

	} else {

		__LOG.Printf("%s | %s \n\t%v\n", "INFO", "Created ConfigMap", res.Namespace+"/"+res.Name+" : "+res.CreationTimestamp.String())

		return nil
	}
}

// returns two lists of v1.ConfigMap (active configmaps and inactive configmaps)
func getConfigMaps(_namespace string) ([]v1.ConfigMap, []v1.ConfigMap) {

	if configmapsfromsourcens, err := __K8SCLIENT.CoreV1().ConfigMaps(_namespace).List(context.TODO(), metav1.ListOptions{}); err != nil {

		__LOG.Printf("%s | %s\n", "ERROR", err.Error())

		return []v1.ConfigMap{}, []v1.ConfigMap{}

	} else {

		activesc, inactivesc := []v1.ConfigMap{}, []v1.ConfigMap{}

		for _, configmap := range configmapsfromsourcens.Items {

			if configmap.Annotations[ANNOTATION_ACTIVE] == "true" {

				activesc = append(activesc, configmap)

			} else if configmap.Annotations[ANNOTATION_ACTIVE] == "false" {

				inactivesc = append(inactivesc, configmap)
			}
		}

		// listing the active configmaps
		__LOG.Printf("%s | %s\n", "INFO", "Active ConfigMaps from Namespace["+_namespace+"]:")
		__LOG.Printf("%s | %s\n", "----", "--------------------------------")

		for _, item := range activesc {

			__LOG.Printf("Name: %s\n", item.Name)
		}
		__LOG.Printf("%s | %s\n", "----", "--------------------------------")

		// listing the inactive configmaps
		__LOG.Printf("%s | %s\n", "INFO", "Inactive ConfigMaps from Namespace["+_namespace+"]:")
		__LOG.Printf("%s | %s\n", "----", "--------------------------------")

		for _, item := range inactivesc {

			__LOG.Printf("Name: %s\n", item.Name)
		}
		__LOG.Printf("%s | %s\n", "----", "--------------------------------")

		return activesc, inactivesc
	}
}

// Initcommand for ConfigMap Distribution
func InitConfigMapHandler(_sourcens string, _k8sclient *kubernetes.Clientset, _log *log.Logger, _avoidns []string, _immutablereplicas bool) {

	for {

		// receive relevant configmaps from the source namespace
		activesc, inactivesc := getConfigMaps(_sourcens)

		// handle active configmaps
		for _, configmap := range activesc {

			for _, namespace := range GetDistributeNamespaces(_avoidns, configmap.Annotations) {

				// get current configmap from namespace, if existent -> delete rist and then create
				if currconfigmap, err := __K8SCLIENT.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configmap.Name, metav1.GetOptions{}); err != nil {

					__LOG.Printf("%s | %s\n", "WARNING", "Namespace["+namespace+"] - "+err.Error())

					createConfigMap(configmap, namespace)

				} else {

					if configmap.ResourceVersion != currconfigmap.Annotations[ANNOTATION_ORIGINAL_RV] {

						if currconfigmap.Annotations[ANNOTATION_REPLICA] == "true" {

							if err := __K8SCLIENT.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), currconfigmap.Name, metav1.DeleteOptions{}); err != nil {

								__LOG.Printf("%s | %s\n", "ERROR", err.Error())

							} else {

								createConfigMap(configmap, namespace)
							}
						}
					}
				}
			}
		}

		// handle inactive configmaps
		for _, configmap := range inactivesc {

			__LOG.Printf("%s | %s\n", "INFO", "Remove inactive ConfigMap "+configmap.Name)

			for _, namespace := range GetAllNamespaces() {

				if namespace != _sourcens {

					if err := __K8SCLIENT.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), configmap.Name, metav1.DeleteOptions{}); err != nil {

						__LOG.Printf("%s | %s\n", "ERROR ["+namespace+"]", err.Error())

					} else {

						__LOG.Printf("%s | %s\n", "ERROR", "Deleted ConfigMap - "+configmap.Namespace+"/"+configmap.Name)
					}
				}
			}
		}

		time.Sleep(60 * time.Second)
	}
}
