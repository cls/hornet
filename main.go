// Copyright (C) 2022, Connor Lane Smith <cls@lubutu.com>

package main

func main() {
	objType := &AtomicType{}
	dupType := &NonlinearType{objType}
	funType := &FunctionType{dupType, objType}
	appType := &FunctionType{objType, funType}
	lamType := &FunctionType{funType, objType}
	app := Symbol{"app", appType}
	lam := Symbol{"lam", lamType}
	lhsBody := &ApplyTerm{&ApplyTerm{&SymbolTerm{app}, &ApplyTerm{&SymbolTerm{lam}, &VarTerm{1}}}, &VarTerm{0}}
	lhs := &LambdaTerm{funType, &LambdaTerm{dupType, lhsBody}}
	rhsBody := &ApplyTerm{&VarTerm{1}, &VarTerm{0}}
	rhs := &LambdaTerm{funType, &LambdaTerm{dupType, rhsBody}}
	rule := Rule{lhs, rhs}
	rule.Print()
	rule.Compile()

	for dtorTag, dtorRules := range ByDtorTag {
		print(dtorTag)
		print(": (")
		print(dtorRules.DtorIndex)
		print(", {")
		var comma1 bool
		for ctorTag, ctorRules := range dtorRules.ByCtorTag {
			if comma1 {
				print(", ")
			}
			print(ctorTag)
			print(": (")
			print(ctorRules.CtorIndex)
			print(", [")
			var comma2 bool
			for _, nodeSpec := range ctorRules.NodeSpecs {
				if comma2 {
					print(", ")
				}
				print(nodeSpec.Tag)
				print("#")
				print(nodeSpec.Arity)
				comma2 = true
			}
			print("], [")
			var comma3 bool
			for _, edgeSpec := range ctorRules.EdgeSpecs {
				if comma3 {
					print(", ")
				}
				print(edgeSpec.Source.NodeIndex)
				print("#")
				print(edgeSpec.Source.PortIndex)
				print("--")
				print(edgeSpec.Target.NodeIndex)
				print("#")
				print(edgeSpec.Target.PortIndex)
				comma3 = true
			}
			print("])")
			comma1 = true
		}
		print("})\n")
	}
}
