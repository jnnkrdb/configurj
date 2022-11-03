package handler

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// creates the secret from a template (_secret) in the namespace (_namespace)
func createSecret(_secret v1.Secret, _namespace string) error {

	var temp = v1.Secret{}
	temp.Name = _secret.Name
	temp.Namespace = _namespace
	temp.Annotations = GetAnnotations(_secret.Namespace, _secret.Name, _secret.ResourceVersion)
	temp.Labels = _secret.Labels
	temp.Labels[LABELS_K8S_INSTANCE] = _namespace + "--" + temp.Name + "--" + _secret.ResourceVersion
	temp.Data = _secret.Data
	temp.StringData = _secret.StringData
	temp.Immutable = &__IMMUTABLE
	temp.Type = _secret.Type

	if res, err := __K8SCLIENT.CoreV1().Secrets(_namespace).Create(context.TODO(), &temp, metav1.CreateOptions{}); err != nil {

		__LOG.Printf("%s | %s\n", "ERROR", err.Error())

		return err

	} else {

		__LOG.Printf("%s | %s \n\t%v\n", "INFO", "Created Secret", res.Namespace+"/"+res.Name+" : "+res.CreationTimestamp.String())

		return nil
	}
}

// returns two lists of v1.Secret (active secrets and inactive secrets)
func getSecrets(_namespace string) ([]v1.Secret, []v1.Secret) {

	if secretsfromsourcens, err := __K8SCLIENT.CoreV1().Secrets(_namespace).List(context.TODO(), metav1.ListOptions{}); err != nil {

		__LOG.Printf("%s | %s\n", "ERROR", err.Error())

		return []v1.Secret{}, []v1.Secret{}

	} else {

		activesc, inactivesc := []v1.Secret{}, []v1.Secret{}

		for _, secret := range secretsfromsourcens.Items {

			if secret.Annotations[ANNOTATION_ACTIVE] == "true" {

				activesc = append(activesc, secret)

			} else if secret.Annotations[ANNOTATION_ACTIVE] == "false" {

				inactivesc = append(inactivesc, secret)
			}
		}

		// listing the active secrets
		__LOG.Printf("%s | %s\n", "INFO", "Active Secrets from Namespace["+_namespace+"]:")
		__LOG.Printf("%s | %s\n", "----", "--------------------------------")

		for _, item := range activesc {

			__LOG.Printf("Name: %s\n", item.Name)
		}
		__LOG.Printf("%s | %s\n", "----", "--------------------------------")

		// listing the inactive secrets
		__LOG.Printf("%s | %s\n", "INFO", "Inactive Secrets from Namespace["+_namespace+"]:")
		__LOG.Printf("%s | %s\n", "----", "--------------------------------")

		for _, item := range inactivesc {

			__LOG.Printf("Name: %s\n", item.Name)
		}
		__LOG.Printf("%s | %s\n", "----", "--------------------------------")

		return activesc, inactivesc
	}
}

// Initcommand for Secret Distribution
func InitSecretHandler(_avoidns []string) {

	for {

		// receive relevant secrets from the source namespace
		active, inactive := getSecrets(__SOURCENS)

		// handle active secrets
		for _, secret := range active {

			for _, namespace := range GetDistributeNamespaces(_avoidns, secret.Annotations) {

				// get current secret from namespace, if existent -> delete first and then create
				if currsecret, err := __K8SCLIENT.CoreV1().Secrets(namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{}); err != nil {

					__LOG.Printf("%s | %s\n", "WARNING", "Namespace["+namespace+"] - "+err.Error())

					createSecret(secret, namespace)

				} else {

					if currsecret.Annotations[ANNOTATION_REPLICA] == "true" {

						if secret.ResourceVersion != currsecret.Annotations[ANNOTATION_ORIGINAL_RV] {

							if err := __K8SCLIENT.CoreV1().Secrets(namespace).Delete(context.TODO(), currsecret.Name, metav1.DeleteOptions{}); err != nil {

								__LOG.Printf("%s | %s\n", "ERROR", err.Error())

							} else {

								createSecret(secret, namespace)
							}
						}
					}
				}
			}
		}

		// handle inactive secrets
		for _, secret := range inactive {

			__LOG.Printf("%s | %s\n", "INFO", "Remove inactive Secret "+secret.Name)

			for _, namespace := range GetAllNamespaces() {

				if namespace != __SOURCENS {

					if err := __K8SCLIENT.CoreV1().Secrets(namespace).Delete(context.TODO(), secret.Name, metav1.DeleteOptions{}); err != nil {

						__LOG.Printf("%s | %s\n", "ERROR ["+namespace+"]", err.Error())

					} else {

						__LOG.Printf("%s | %s\n", "ERROR", "Deleted Secret - "+secret.Namespace+"/"+secret.Name)
					}
				}
			}
		}

		time.Sleep(60 * time.Second)
	}
}
