package handler

import (
	"context"

	"github.com/jnnkrdb/configurj-engine/core"
	"github.com/jnnkrdb/configurj-engine/int/v1alpha1"
	"github.com/jnnkrdb/corerdb/prtcl"
)

// request global configs
func get_GlobalSecrets() v1alpha1.GlobalSecretList {

	prtcl.Log.Println("requesting list of globalsecrets")

	res := v1alpha1.GlobalSecretList{}

	// look up globalsecrets
	if err := core.CRDCLIENT.Get().Resource("globalsecrets").Do(context.TODO()).Into(&res); err != nil {

		prtcl.Log.Println("error receiving globalsecrets list:", err)

		prtcl.PrintObject(core.CRDCLIENT, res, err)
	}

	return res
}
