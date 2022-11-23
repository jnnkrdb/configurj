package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jnnkrdb/configurj-engine/core"
	"github.com/jnnkrdb/configurj-engine/handler"
	"github.com/jnnkrdb/configurj-engine/int/v1alpha1"

	"github.com/jnnkrdb/corerdb/prtcl"
	"github.com/jnnkrdb/httprdb"
)

var (
	TimeOutMin float64 = 3
	Debugging  bool    = false
)

func main() {

	// initialize the env configs
	if debug, err := strconv.ParseBool(os.Getenv("DEBUGGING")); err == nil {
		Debugging = debug
	}

	if tom, err := strconv.ParseFloat(os.Getenv("TIMEOUTMINUTES"), 64); err == nil {
		TimeOutMin = tom
	}

	// activate logging
	prtcl.SetDebugging(Debugging)

	prtcl.Log.Println("successfully initialized globals.jnnkrdb.de-operator")

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

			time.Sleep(time.Duration(TimeOutMin) * time.Minute)
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
