// Copyright (C) 2022, Connor Lane Smith <cls@lubutu.com>

package main

type Rule struct {
	Lhs Term
	Rhs Term
}

type NodeSpec struct {
	Tag   int
	Arity int
}

type PortSpec struct {
	NodeIndex int
	PortIndex int
}

type EdgeSpec struct {
	Source PortSpec
	Target PortSpec
}

type GraphRule struct {
	CtorIndex int
	NodeSpecs []NodeSpec
	EdgeSpecs []EdgeSpec
}

type GraphRules struct {
	DtorIndex int
	ByCtorTag map[int]GraphRule
}

var ByDtorTag map[int]GraphRules

func (rule Rule) Print() {
	rule.Lhs.Print()
	print(" => ")
	rule.Rhs.Print()
	print("\n")
}

func (rule Rule) Compile() {
	lhs := rule.Lhs.Graph(0, nil)
	rhs := rule.Rhs.Graph(0, nil)

	context := lhs.Head()
	dtorPort := context.Port()
	dtor := dtorPort.Node
	var dtorRules GraphRules
	var ok bool
	if ByDtorTag == nil {
		ByDtorTag = make(map[int]GraphRules)
	} else {
		dtorRules, ok = ByDtorTag[dtor.Tag]
	}
	dtorIndex, ctorPort := findCtor(dtor)
	ctor := ctorPort.Node
	if !ok {
		dtorRules.DtorIndex = dtorIndex
		dtorRules.ByCtorTag = make(map[int]GraphRule)
	} else if dtorIndex != dtorRules.DtorIndex {
		panic("Unexpected dtor index")
	}
	ctorRules, ok := dtorRules.ByCtorTag[ctor.Tag]
	if ok {
		panic("Unexpected dtor-ctor pair")
	}

	nodes := make(map[*Node]int)
	nodes[dtor] = 0
	nodes[ctor] = 1
	lhs.Unify(rhs)
	lhs.ForEach(func(arrow *Arrow) {
		port := arrow.Port()
		collectNodes(port.Node, nodes)
	})

	ctorRules.CtorIndex = ctorPort.Index
	ctorRules.NodeSpecs = make([]NodeSpec, len(nodes))
	for node, i := range nodes {
		ctorRules.NodeSpecs[i] = NodeSpec{node.Tag, len(node.Ports)}
		for sourceIndex, targetPort := range node.Ports {
			j := nodes[targetPort.Node]
			if i < j {
				var edgeSpec EdgeSpec
				edgeSpec.Source = PortSpec{i, sourceIndex}
				edgeSpec.Target = PortSpec{j, targetPort.Index}
				ctorRules.EdgeSpecs = append(ctorRules.EdgeSpecs, edgeSpec)
			}
		}
	}

	dtorRules.ByCtorTag[ctor.Tag] = ctorRules
	ByDtorTag[dtor.Tag] = dtorRules
}

func collectNodes(node *Node, nodes map[*Node]int) {
	_, ok := nodes[node]
	if ok {
		return
	}
	nodes[node] = len(nodes)
	for _, port := range node.Ports {
		if port.Node != nil {
			collectNodes(port.Node, nodes)
		}
	}
}

func findCtor(dtor *Node) (int, Port) {
	for index, port := range dtor.Ports {
		return index, port
	}
	panic("Node not found")
}
