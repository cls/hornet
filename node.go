// Copyright (C) 2022, Connor Lane Smith <cls@lubutu.com>

package main

type Node struct {
	Tag   int
	Level int
	Ports []Port
}

type Port struct {
	Node  *Node
	Index int
}

const Fan = -4
const Croissant = -3
const Bracket = -2
const Eraser = -1

type Symbol struct {
	Name string
	Type Type
}

var Symbols []Symbol

func Tag(symbol Symbol) int {
	for tag, tagged := range Symbols {
		if tagged.Name == symbol.Name {
			return tag
		}
	}
	Symbols = append(Symbols, symbol)
	tag := len(Symbols) - 1
	return tag
}

func (lhs Port) Connect(rhs Port) {
	lhs.Node.Ports[lhs.Index] = rhs
	rhs.Node.Ports[rhs.Index] = lhs
}
