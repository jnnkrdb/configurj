package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"configurj/handler"
	"configurj/probes"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Settings struct {
	ImmutableReplica bool     `json:"immutablereplicas"`
	HealthPort       string   `json:"healthport"`
	SourceNS         string   `json:"sourcenamespace"`
	AvoidSecrets     []string `json:"avoidsecrets"`
	AvoidConfigMaps  []string `json:"avoidconfigmaps"`
}

const (
	PREFIX       = "[jr] "
	settingsjson = "/configs/settings.json"
)

func main() {

	// setting liveness to 500, because of booting sevice
	// will be changed, when last handler starts
	probes.LIVENESS = 500

	// Setting Loogers for packages
	// Logging
	_LOG := log.New(os.Stdout, PREFIX, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	_LOG.Printf("%s | %s\n", "INFO", "Starting the Service ConfiguRJ")

	// loading settings.json
	settings := Settings{}
	// ------------------------------------------------------------------- change location of service config.json here
	if file, err := ioutil.ReadFile(settingsjson); err != nil {

		_LOG.Printf("%s | %s\n", "ERROR", err.Error())

	} else {

		if err := json.Unmarshal(file, &settings); err != nil {

			_LOG.Printf("%s | %s\n", "ERROR", err.Error())
		}
	}

	if settings.SourceNS == "" {

		_LOG.Printf("%s | %s\n", "WARNING", "No Source Namespace configured. Aborting Servicestart.")

	} else {

		// create k8s client
		_LOG.Printf("%s | %s\n", "INFO", "Loading InClusterConfig")
		if cnf, err := rest.InClusterConfig(); err != nil {

			_LOG.Printf("%s | %s\n", "ERROR", err.Error())

		} else {

			if cs, err := kubernetes.NewForConfig(cnf); err != nil {

				_LOG.Printf("%s | %s\n", "ERROR", err.Error())

			} else {

				// append sourcens to avoids
				settings.AvoidConfigMaps = append(settings.AvoidConfigMaps, settings.SourceNS)
				settings.AvoidSecrets = append(settings.AvoidSecrets, settings.SourceNS)

				// append empty string to avoids
				settings.AvoidConfigMaps = append(settings.AvoidConfigMaps, "")
				settings.AvoidSecrets = append(settings.AvoidSecrets, "")

				// Set Environment of the service
				handler.SetEnvironment(cs, _LOG, settings.ImmutableReplica, settings.SourceNS)

				// start namespacelisting
				go handler.AsyncWatchNamespaces(_LOG, cs)

				// init secret handler
				go handler.InitSecretHandler(settings.AvoidSecrets)

				// init configmap handler
				go handler.InitConfigMapHandler(settings.AvoidConfigMaps)

				// setting healthstatus to live
				probes.LIVENESS = 200
			}
		}
	}

	// starting the healthz-handler
	probes.StartHealthz(settings.HealthPort, _LOG)
}
