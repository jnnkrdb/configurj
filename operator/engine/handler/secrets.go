package handler

import (
	"context"

	"github.com/jnnkrdb/configurj-engine/env"
	"github.com/jnnkrdb/configurj-engine/int/v1alpha1"

	"github.com/jnnkrdb/corerdb/fnc"
	"github.com/jnnkrdb/k8s/operator"
	"github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// function to create, update and delete secrets from a globalsecret
func CRUD_Secrets(gs v1alpha1.GlobalSecret) {

	crudlog := env.Log().WithFields(logrus.Fields{
		"current.globalsecret.name":      gs.Name,
		"current.globalsecret.namespace": gs.Namespace,
	})

	crudlog.Debug("processing globalconfig")

	var (
		all_namespaces     []string
		matched_namespaces []string
		err                error
	)

	crudlog.WithFields(logrus.Fields{
		"all_namespaces":     all_namespaces,
		"matched_namespaces": matched_namespaces,
		"err":                err,
	}).Trace("allocated cache for objects")

	crudlog.Debug("requesting namespace lists")
	// request the namespace lists (all namespaces, avoided via regex, matched via regex)
	all_namespaces, matched_namespaces, err = GetNamespaceLists(gs.Spec.Namespaces.AvoidRegex, gs.Spec.Namespaces.MatchRegex)

	crudlog.WithFields(logrus.Fields{
		"all_namespaces":     all_namespaces,
		"matched_namespaces": matched_namespaces,
		"err":                err,
	}).Trace("received namespaces lists")

	// DELETE
	crudlog.Debug("delete non-matching secrets")

	for _, clusternamespace := range all_namespaces {

		delTrace := crudlog.WithFields(logrus.Fields{
			"current.namespace":   clusternamespace,
			"matching.namespaces": matched_namespaces,
		})

		delTrace.Trace("checking namespace match")

		if !fnc.StringInList(clusternamespace, matched_namespaces) {

			delTrace.Trace("namespace does not match with required namespaces")

			_DeleteSecret(clusternamespace, gs.Spec.Name)

		} else {

			delTrace.Trace("namespace does match with required namespaces")
		}
	}

	crudlog.Debug("creating/updating matching secrets")

	for _, matchednamespace := range matched_namespaces {

		llog := crudlog.WithField("destination.secret", matchednamespace+"/"+gs.Spec.Name)

		llog.Trace("checking the existence of the secret")

		if scrt, err := operator.K8S().CoreV1().Secrets(matchednamespace).Get(context.TODO(), gs.Spec.Name, metav1.GetOptions{}); err != nil {

			llog.Trace("secret does not exist")

			// CREATE
			_CreateSecret(matchednamespace, gs)

		} else {

			llog.Trace("secret does exist")

			vercomplog := llog.WithFields(logrus.Fields{
				"secret.annotation.resourceversion":    scrt.Annotations[ANNOTATION_RESOURCEVERSION],
				"current.globalsecret.resourceversion": gs.ResourceVersion,
			})

			vercomplog.Trace("comparing secret versions")

			// UPDATE
			if scrt.Annotations[ANNOTATION_RESOURCEVERSION] != gs.ResourceVersion {

				vercomplog.Trace("versions do not match, updating secret")

				// delete the old secret
				if _DeleteSecret(matchednamespace, scrt.Name) == nil {

					_CreateSecret(matchednamespace, gs)
				}
			}
		}
	}
}

// create function for secrets
func _CreateSecret(namespace string, gs v1alpha1.GlobalSecret) (err error) {

	createlog := env.Log().WithFields(logrus.Fields{
		"destination.secret":  namespace + "/" + gs.Spec.Name,
		"source.globalsecret": gs.Namespace + "/" + gs.Name,
	})

	createlog.Debug("creating secret")

	createlog.Trace("initializing new secret in cache")

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

	createlog.WithField("secret", new).Trace("initialized new secret in cache")

	var scrt_result *v1.Secret

	scrt_result, err = operator.K8S().CoreV1().Secrets(new.Namespace).Create(context.TODO(), &new, metav1.CreateOptions{})

	createlog.WithFields(logrus.Fields{
		"secret.cached":         new,
		"secret.creationresult": *scrt_result,
	}).Debug("secrets content")

	if err != nil {

		createlog.WithError(err).Error("error while creating secret")

	} else {

		createlog.Debug("secret created")
	}

	return
}

// delete function for secrets
func _DeleteSecret(namespace, name string) (err error) {

	env.Log().WithField("delete.secret", namespace+"/"+name).Debug("deleting secret")

	err = operator.K8S().CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

	if err != nil {

		env.Log().WithField("delete.secret", namespace+"/"+name).WithError(err).Error("error deleting secret")

	} else {

		env.Log().WithField("delete.secret", namespace+"/"+name).Trace("deleted secret")
	}

	return
}
