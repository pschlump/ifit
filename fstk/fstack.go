package fstk

/*
Copyright (C) Philip Schlump, 2016.

MIT Licensed.
*/

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type FileStackElement struct {
	S_LineNo int
	C_LineNo int
	File     *os.File
	Name     string // name of the item
	Scanner  *bufio.Scanner
}
type FileStackType struct {
	Stack []FileStackElement
	Top   int
}

func NewFileStackType() (rv *FileStackType) {
	return &FileStackType{Stack: make([]FileStackElement, 0, 10), Top: -1}
}

func (ns *FileStackType) IsEmpty() bool {
	return ns.Top <= -1
}

func (ns *FileStackType) Push(S, C int, fp *os.File, name string) {
	ns.Top++
	if len(ns.Stack) <= ns.Top {
		ns.Stack = append(ns.Stack, FileStackElement{S, C, fp, name, nil})
	} else {
		ns.Stack[ns.Top] = FileStackElement{S, C, fp, name, nil}
	}
}

var ErrEmptyStack = errors.New("Empty Stack")

func (ns *FileStackType) Peek() (FileStackElement, error) {
	if !ns.IsEmpty() {
		return ns.Stack[ns.Top], nil
	} else {
		return FileStackElement{}, ErrEmptyStack
	}
}

func (ns *FileStackType) Pop() {
	if !ns.IsEmpty() {
		ns.Top--
	}
}

func (ns *FileStackType) Length() int {
	return ns.Top + 1
}

func (ns *FileStackType) Dump1() {
	fmt.Printf("File Stack Dump 1\n")
	fmt.Printf("  Top = %d\n", ns.Top)
	for ii, vv := range ns.Stack {
		if ii <= ns.Top {
			fmt.Printf("   %d: Name [%s] LineNo: %d\n", ii, vv.Name, vv.C_LineNo)
		}
	}
}

func (ns *FileStackType) GetNames() (names []string) {
	for ii, vv := range ns.Stack {
		if ii <= ns.Top {
			names = append(names, vv.Name)
		}
	}
	return
}

func (ns *FileStackType) SetLineNo(n int) {
	if !ns.IsEmpty() {
		ns.Stack[ns.Top].C_LineNo = n
	}
}

func (ns *FileStackType) SetScanner(ss *bufio.Scanner) {
	if !ns.IsEmpty() {
		ns.Stack[ns.Top].Scanner = ss
	}
}
