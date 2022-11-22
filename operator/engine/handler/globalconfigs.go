package handler

import (
	"context"

	"github.com/jnnkrdb/configurj-engine/core"
	"github.com/jnnkrdb/configurj-engine/int/v1alpha1"
	"github.com/jnnkrdb/corerdb/prtcl"
)

// request globalconfigs
func get_GlobalConfigs() v1alpha1.GlobalConfigList {

	prtcl.Log.Println("requesting list of globalconfigs")

	res := v1alpha1.GlobalConfigList{}

	// look up globalconfigs
	if err := core.CRDCLIENT.Get().Resource("globalconfigs").Do(context.TODO()).Into(&res); err != nil {

		prtcl.Log.Println("error receiving globalconfigs list:", err)

		prtcl.PrintObject(core.CRDCLIENT, res, err)
	}

	return res
}
