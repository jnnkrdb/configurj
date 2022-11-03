package handler

import (
	"log"

	"k8s.io/client-go/kubernetes"
)

var (

	// Kubernetes Interface
	__K8SCLIENT *kubernetes.Clientset

	// Logging architecture
	__LOG *log.Logger

	// Immmutable Replicas
	__IMMUTABLE bool = false

	// Source Namespace
	__SOURCENS string = ""
)

// set environment for the service
func SetEnvironment(_k8s *kubernetes.Clientset, _log *log.Logger, _immutable bool, _sourcens string) {

	__K8SCLIENT = _k8s
	__LOG = _log
	__IMMUTABLE = _immutable
	__SOURCENS = _sourcens
}
