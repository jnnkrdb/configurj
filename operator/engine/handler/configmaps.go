package handler

import (
	"context"

	"github.com/jnnkrdb/configurj-engine/env"
	"github.com/jnnkrdb/configurj-engine/int/v1alpha1"

	"github.com/jnnkrdb/k8s/operator"
	"github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// function to create, update and delete configmaps from a globalconfig
func CRUD_Configmaps(gc v1alpha1.GlobalConfig) {

	crudlog := env.Log().WithFields(logrus.Fields{
		"current.globalconfig.name":      gc.Name,
		"current.globalconfig.namespace": gc.Namespace,
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
	all_namespaces, matched_namespaces, err = GetNamespaceLists(gc.Spec.Namespaces.AvoidRegex, gc.Spec.Namespaces.MatchRegex)

	crudlog.WithFields(logrus.Fields{
		"all_namespaces":     all_namespaces,
		"matched_namespaces": matched_namespaces,
		"err":                err,
	}).Trace("received namespaces lists")

	// DELETE
	crudlog.Debug("delete non-matching configmaps")

	for _, clusternamespace := range all_namespaces {

		delTrace := crudlog.WithFields(logrus.Fields{
			"current.namespace":   clusternamespace,
			"matching.namespaces": matched_namespaces,
		})

		delTrace.Trace("checking namespace match")

		if !StringInList(clusternamespace, matched_namespaces) {

			delTrace.Trace("namespace does not match with required namespaces")

			_DeleteConfigMap(clusternamespace, gc.Spec.Name)

		} else {

			delTrace.Trace("namespace does match with required namespaces")
		}
	}

	crudlog.Debug("creating/updating matching configmaps")

	for _, matchednamespace := range matched_namespaces {

		llog := crudlog.WithField("destination.configmap", matchednamespace+"/"+gc.Spec.Name)

		llog.Trace("checking the existence of the configmap")

		if cm, err := operator.K8S().CoreV1().ConfigMaps(matchednamespace).Get(context.TODO(), gc.Spec.Name, metav1.GetOptions{}); err != nil {

			llog.Trace("configmap does not exist")

			// CREATE
			_CreateConfigMap(matchednamespace, gc)

		} else {

			llog.Trace("configmap does exist")

			vercomplog := llog.WithFields(logrus.Fields{
				"configmap.annotation.resourceversion":    cm.Annotations[ANNOTATION_RESOURCEVERSION],
				"current.globalconfigmap.resourceversion": gc.ResourceVersion,
			})

			vercomplog.Trace("comparing configmap versions")

			// UPDATE
			if cm.Annotations[ANNOTATION_RESOURCEVERSION] != gc.ResourceVersion {

				vercomplog.Trace("versions do not match, updating configmap")

				// delete the old configmap
				if _DeleteConfigMap(matchednamespace, cm.Name) == nil {

					_CreateConfigMap(matchednamespace, gc)
				}
			}
		}
	}
}

// create function for configmaps
func _CreateConfigMap(namespace string, gc v1alpha1.GlobalConfig) (err error) {

	createlog := env.Log().WithFields(logrus.Fields{
		"destination.configmap": namespace + "/" + gc.Spec.Name,
		"source.globalconfig":   gc.Namespace + "/" + gc.Name,
	})

	createlog.Debug("creating configmap")

	createlog.Trace("initializing new configmap in cache")

	var new v1.ConfigMap

	new.Name = gc.Spec.Name

	new.Namespace = namespace

	new.Annotations = func() map[string]string {

		result := make(map[string]string)

		result[ANNOTATION_RESOURCEVERSION] = gc.ResourceVersion

		return result
	}()

	new.Immutable = &gc.Spec.Immutable

	new.Data = gc.Spec.Data

	createlog.WithField("configmap", new).Trace("initialized new configmap in cache")

	var cm_result *v1.ConfigMap

	cm_result, err = operator.K8S().CoreV1().ConfigMaps(new.Namespace).Create(context.TODO(), &new, metav1.CreateOptions{})

	createlog.WithFields(logrus.Fields{
		"configmap.cached":         new,
		"configmap.creationresult": *cm_result,
	}).Debug("configmaps content")

	if err != nil {

		createlog.WithError(err).Error("error while creating configmap")

	} else {

		createlog.Debug("configmap created")
	}

	return
}

// delete function for configmaps
func _DeleteConfigMap(namespace, name string) (err error) {

	env.Log().WithField("delete.configmap", namespace+"/"+name).Debug("deleting configmap")

	var cm *v1.ConfigMap

	if cm, err = operator.K8S().CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{}); err != nil {

		env.Log().WithField("delete.configmap", namespace+"/"+name).Trace("configmap does not exist")

	} else {

		err = operator.K8S().CoreV1().ConfigMaps(cm.Namespace).Delete(context.TODO(), cm.Name, metav1.DeleteOptions{})

		if err != nil {

			env.Log().WithField("delete.configmap", cm.Namespace+"/"+cm.Name).WithError(err).Error("error deleting configmap")

		} else {

			env.Log().WithField("delete.configmap", cm.Namespace+"/"+cm.Name).Trace("deleted configmap")
		}
	}

	return
}
