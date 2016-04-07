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
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pschlump/filelib"
	"github.com/pschlump/godebug" // fopen
	"github.com/pschlump/ifit/fstk"
	"github.com/pschlump/ifit/ifitlib"
	"github.com/pschlump/ifit/stk"
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
	PattType{"<!-- !! include ", 4, "include"},
	PattType{"// !! include ", 4, "include"},
	PattType{"/* !! include ", 4, "include"},
	PattType{"!! include ", 3, "include"},
	PattType{"<!-- !! include_once ", 4, "include_once"},
	PattType{"// !! include_once ", 4, "include_once"},
	PattType{"/* !! include_once ", 4, "include_once"},
	PattType{"!! include_once ", 3, "include_once"},
	PattType{"<!-- !! set_path ", -4, "set_path"},
	PattType{"// !! set_path ", -4, "set_path"},
	PattType{"/* !! set_path ", -4, "set_path"},
	PattType{"!! set_path ", -3, "set_path"},
	PattType{"<!-- !! define ", -4, "define"},
	PattType{"// !! define ", -4, "define"},
	PattType{"/* !! define ", -4, "define"},
	PattType{"!! define ", -3, "define"},
	PattType{"<!-- !! undef ", 4, "undef"},
	PattType{"// !! undef ", 4, "undef"},
	PattType{"/* !! undef ", 4, "undef"},
	PattType{"!! undef ", 3, "undef"},
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
	var openedOnece map[string]bool // for include_once
	openedOnece = make(map[string]bool)
	var SearchPath []string = []string{"./"}
	var PathOfInput string = ""
	var hasSub *regexp.Regexp
	hasSub = regexp.MustCompile("\\$\\$[a-zA-Z_][a-zA-Z0-9_]*\\$\\$")

	flag.Parse()
	fns := flag.Args()

	if len(fns) == 0 && *SubFN == "" {
		fmt.Fprintf(os.Stderr, "Requried option is missing\n")
		os.Exit(1)
	}

	fStack := fstk.NewFileStackType()

	fi, err := filelib.Fopen(*InputFN, "r")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file %s, Error: %s\n", *InputFN, err)
		os.Exit(1)
	}
	fStack.Push(1, 1, fi, *InputFN) // push the file at this point
	openedOnece[*InputFN] = true

	if filepath.IsAbs(*InputFN) || strings.Contains(*InputFN, string(os.PathSeparator)) {
		// if filepath.IsAbs(*InputFN) || strings.Contains(*InputFN, "/") {
		PathOfInput = filepath.Dir(filepath.Clean(*InputFN))
		PathOfInput += "/"
		// fmt.Printf("PathOfInput= [%s]\n", PathOfInput)
	}

	fo, err := filelib.Fopen(*OutputFN, "w")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output file %s, Error: %s\n", *OutputFN, err)
		os.Exit(1)
	}
	defer fo.Close()
	openedOnece[*OutputFN] = true // not certain but put the output file on the list of opened files, not good to include the output in the input.

	sub_top := make(map[string]map[string]string)
	sub := make(map[string]string)
	if *SubFN != "" {
		s, err := ioutil.ReadFile(*SubFN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening substitution JSON file %s, Error: %s\n", *SubFN, err)
			os.Exit(1)
		}
		sub_top, err = ifitlib.JsonStringToStringString(string(s))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON file %s, Error: %s\n", *SubFN, err)
			os.Exit(1)
		}
		var ok bool
		sub, ok = sub_top[*Mode]
		if !ok {
			fmt.Fprintf(os.Stderr, "ifit: Warning - mode %s not defined in %s - using an empty configuration\n", *Mode, *SubFN)
		}
	}

	for _, vv := range fns {
		// sub[vv] = "on"
		name, value, err := ifitlib.ParseNameValueOpt(vv)
		if err != nil {
			fmt.Printf("ifit: Invalid command line options. Error: %s Got: %s -- Assuming %s as name with vaolue of 'on'\n", err, vv, vv)
		}
		godebug.Printf(db2, "Option: [%s] name=[%s] value=[%s]\n", vv, name, value)
		sub[name] = value
	}

	outputOn := true
	ifStack := stk.NewNameStackType()               // Create the stack
	ifStack.Push(1, 1, outputOn, "**main**", false) // Push the empty frame - assume output on to start

	sub["__FILE__"] = *InputFN
	now := time.Now()
	sub["__DATE__"] = now.Format("2006-01-02")
	sub["__TIME__"] = now.Format("15:04:05")
	sub["__TSTAMP__"] = now.Format(time.RFC3339)
	sub["__Mode__"] = *Mode
	sub["__Output__"] = *OutputFN
	strs := ifitlib.KeysSorted(sub)
	sort.Strings(strs)
	sub["__TRUE_ITEMS__"] = ifitlib.CommaList(strs)
	sub["__PATH__"] = ifitlib.CommaList(SearchPath)
	stkNames := fStack.GetNames()
	sub["__OPENED_FILES__"] = ifitlib.CommaList(stkNames)

	// fmt.Printf("AT: %s\n", godebug.LF())
	var line_no = 1
	scanner := bufio.NewScanner(fi)
	fStack.SetScanner(scanner)
	for !fStack.IsEmpty() { // Loop until file stack is empty from pop.

		// fmt.Printf("AT: %s\n", godebug.LF())
		for ; scanner.Scan(); line_no++ {
			// fmt.Printf("AT: -at top of per-line - %s\n", godebug.LF())
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
				name := ifitlib.GetItemN(line[foundAt:], Pattern[pos].NthItem)
				if *Debug {
					fmt.Printf("pos=%v %s %s\n", pos, name, itemType)
				}
				if itemType == "include" || itemType == "include_once" {
					fname := ifitlib.FindFile(PathOfInput, name, SearchPath)
					if itemType == "include" || (itemType == "include_once" && !openedOnece[fname]) {
						godebug.Printf(db4, "Found %s [%s] - opeing file, %s\n", itemType, fname, godebug.LF())
						fStack.SetLineNo(line_no + 1)
						openedOnece[fname] = true
						ft, err := filelib.Fopen(fname, "r") // if it is an "include" then ... open file, push with new line number
						if err != nil {
							fmt.Fprintf(os.Stderr, "ifit: Error opening included input file %s, Error: %s\n", fname, err)
						}
						fStack.Push(1, 1, ft, fname) // push the file at this point
						fi = ft
						line_no = 0
						scanner = bufio.NewScanner(fi)
						fStack.SetScanner(scanner)
						sub["__FILE__"] = fname
						sub["__LINE__"] = fmt.Sprintf("%d", 1)
						stkNames := fStack.GetNames() // set __OPENED_FILES__
						sub["__OPENED_FILES__"] = ifitlib.CommaList(stkNames)
						if db4 {
							godebug.Printf(db4, "include - at bottom\n")
							fStack.Dump1()
						}
					}
				}
				if itemType == "define" {
					set := ifitlib.GetItemSet(line[foundAt:], -Pattern[pos].NthItem)
					if len(set) >= 2 {
						sub[set[0]] = set[1]
					} else if len(set) >= 1 {
						sub[set[0]] = "on"
					} else {
						fmt.Printf("ifit: Syntax error, define needs a name to define, line %d\n", line_no)
					}
				}
				if itemType == "undef" {
					if _, ok := sub[name]; ok {
						delete(sub, name)
					}
				}
				if itemType == "set_path" {
					// fmt.Printf("AT: %s\n", godebug.LF())
					set := ifitlib.GetItemSet(line[foundAt:], -Pattern[pos].NthItem)
					SearchPath = set
					// fmt.Printf("AT: set >%s< %s\n", godebug.SVar(set), godebug.LF())
					sub["__PATH__"] = ifitlib.CommaList(SearchPath)
					// fmt.Printf("AT: %s\n", godebug.LF())
				}
				if itemType == "if" {
					godebug.Printf(db1, "db: Found *if*/top, stack=%d, outputOn=%v, line_no=%d, %s\n", ifStack.Length(), outputOn, line_no, godebug.LF())
					if outputOn {
						_, inHash := sub[name]
						if inHash {
							if *Debug {
								fmt.Printf("Found in array %s\n", name)
							}
							outputOn = true
							ifStack.Push(line_no, line_no, outputOn, name, false) // Push the empty frame - assume output on to start
						} else {
							outputOn = false
							ifStack.Push(line_no, line_no, outputOn, name, false) // Push the empty frame - assume output on to start
						}
					} else {
						ifStack.Push(line_no, line_no, outputOn, name, true) // Push the empty frame - assume output on to start -- nested!
					}
					godebug.Printf(db1, "db: Found *if*/bot, stack=%d, outputOn=%v, line_no=%d, %s\n", ifStack.Length(), outputOn, line_no, godebug.LF())
				}
				if itemType == "else" {
					godebug.Printf(db1, "db: Found *else*/top, stack=%d, outputOn=%v, line_no=%d, %s\n", ifStack.Length(), outputOn, line_no, godebug.LF())
					x, err := ifStack.Peek()
					godebug.Printf(db1, "db: x from Peek on stack = %s\n", godebug.SVar(x))
					if err != nil {
						fmt.Fprintf(os.Stderr, "ifit: Error detected on line %d - else found with no if\n", line_no)
					} else if x.Tag == name {
						if !x.Nested {
							outputOn = !x.TF
						}
					} else {
						fmt.Fprintf(os.Stderr, "ifit: Error detected on line %d - mis matched else or invalid name on else, Started on line %d\n", line_no, x.S_LineNo)
						if db3 {
							fmt.Fprintf(os.Stderr, "ifit: x.Tag = [%s] name = [%s]\n", x.Tag, name)
							fmt.Fprintf(os.Stderr, "ifit: line = [%s] no = %d\n", line, line_no)
							ifStack.Dump1()
						}
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

		// fmt.Printf("AT: %s\n", godebug.LF())
		fi.Close()
		fStack.Pop() // Pop stack to restore previous file - loop

		if !fStack.IsEmpty() { // Loop until file stack is empty from pop.
			// fmt.Printf("AT: %s\n", godebug.LF())
			ff, _ := fStack.Peek() // peek to get name/line no back

			fi = ff.File
			scanner = ff.Scanner
			sub["__FILE__"] = ff.Name
			sub["__LINE__"] = fmt.Sprintf("%d", ff.C_LineNo)
			line_no = ff.C_LineNo

			stkNames := fStack.GetNames() // set __OPENED_FILES__
			sub["__OPENED_FILES__"] = ifitlib.CommaList(stkNames)

			if db4 {
				godebug.Printf(db4, "include - after pop\n")
				fStack.Dump1()
			}
		}

		// fmt.Printf("AT: %s\n", godebug.LF())
	}
	// fmt.Printf("AT: %s\n", godebug.LF())
}

const db1 = false // if/else/end
const db2 = false // command line name=value processing
const db3 = false // if/else/end - more details
const db4 = false // include related
