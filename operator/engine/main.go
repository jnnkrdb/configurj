package main

import (
	"net/http"
	"time"

	"github.com/jnnkrdb/configurj-engine/env"
	"github.com/jnnkrdb/configurj-engine/handler"
	"github.com/jnnkrdb/configurj-engine/int/v1alpha1"

	"github.com/jnnkrdb/k8s/operator"
	"github.com/sirupsen/logrus"
)

var (

	// cached error
	err error = nil

	// lists of the resources

	gcList v1alpha1.GlobalConfigList
	gsList v1alpha1.GlobalSecretList
)

func main() {

	env.Log().WithFields(logrus.Fields{
		"supported-versions": []string{
			"globalsecrets." + v1alpha1.GroupName + "/" + v1alpha1.GroupVersion,
			"globalconfigs." + v1alpha1.GroupName + "/" + v1alpha1.GroupVersion,
		},
	}).Print("initializing configurj")

	if err = operator.InitK8sOperatorClient(); err != nil {

		env.Log().WithError(err).Error("error while initializing base kubernetes operator")

	} else {

		env.Log().Debug("successfully initialized base kubernetes operator")

		if err = operator.InitCRDOperatorRestClient(v1alpha1.GroupName, v1alpha1.GroupVersion, v1alpha1.AddToScheme); err != nil {

			env.Log().WithError(err).Error("error while initializing crds kubernetes operator")

		} else {

			env.Log().Debug("successfully initialized crds kubernetes operator")

			env.Log().WithField("globalconfiglist", gcList).Trace("allocated space for new globalconfiglist")

			env.Log().WithField("globalsecretlist", gsList).Trace("allocated space for new globalsecretlist")

			go func() {

				env.Log().WithField("TimeOutSeconds", env.TIMEOUTSECONDS).Debug("starting routine for resource processing")

				for {

					env.Log().Debug("requesting list of globalconfigs")

					if gcList, err = v1alpha1.GetGlobalConfigList(); err != nil {

						env.Log().WithError(err).Error("error receiving list of globalconfigs")

					} else {

						env.Log().WithField("globalconfiglist", gcList).Trace("received list of globalconfigs")

						env.Log().Debug("starting routine for globalconfigs")

						for _, gc := range gcList.Items {

							handler.CRUD_Configmaps(gc)
						}
					}

					// empty list and free storage
					env.Log().Trace("delete cached globalconfigs list")
					gcList = v1alpha1.GlobalConfigList{}

					env.Log().Trace("requesting list of globalsecrets")

					if gsList, err = v1alpha1.GetGlobalSecretList(); err != nil {

						env.Log().WithError(err).Error("error receiving list of globalsecrets")

					} else {

						env.Log().WithField("globalsecretlist", gsList).Trace("received list of globalsecrets")

						env.Log().Debug("starting routine for globalsecrets")

						for _, gs := range gsList.Items {

							handler.CRUD_Secrets(gs)
						}
					}

					// empty list and free storage
					env.Log().Trace("delete cached globalsecrets list")
					gsList = v1alpha1.GlobalSecretList{}

					env.Log().WithField("TimeOutSeconds", env.TIMEOUTSECONDS).Debugf("freezing routine for %v seconds", env.TIMEOUTSECONDS)

					time.Sleep(time.Duration(env.TIMEOUTSECONDS) * time.Second)
				}
			}()
		}
	}

	http.HandleFunc("/healthz/live", func(w http.ResponseWriter, r *http.Request) {

		env.Log().WithFields(logrus.Fields{
			"healthz.live":       true,
			"code":               200,
			"request.remoteaddr": r.RemoteAddr,
			"request.protocol":   r.Proto,
			"request.requesturi": r.RequestURI,
		}).Trace("health requested: liveness")

		w.WriteHeader(200)
	})

	http.HandleFunc("/healthz/ready", func(w http.ResponseWriter, r *http.Request) {

		env.Log().WithFields(logrus.Fields{
			"healthz.ready":      true,
			"code":               200,
			"request.remoteaddr": r.RemoteAddr,
			"request.protocol":   r.Proto,
			"request.requesturi": r.RequestURI,
		}).Trace("health requested: readyness")

		w.WriteHeader(200)
	})

	if err = http.ListenAndServe(":80", nil); err != nil {
		env.Log().WithField("endpoint", ":80").WithError(err).Error("http server finished")
	}
}
