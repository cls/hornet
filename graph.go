// Copyright (C) 2022, Connor Lane Smith <cls@lubutu.com>

package hornet

type Arrow struct {
	Port   Port
	parent *Arrow
}

func (arrow *Arrow) Connect(node *Node, index int) {
	root := arrow.Root()
	port := Port{node, index}
	if root.Port.Node != nil {
		root.Port.Connect(port)
	} else {
		root.Port = port
	}
}

func (arrow *Arrow) Root() *Arrow {
	for arrow.parent != nil {
		arrow = arrow.parent
	}
	return arrow
}

type Graph interface {
	ForEach(func(arrow *Arrow))
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

func (t *AtomicGraph) Map(f func(arrow *Arrow) *Arrow) Graph {
	arrow := f(t.Arrow)
	return &AtomicGraph{arrow}
}

func (lhs *AtomicGraph) Unify(graph Graph) {
	rhs := graph.(*AtomicGraph)
	lhsRoot := lhs.Arrow.Root()
	rhsRoot := rhs.Arrow.Root()
	if lhsRoot.Port.Node == nil {
		lhsRoot.parent = rhsRoot
	} else if rhsRoot.Port.Node == nil {
		rhsRoot.parent = lhsRoot
	} else {
		lhsRoot.Port.Connect(rhsRoot.Port)
	}
}

func (t *FunctionGraph) ForEach(f func(arrow *Arrow)) {
	t.Domain.ForEach(f)
	t.Codomain.ForEach(f)
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

func (t *NonlinearGraph) Map(f func(arrow *Arrow) *Arrow) Graph {
	base := t.Base.Map(f)
	return &NonlinearGraph{base}
}

func (lhs *NonlinearGraph) Unify(graph Graph) {
	rhs := graph.(*NonlinearGraph)
	lhs.Base.Unify(rhs.Base)
}
