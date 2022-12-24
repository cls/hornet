// Copyright (C) 2022, Connor Lane Smith <cls@lubutu.com>

package hornet

type Type interface {
	Graph() Graph
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

func (t *FunctionType) Graph() Graph {
	domain := t.Domain.Graph()
	codomain := t.Codomain.Graph()
	return &FunctionGraph{domain, codomain}
}

func (t *NonlinearType) Graph() Graph {
	base := t.Base.Graph()
	return &NonlinearGraph{base}
}
