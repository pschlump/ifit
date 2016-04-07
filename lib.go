package main

/*
Copyright (C) Philip Schlump, 2016.

MIT Licensed.
*/

import (
	"fmt"
	"os"
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
	nthItem--
	if nthItem < len(w) && nthItem >= 0 {
		name = w[nthItem]
	}
	return
}

func GetItemSet(line string, nthItem int) (set []string) {
	w := ParseLineIntoWords(line)
	set = make([]string, 0, len(w))
	nthItem--
	// fmt.Printf("nthItem (after sub = %d, words are >%s<\n", nthItem, godebug.SVar(w))
	for ; nthItem < len(w); nthItem++ {
		name := w[nthItem]
		if name == "!!" {
			break
		}
		set = append(set, name)
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

// Use the search path to find a file
func FindFile(fn string, sp []string) (rv string) {
	for _, vv := range sp {
		s := vv + "/" + fn
		if Exists(s) {
			rv = s
			return
		}
	}
	return fn
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
