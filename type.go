// Copyright (C) 2022, Connor Lane Smith <cls@lubutu.com>

package main

type Type interface {
	Graph() Graph
	Print()
}

type AtomicType struct{}

type FunctionType struct {
	Domain   Type
	Codomain Type
}

type NonlinearType struct {
	Base Type
}

func (t *AtomicType) Graph() Graph {
	arrow := new(Arrow)
	return &AtomicGraph{arrow}
}

func (t *AtomicType) Print() {
	print("0")
}

func (t *FunctionType) Graph() Graph {
	domain := t.Domain.Graph()
	codomain := t.Codomain.Graph()
	return &FunctionGraph{domain, codomain}
}

func (t *FunctionType) Print() {
	print("(")
	t.Domain.Print()
	print(" -o ")
	t.Codomain.Print()
	print(")")
}

func (t *NonlinearType) Graph() Graph {
	base := t.Base.Graph()
	return &NonlinearGraph{base}
}

func (t *NonlinearType) Print() {
	print("!")
	t.Base.Print()
}
