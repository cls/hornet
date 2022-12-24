// Copyright (C) 2022, Connor Lane Smith <cls@lubutu.com>

package hornet

type Term interface {
	Graph(level int, basis []Graph) Graph
}

type SymbolTerm struct {
	Symbol Symbol
}

type VarTerm struct {
	Index int
}

type LambdaTerm struct {
	BindType Type
	Body     Term
}

type ApplyTerm struct {
	Function Term
	Argument Term
}

type PromoteTerm struct {
	Values []Term
	Body   Term
}

type DerelictTerm struct {
	Value Term
}

type CopyTerm struct {
	Value Term
	Body  Term
}

type DiscardTerm struct {
	Value Term
	Body  Term
}

func (t *SymbolTerm) Graph(level int, basis []Graph) Graph {
	node := new(Node)
	node.Tag = Tag(t.Symbol)
	node.Level = level
	symbol := t.Symbol.Type.Graph()
	return symbol.Map(func(arrow *Arrow) *Arrow {
		node.Ports = append(node.Ports, Port{})
		arrow.Connect(node, len(node.Ports)-1)
		return arrow
	})
}

func (t *VarTerm) Graph(level int, basis []Graph) Graph {
	return basis[t.Index]
}

func (t *LambdaTerm) Graph(level int, basis []Graph) Graph {
	bind := t.BindType.Graph()
	body := t.Body.Graph(level, append(basis, bind))
	return &FunctionGraph{bind, body}
}

func (t *ApplyTerm) Graph(level int, basis []Graph) Graph {
	function := t.Function.Graph(level, basis).(*FunctionGraph)
	argument := t.Argument.Graph(level, basis)
	function.Domain.Unify(argument)
	return function.Codomain
}

func (t *PromoteTerm) Graph(level int, basis []Graph) Graph {
	newBasis := make([]Graph, len(t.Values))
	for i, tvalue := range t.Values {
		value := tvalue.Graph(level, basis)
		newBasis[i] = value.Map(func(arrow *Arrow) *Arrow {
			node := new(Node)
			node.Tag = Bracket
			node.Level = level
			node.Ports = make([]Port, 2)
			arrow.Connect(node, 0)
			newArrow := new(Arrow)
			newArrow.Connect(node, 1)
			return newArrow
		})
	}
	body := t.Body.Graph(level+1, newBasis)
	return body
}

func (t *DerelictTerm) Graph(level int, basis []Graph) Graph {
	value := t.Value.Graph(level, basis)
	return value.Map(func(arrow *Arrow) *Arrow {
		node := new(Node)
		node.Tag = Croissant
		node.Level = level
		node.Ports = make([]Port, 2)
		arrow.Connect(node, 0)
		newArrow := new(Arrow)
		newArrow.Connect(node, 1)
		return newArrow
	})
}

func (t *CopyTerm) Graph(level int, basis []Graph) Graph {
	value := t.Value.Graph(level, basis)
	left := value.Map(func(arrow *Arrow) *Arrow {
		node := new(Node)
		node.Tag = Fan
		node.Level = level
		node.Ports = make([]Port, 3)
		arrow.Connect(node, 0)
		newArrow := new(Arrow)
		newArrow.Connect(node, 1)
		return newArrow
	})
	right := left.Map(func(arrow *Arrow) *Arrow {
		node := arrow.Port.Node
		newArrow := new(Arrow)
		newArrow.Connect(node, 2)
		return newArrow
	})
	body := t.Body.Graph(level, append(basis, left, right))
	return body
}

func (t *DiscardTerm) Graph(level int, basis []Graph) Graph {
	value := t.Value.Graph(level, basis)
	value.ForEach(func(arrow *Arrow) {
		node := new(Node)
		node.Tag = Eraser
		node.Level = level
		node.Ports = make([]Port, 1)
		arrow.Connect(node, 0)
	})
	body := t.Body.Graph(level, basis)
	return body
}
