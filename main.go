package main

/*
Copyright (C) Philip Schlump, 2016.

MIT Licensed.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pschlump/filelib" // fopen
	"github.com/pschlump/json"    // Modifed from: "encoding/json"
	"github.com/pschlump/pw"      // word parser
)

var InputFN = flag.String("input", "", "Input Meta File") // 0
var OutputFN = flag.String("output", "", "Output Code")   // 1
var SubFN = flag.String("sub", "", "Substution Values")   // 1
var Debug = flag.Bool("debug", false, "Debug Flag")       // 2
func init() {
	flag.StringVar(InputFN, "i", "", "Input Meta File")
	flag.StringVar(OutputFN, "o", "", "Output Code")
	flag.StringVar(SubFN, "s", "", "Substution Values")
	flag.BoolVar(Debug, "D", false, "Debug Flag")
}

type AFx func(s string) (string, string)

type PattType struct {
	Pat string
	Fx  AFx
}

var Patterns = []PattType{
	PattType{"<!-- !! if ", func(s string) (string, string) { return GetItemN(s, 4, "if") }},
	PattType{"<!-- !! end ", func(s string) (string, string) { return GetItemN(s, 4, "end") }},
	PattType{"!! if ", func(s string) (string, string) { return GetItemN(s, 3, "if") }},
	PattType{"!! end ", func(s string) (string, string) { return GetItemN(s, 3, "end") }},
	PattType{"/* !! if ", func(s string) (string, string) { return GetItemN(s, 4, "if") }},
	PattType{"/* !! end ", func(s string) (string, string) { return GetItemN(s, 4, "end") }},
	PattType{"// !! if ", func(s string) (string, string) { return GetItemN(s, 4, "if") }},
	PattType{"// !! end ", func(s string) (string, string) { return GetItemN(s, 4, "end") }},
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
func GetItemN(line string, pos int, tag string) (name, rv string) {
	rv = "**error**"
	w := ParseLineIntoWords(line)
	if len(w) >= pos {
		name = w[pos-1]
		rv = tag
	}
	return
}

func HasPrefix(s string) (int, AFx) {
	for ii, vv := range Patterns {
		if strings.HasPrefix(s, vv.Pat) {
			if *Debug {
				fmt.Printf("Found %s at %d in >>%s<<\n", vv, ii, s)
			}
			return ii, vv.Fx
		}
	}
	return -1, nil
}

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

func main() {

	flag.Parse()
	fns := flag.Args()

	if len(fns) == 0 {
		fmt.Fprintf(os.Stderr, "Requried option is missing\n")
		os.Exit(1)
	}

	fi, err := filelib.Fopen(*InputFN, "r")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file %s, Error: %s\n", *InputFN, err)
		os.Exit(1)
	}
	defer fi.Close()

	fo, err := filelib.Fopen(*OutputFN, "w")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output file %s, Error: %s\n", *OutputFN, err)
		os.Exit(1)
	}
	defer fo.Close()

	sub := make(map[string]string)
	if *SubFN != "" {
		s, err := ioutil.ReadFile(*SubFN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening substitution JSON file %s, Error: %s\n", *SubFN, err)
			os.Exit(1)
		}
		sub, err = JsonStringToString(string(s))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON file %s, Error: %s\n", *SubFN, err)
			os.Exit(1)
		}
	}
	_ = sub

	outputOn := true
	hasSub := regexp.MustCompile("\\$\\$[a-zA-Z][a-zA-Z0-9_]*\\$\\$")

	scanner := bufio.NewScanner(fi)
	for line_no := 1; scanner.Scan(); line_no++ {
		line := scanner.Text()
		if *Debug {
			fmt.Fprintf(fo, "%4d: %s\n", line_no, line)
		}
		if hasSub.MatchString(line) {
			if *Debug {
				fmt.Printf("found matching line, %d\n", line_no)
			}
			line = hasSub.ReplaceAllStringFunc(line, func(in string) (out string) {
				in = in[2 : len(in)-2]
				// fmt.Printf("in [%s]\n", in)
				var ok bool
				if out, ok = sub[in]; !ok {
					fmt.Fprintf(os.Stderr, "Warning: substitution replacement for %s on line %d did not match - using empty string as replacment.", in, line_no)
				}
				return
			})
		}
		pos, fx := HasPrefix(line)
		if pos >= 0 {
			name, itemType := fx(line)
			if *Debug {
				fmt.Printf("pos=%v %s %s\n", pos, name, itemType)
			}
			if itemType == "if" {
				if InArray(name, fns) {
					if *Debug {
						fmt.Printf("Found in array %s\n", name)
					}
					outputOn = true
				} else {
					outputOn = false
				}
			}
			if itemType == "end" {
				outputOn = true
			}
		} else if outputOn {
			// fmt.Fprintf(fo, "%4s: %s\n", "++++", line)
			fmt.Fprintf(fo, "%s\n", line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
