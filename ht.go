package main

/*
Copyright (C) Philip Schlump, 2012-2021.

BSD 3 Clause Licensed.
*/

import (
	"fmt"

	"github.com/pschlump/HashStr"
	"github.com/pschlump/pluto/comparable"
	hash_tab "github.com/pschlump/pluto/hash_tab_dll"
)

// comparable "github.com/pschlump/pluto/comparable"
// TestData is an Inteface Matcing data type for the Nodes that supports the Comparable
// interface.  This means that it has a Compare fucntion.

// type TestData struct {
// 	S string
// }

// At compile time verify that this is a correct type/interface setup.
var _ comparable.Comparable = (*DefinedItem)(nil)
var _ hash_tab.Hashable = (*DefinedItem)(nil)
var _ comparable.Equality = (*DefinedItem)(nil)

// Compare implements the Compare function to satisfy the interface requirements.
func (aa DefinedItem) Compare(x comparable.Comparable) int {
	if bb, ok := x.(DefinedItem); ok {
		if aa.Name < bb.Name {
			return -1
		} else if aa.Name > bb.Name {
			return 1
		}
	} else if bb, ok := x.(*DefinedItem); ok {
		if aa.Name < bb.Name {
			return -1
		} else if aa.Name > bb.Name {
			return 1
		}
	} else {
		panic(fmt.Sprintf("Passed invalid type %T to a Compare function.", x))
	}
	return 0
}

func (aa DefinedItem) IsEqual(x comparable.Equality) bool {
	if bb, ok := x.(DefinedItem); ok {
		if aa.Name == bb.Name {
			return true
		}
		return false
	} else if bb, ok := x.(*DefinedItem); ok {
		if aa.Name == bb.Name {
			return true
		}
		return false
	} else {
		panic(fmt.Sprintf("Passed invalid type %T to a Compare function.", x))
	}
	return false
}

func (aa DefinedItem) HashKey(x interface{}) (rv int) {
	if v, ok := x.(*DefinedItem); ok {
		rv = HashStr.HashStr([]byte(v.Name))
		return
	}
	if v, ok := x.(DefinedItem); ok {
		rv = HashStr.HashStr([]byte(v.Name))
		return
	}
	return
}
