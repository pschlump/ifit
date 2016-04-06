package main

/*
Copyright (C) Philip Schlump, 2016.

MIT Licensed.
*/

import (
	"errors"
	"fmt"
)

type NameStackElement struct {
	S_LineNo int
	C_LineNo int
	TF       bool
	Tag      string // name of the item
	Nested   bool
}
type NameStackType struct {
	Stack []NameStackElement
	Top   int
}

func NewNameStackType() (rv *NameStackType) {
	return &NameStackType{Stack: make([]NameStackElement, 0, 10), Top: -1}
}

func (ns *NameStackType) IsEmpty() bool {
	return ns.Top <= -1
}

func (ns *NameStackType) Push(S, C int, tf bool, tag string, N bool) {
	ns.Top++
	if len(ns.Stack) <= ns.Top {
		ns.Stack = append(ns.Stack, NameStackElement{S, C, tf, tag, N})
	} else {
		ns.Stack[ns.Top] = NameStackElement{S, C, tf, tag, N}
	}
}

var ErrEmptyStack = errors.New("Empty Stack")

func (ns *NameStackType) Peek() (NameStackElement, error) {
	if !ns.IsEmpty() {
		return ns.Stack[ns.Top], nil
	} else {
		return NameStackElement{}, ErrEmptyStack
	}
}

func (ns *NameStackType) Pop() {
	if !ns.IsEmpty() {
		ns.Top--
	}
}

func (ns *NameStackType) Length() int {
	return ns.Top + 1
}

func (ns *NameStackType) Dump1() {
	fmt.Printf("Stack Dump 1\n")
	fmt.Printf("  Top = %d\n", ns.Top)
	for ii, vv := range ns.Stack {
		if ii <= ns.Top {
			fmt.Printf("   %d: Tag [%s] TF %v Nexted %v LineNo: %d\n", ii, vv.Tag, vv.TF, vv.Nested, vv.S_LineNo)
		}
	}
}
