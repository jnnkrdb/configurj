package handler

import (
	"github.com/jnnkrdb/corerdb/prtcl"
)

// goroutine for the handling of the crds
//
// globalconfigs/globalsecrets will be parsed into configmaps and secrets,
// which will be created, updated or deleted
func HandleCRDs(avoidNamespaces []string) {

	prtcl.Log.Println("beginning resource processing")

	// receive allowed namespaces
	allowednamespaces := calculateAllowedNamespaces(avoidNamespaces)
	prtcl.PrintObject(avoidNamespaces, allowednamespaces)

	// calculate globalconfigs
	gConfigs := get_GlobalConfigs()

	for _, gc := range gConfigs.Items {

		crud_Configmaps(allowednamespaces, gc)
	}

	// calculate globalsecrets
	gSecrets := get_GlobalSecrets()

	for _, gs := range gSecrets.Items {

		crud_Secrets(allowednamespaces, gs)
	}
}
