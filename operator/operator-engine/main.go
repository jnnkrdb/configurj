package main

import (
	"github.com/jnnkrdb/corerdb/prtcl"
	"k8s.io/client-go/rest"
)

func main() {

	if config, err := rest.InClusterConfig(); err != nil {

		prtcl.PrintObject(config, err)

	} else {

	}

}
