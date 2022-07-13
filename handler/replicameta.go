package handler

import (
	"strconv"
	"time"
)

const (
	// annotations

	_URI = "configurj.jnnkrdb.de"

	//replica

	ANNOTATION_REPLICA     = _URI + "/replica"     // value will be "true"
	ANNOTATION_TIMESTAMP   = _URI + "/timestamp"   // creates a timestamp "YYYY/MM/DD"
	ANNOTATION_ORIGINAL    = _URI + "/original"    // name of the original object
	ANNOTATION_ORIGINAL_NS = _URI + "/original-ns" // namespace of the original object
	ANNOTATION_ORIGINAL_RV = _URI + "/original-rv" // resourceversion of the original object

	// originals

	ANNOTATION_ACTIVE = _URI + "/active"
	ANNOTATION_AVOID  = _URI + "/avoid"
	ANNOTATION_MATCH  = _URI + "/match"

	// labels

	_K8S_URI            = "app.kubernetes.io"
	LABELS_K8S_INSTANCE = _K8S_URI + "/instance"
)

// returns the annotations that will be implemented in the replica
func GetAnnotations(_namespace, _name, _resourceversion string) map[string]string {

	annotations := make(map[string]string)

	annotations[ANNOTATION_REPLICA] = "true"

	annotations[ANNOTATION_TIMESTAMP] = strconv.Itoa(time.Now().Year()) + "/" + strconv.Itoa(int(time.Now().Month())) + "/" + strconv.Itoa(time.Now().Day())

	annotations[ANNOTATION_ORIGINAL] = _name

	annotations[ANNOTATION_ORIGINAL_NS] = _namespace

	annotations[ANNOTATION_ORIGINAL_RV] = _resourceversion

	return annotations
}
