package handler

import (
	"regexp"
)

const ANNOTATION_RESOURCEVERSION string = "configurj.jnnkrdb.de/version"

// --- helper functions

// find a string in a list of strings
func StringInList(searchstring string, listofstring []string) bool {

	for index := range listofstring {

		if listofstring[index] == searchstring {

			return true
		}
	}

	return false
}

// find a string in a list full of regex
func FindStringInRegexpList(compare string, listOfRegexp []string) (result bool, err error) {

	for index := range listOfRegexp {

		if result, err = regexp.MatchString(listOfRegexp[index], compare); result {

			break
		}
	}

	return
}
