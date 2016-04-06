package main

/*
Copyright (C) Philip Schlump, 2016.

MIT Licensed.
*/

import "flag"

// fopen
// Modifed from: "encoding/json"
// word parser

var InputFN = flag.String("input", "", "Input Meta File") // 0
var OutputFN = flag.String("output", "", "Output Code")   // 1
var SubFN = flag.String("sub", "", "Substution Values")   // 1
var Mode = flag.String("mode", "", "Mode Values")         // 1
var Debug = flag.Bool("debug", false, "Debug Flag")       // 2
func init() {
	flag.StringVar(InputFN, "i", "", "Input Meta File")
	flag.StringVar(OutputFN, "o", "", "Output Code")
	flag.StringVar(SubFN, "s", "", "Substution Values")
	flag.StringVar(Mode, "m", "", "Mode Values")
	flag.BoolVar(Debug, "D", false, "Debug Flag")
}
