// Copyright (C) 2022, Connor Lane Smith <cls@lubutu.com>

package main

type Arrow struct {
	port   Port
	parent *Arrow
}

func (arrow *Arrow) Connect(node *Node, index int) {
	root := arrow.Root()
	port := Port{node, index}
	if root.port.Node != nil {
		root.port.Connect(port)
	} else {
		root.port = port
	}
}

func (arrow *Arrow) Port() Port {
	return arrow.Root().port
}

func (arrow *Arrow) Root() *Arrow {
	for arrow.parent != nil {
		arrow = arrow.parent
	}
	return arrow
}

type Graph interface {
	ForEach(func(arrow *Arrow))
	Head() *Arrow
	Map(func(arrow *Arrow) *Arrow) Graph
	Unify(graph Graph)
}

type AtomicGraph struct {
	Arrow *Arrow
}

type FunctionGraph struct {
	Domain   Graph
	Codomain Graph
}

type NonlinearGraph struct {
	Base Graph
}

func (t *AtomicGraph) ForEach(f func(arrow *Arrow)) {
	f(t.Arrow)
}

func (t *AtomicGraph) Head() *Arrow {
	return t.Arrow
}

func (t *AtomicGraph) Map(f func(arrow *Arrow) *Arrow) Graph {
	arrow := f(t.Arrow)
	return &AtomicGraph{arrow}
}

func (lhs *AtomicGraph) Unify(graph Graph) {
	rhs := graph.(*AtomicGraph)
	lhsRoot := lhs.Arrow.Root()
	rhsRoot := rhs.Arrow.Root()
	if lhsRoot.port.Node == nil {
		lhsRoot.parent = rhsRoot
	} else if rhsRoot.port.Node == nil {
		rhsRoot.parent = lhsRoot
	} else {
		lhsRoot.port.Connect(rhsRoot.port)
	}
}

func (t *FunctionGraph) ForEach(f func(arrow *Arrow)) {
	t.Domain.ForEach(f)
	t.Codomain.ForEach(f)
}

func (t *FunctionGraph) Head() *Arrow {
	return t.Codomain.Head()
}

func (t *FunctionGraph) Map(f func(arrow *Arrow) *Arrow) Graph {
	domain := t.Domain.Map(f)
	codomain := t.Codomain.Map(f)
	return &FunctionGraph{domain, codomain}
}

func (lhs *FunctionGraph) Unify(graph Graph) {
	rhs := graph.(*FunctionGraph)
	lhs.Domain.Unify(rhs.Domain)
	lhs.Codomain.Unify(rhs.Codomain)
}

func (t *NonlinearGraph) ForEach(f func(arrow *Arrow)) {
	t.Base.ForEach(f)
}

func (t *NonlinearGraph) Head() *Arrow {
	panic("NonlinearGraph has no head")
}

func (t *NonlinearGraph) Map(f func(arrow *Arrow) *Arrow) Graph {
	base := t.Base.Map(f)
	return &NonlinearGraph{base}
}

func (lhs *NonlinearGraph) Unify(graph Graph) {
	rhs := graph.(*NonlinearGraph)
	lhs.Base.Unify(rhs.Base)
}
