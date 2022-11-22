package main

import (
	"time"

	"github.com/jnnkrdb/configurj-engine/core"
	"github.com/jnnkrdb/configurj-engine/handler"

	"github.com/jnnkrdb/corerdb/fnc"
	"github.com/jnnkrdb/corerdb/prtcl"
)

const CONFIGFILE string = "/config/settings.json"

type Settings struct {
	Debugging             bool     `json:"debugging"`
	GlobalAvoidNamespaces []string `json:"globalavoidnamespaces"`
	TimeOutSec            float64  `json:"timeoutsec"`
}

func main() {

	// activate logging
	prtcl.SetDebugging(true)

	// load service config, init with default confs
	conf := Settings{
		Debugging:             true,
		GlobalAvoidNamespaces: nil,
		TimeOutSec:            5,
	}

	if err := fnc.LoadStructFromFile("json", CONFIGFILE, &conf); err != nil {

		prtcl.Log.Println("failed loading service-configuration:", err)

		prtcl.PrintObject(CONFIGFILE, conf, err)

	} else {

		prtcl.Log.Println("successfully loaded service-configuration")

		prtcl.PrintObject(conf)

		prtcl.SetDebugging(conf.Debugging)

		// kill operator if crd configload fails
		prtcl.ErrorKill(core.LoadRestClient())

		// start the actual routine for the operator
		for {

			// runnning crds handler
			handler.HandleCRDs(conf.GlobalAvoidNamespaces)

			time.Sleep(time.Duration(conf.TimeOutSec) * time.Minute)
		}
	}
}
