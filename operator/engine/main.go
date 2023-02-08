package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/jnnkrdb/configurj-engine/core"
	"github.com/jnnkrdb/configurj-engine/env"
	"github.com/jnnkrdb/configurj-engine/handler"
	"github.com/jnnkrdb/configurj-engine/int/v1alpha1"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/jnnkrdb/corerdb/prtcl"
	"github.com/jnnkrdb/httprdb"
	"github.com/jnnkrdb/k8s/operator"
)

var (
	Debugging bool = false
)

func main() {

	env.Log().WithFields(logrus.Fields{
		"supported-versions": []string{
			"globalsecrets." + v1alpha1.GroupName + "/" + v1alpha1.GroupVersion,
			"globalconfigs." + v1alpha1.GroupName + "/" + v1alpha1.GroupVersion,
		},
	}).Print("initializing configurj")

	// initialize the env configs
	if debug, err := strconv.ParseBool(os.Getenv("DEBUGGING")); err == nil {
		Debugging = debug
	}

	// activate logging
	prtcl.SetDebugging(Debugging)

	var err error
	if err = operator.InitK8sOperatorClient(); err != nil {

		env.Log().WithError(err).Error("error while initializing base kubernetes operator")

	} else {

		env.Log().Debug("successfully initialized base kubernetes operator")

		if err = operator.InitCRDOperatorRestClient(v1alpha1.GroupName, v1alpha1.GroupVersion, v1alpha1.AddToScheme); err != nil {

			env.Log().WithError(err).Error("error while initializing crds kubernetes operator")

		} else {

			env.Log().Debug("successfully initialized crds kubernetes operator")

			go func() {

				env.Log().WithField("TimeOutSeconds", env.TIMEOUTSECONDS).Debug("starting routine for resource processing")

				for {

					time.Sleep(time.Duration(env.TIMEOUTSECONDS) * time.Second)
				}

			}()
		}
	}

	// kill operator if crd configload fails
	prtcl.ErrorKill(core.LoadRestClient())

	// start the actual routine for the operator
	go func() {

		prtcl.Log.Println("beginning resource processing")

		for {

			// calculate globalconfigs
			for _, gc := range func() v1alpha1.GlobalConfigList {

				prtcl.Log.Println("requesting list of globalsecrets")

				result := v1alpha1.GlobalConfigList{}

				// look up globalconfigs
				if err := core.CRDCLIENT.Get().Resource("globalconfigs").Do(context.TODO()).Into(&result); err != nil {

					prtcl.Log.Println("error receiving globalconfigs list:", err)

					prtcl.PrintObject(core.CRDCLIENT, result, err)
				}

				return result

			}().Items {

				handler.CRUD_Configmaps(gc)
			}

			// calculate globalsecrets
			for _, gs := range func() v1alpha1.GlobalSecretList {

				prtcl.Log.Println("requesting list of globalsecrets")

				result := v1alpha1.GlobalSecretList{}

				// look up globalsecrets
				if err := core.CRDCLIENT.Get().Resource("globalsecrets").Do(context.TODO()).Into(&result); err != nil {

					prtcl.Log.Println("error receiving globalsecrets list:", err)

					prtcl.PrintObject(core.CRDCLIENT, result, err)
				}

				return result

			}().Items {

				handler.CRUD_Secrets(gs)
			}

			time.Sleep(time.Duration(env.TIMEOUTSECONDS) * time.Second)
		}
	}()

	prtcl.Log.Println("initializing healthz-endpoint")

	healthapi := httprdb.CreateApiEndpoint(":80", "release", []httprdb.Route{
		{
			Request: "GET",
			SubPath: "/healthz/live",
			Handler: func(ctx *gin.Context) {
				ctx.IndentedJSON(200,
					struct {
						Code   int    `json:"code"`
						Status string `json:"status"`
					}{
						Code:   200,
						Status: "health/live: OK",
					})
			},
		},
	})

	healthapi.Boot()
}
