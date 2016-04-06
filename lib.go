package main

/*
Copyright (C) Philip Schlump, 2016.

MIT Licensed.
*/

import (
	"github.com/pschlump/json"
	"github.com/pschlump/pw"
) // Modifed from: "encoding/json"

func InArray(lookFor string, inArr []string) bool {
	for _, v := range inArr {
		if lookFor == v {
			return true
		}
	}
	return false
}

func JsonStringToString(s string) (theJSON map[string]string, err error) {
	err = json.Unmarshal([]byte(s), &theJSON)
	if err != nil {
		theJSON = make(map[string]string)
	}
	return
}

func ParseLineIntoWords(line string) []string {
	// rv := strings.Fields ( line )
	Pw := pw.NewParseWords()
	Pw.SetOptions("C", true, true)
	Pw.SetLine(line)
	rv := Pw.GetWords()
	return rv
}

// func GetItemN(s,4,"if") {
func GetItemN(line string, nthItem int) (name string) {
	w := ParseLineIntoWords(line)
	if len(w) >= nthItem {
		name = w[nthItem-1]
	}
	return
}
