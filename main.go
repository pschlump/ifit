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
	"time"

	"github.com/pschlump/filelib"
	"github.com/pschlump/godebug" // fopen
)

type PattType struct {
	Pat      string
	NthItem  int
	ItemType string
}

var Pattern = []PattType{
	PattType{"<!-- !! if ", 4, "if"},
	PattType{"<!-- !! end ", 4, "end"},
	PattType{"<!-- !! else ", 4, "else"},
	PattType{"/* !! if ", 4, "if"},
	PattType{"/* !! end ", 4, "end"},
	PattType{"/* !! else ", 4, "else"},
	PattType{"// !! if ", 4, "if"},
	PattType{"// !! end ", 4, "end"},
	PattType{"// !! else ", 4, "else"},
	PattType{"!! if ", 3, "if"},
	PattType{"!! end ", 3, "end"},
	PattType{"!! else ", 3, "else"},
}

var fv_re *regexp.Regexp
var f_re *regexp.Regexp
var hasSub *regexp.Regexp

func init() {
	fv_re = regexp.MustCompile("([a-zA-Z][a-zA-Z_0-9]*)=(.*)")
	f_re = regexp.MustCompile("[a-zA-Z][a-zA-Z_0-9]*")
	hasSub = regexp.MustCompile("\\$\\$[a-zA-Z][a-zA-Z0-9_]*\\$\\$")
}

func SetFlag(s string) {

	var name, value string
	if fv_re.MatchString(s) {
		// xyzzy - pull out of r.e. match
	} else if f_re.MatchString(s) {
		name = s
		value = "on"
	} else {
		// error
		return
	}
	SetNameValue(name, value)
}

var g_ds map[string]map[string]string

func SetNameValue(name, value string) {
	if ds, ok := g_ds[*Mode]; ok {
		ds[name] = value
		g_ds[*Mode] = ds
	}
}

func IsSet(key, name string) bool {
	if vv, ok := g_ds[key]; ok {
		if _, ok1 := vv[name]; ok1 {
			return true
		}
	}
	return false
}

func IsSetValue(key, name string) string {
	if vv, ok := g_ds[key]; ok {
		if ww, ok1 := vv[name]; ok1 {
			return ww
		}
	}
	return ""
}

func HasIfItTag(s string) (patternNo int, foundAt int) {
	for ii, vv := range Pattern {
		if at := strings.Index(s, vv.Pat); at >= 0 {
			if *Debug {
				fmt.Printf("Found pat= -[%s]- at %d in input line -[%s]-, positon in line %d\n", vv.Pat, ii, s, at)
			}
			return ii, at
		}
	}
	return -1, -1
}

func main() {

	flag.Parse()
	fns := flag.Args()

	if len(fns) == 0 && *SubFN == "" {
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
	for _, vv := range fns {
		sub[vv] = "on"
		// xyzzy - parse fns and do name=value and just name at this point
	}

	outputOn := true
	ifStack := NewNameStackType()    // Create the stack
	ifStack.Push(1, 1, outputOn, "") // Push the empty frame - assume output on to start

	sub["__FILE__"] = *InputFN
	now := time.Now()
	sub["__DATE__"] = now.Format("2006-01-02")
	sub["__TIME__"] = now.Format("15:04:05")
	sub["__TSTAMP__"] = now.Format(time.RFC3339)
	sub["__Mode__"] = *Mode
	sub["__Output__"] = *OutputFN

	scanner := bufio.NewScanner(fi)
	for line_no := 1; scanner.Scan(); line_no++ {
		sub["__LINE__"] = fmt.Sprintf("%d", line_no)
		line := scanner.Text()
		if *Debug {
			fmt.Fprintf(fo, "%4d: %s\n", line_no, line)
		}
		if outputOn && hasSub.MatchString(line) {
			if *Debug {
				fmt.Printf("found matching line, %d\n", line_no)
			}
			line = hasSub.ReplaceAllStringFunc(line, func(in string) (out string) {
				in = in[2 : len(in)-2]
				// fmt.Printf("in [%s]\n", in)
				var ok bool
				if out, ok = sub[in]; !ok {
					fmt.Fprintf(os.Stderr, "ifit: Warning: substitution replacement for %s on line %d did not match - using empty string as replacment.", in, line_no)
				}
				return
			})
		}
		pos, foundAt := HasIfItTag(line)
		if pos >= 0 {
			itemType := Pattern[pos].ItemType
			name := GetItemN(line[foundAt:], Pattern[pos].NthItem)
			if *Debug {
				fmt.Printf("pos=%v %s %s\n", pos, name, itemType)
			}
			if itemType == "if" {
				godebug.Printf(db1, "db: Found *if*/top, stack=%d, outputOn=%v, line_no=%d, %s\n", ifStack.Length(), outputOn, line_no, godebug.LF())
				if outputOn {
					_, inHash := sub[name]
					// if InArray(name, fns) || inHash {
					// xyzzy- use function to check -
					// xyzzy- use Mode to get correct set
					if inHash {
						if *Debug {
							fmt.Printf("Found in array %s\n", name)
						}
						outputOn = true
						ifStack.Push(line_no, line_no, outputOn, "if") // Push the empty frame - assume output on to start
					} else {
						outputOn = false
						ifStack.Push(line_no, line_no, outputOn, "if") // Push the empty frame - assume output on to start
					}
				} else {
					ifStack.Push(line_no, line_no, outputOn, "if") // Push the empty frame - assume output on to start
				}
				godebug.Printf(db1, "db: Found *if*/bot, stack=%d, outputOn=%v, line_no=%d, %s\n", ifStack.Length(), outputOn, line_no, godebug.LF())
			}
			if itemType == "else" {
				godebug.Printf(db1, "db: Found *else*/top, stack=%d, outputOn=%v, %s\n", ifStack.Length(), outputOn, line_no, godebug.LF())
				x, err := ifStack.Peek()
				if err != nil {
					fmt.Fprintf(os.Stderr, "ifit: Error detected on line %d - else found with no if\n", line_no)
				} else if x.Tag == name {
					outputOn = !x.TF
				} else {
					fmt.Fprintf(os.Stderr, "ifit: Error detected on line %d - mis matched else or invalid name on else, Started on line %d\n", line_no, x.S_LineNo)
				}
				godebug.Printf(db1, "db: Found *else*/bot, stack=%d, outputOn=%v, line_no=%d, %s\n", ifStack.Length(), outputOn, line_no, godebug.LF())
			}
			if itemType == "end" {
				godebug.Printf(db1, "db: Found *end*/top, stack=%d, outputOn=%v, line_no=%d, %s\n", ifStack.Length(), outputOn, line_no, godebug.LF())
				if ifStack.Length() > 1 {
					ifStack.Pop()
					x, _ := ifStack.Peek()
					outputOn = x.TF
				} else {
					outputOn = true
					fmt.Fprintf(os.Stderr, "ifit: Error detected on line %d - mis matched end or extra end\n", line_no)
				}
				godebug.Printf(db1, "db: Found *end*/bot, stack=%d, outputOn=%v, line_no=%d, %s\n", ifStack.Length(), outputOn, line_no, godebug.LF())
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

const db1 = false
