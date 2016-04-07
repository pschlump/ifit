package main

/*
Copyright (C) Philip Schlump, 2016.

MIT Licensed.
*/

import (
	"fmt"
	"sort"
	"strings"

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

func JsonStringToStringString(s string) (theJSON map[string]map[string]string, err error) {
	err = json.Unmarshal([]byte(s), &theJSON)
	if err != nil {
		theJSON = make(map[string]map[string]string)
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

func ParseNameValueOpt(s string) (name, value string) {
	if fv_re.MatchString(s) {
		ss := strings.Split(s, "=")
		name = ss[0]
		value = ss[1]
	} else if f_re.MatchString(s) {
		name = s
		value = "on"
	} else {
		name = s
		value = "on"
		fmt.Printf("ifit: Invalid command line options, should be Name or Name=Value, got >%s<\n", s)
	}
	return
}
func CommaList(strs []string) (s string) {
	s = ""
	com := ""
	for _, ii := range strs {
		s = s + com + ii
		com = ", "
	}
	return
}

func KeysSorted(sub map[string]string) (strs []string) {
	strs = make([]string, 0, 20)
	for ii := range sub {
		strs = append(strs, ii)
	}
	sort.Strings(strs)
	return
}
