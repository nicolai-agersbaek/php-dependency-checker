package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/expr"
	"github.com/z7zmey/php-parser/node/name"
	"github.com/z7zmey/php-parser/node/stmt"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/visitor"
	"github.com/z7zmey/php-parser/walker"
	"os"
)

type ImportsResolver struct {
	nsr                                *visitor.NamespaceResolver
	FunctionsUsed, ClassesUsed         []string
	FunctionsProvided, ClassesProvided []string
}

func NewImportsResolver() *ImportsResolver {
	return &ImportsResolver{
		visitor.NewNamespaceResolver(),
		make([]string, 0),
		make([]string, 0),
		make([]string, 0),
		make([]string, 0),
	}
}

func (r *ImportsResolver) addProvidedFunction(fqn string) {
	r.FunctionsProvided = append(r.FunctionsProvided, fqn)
}

func (r *ImportsResolver) addProvidedClass(fqn string) {
	r.ClassesProvided = append(r.ClassesProvided, fqn)
}

func (r *ImportsResolver) EnterNode(w walker.Walkable) bool {
	switch n := w.(type) {
	case *stmt.Namespace:
		if n.NamespaceName == nil {
			r.nsr.Namespace = visitor.NewNamespace("")
		} else {
			NSParts := n.NamespaceName.(*name.Name).Parts
			r.nsr.Namespace = visitor.NewNamespace(concatNameParts(NSParts))
		}

	case *stmt.UseList:
		useType := ""
		if n.UseType != nil {
			useType = n.UseType.(*node.Identifier).Value
		}

		for _, nn := range n.Uses {
			r.nsr.AddAlias(useType, nn, nil)
		}

		// no reason to iterate into depth
		return false

	case *stmt.GroupUse:
		useType := ""
		if n.UseType != nil {
			useType = n.UseType.(*node.Identifier).Value
		}

		for _, nn := range n.UseList {
			r.nsr.AddAlias(useType, nn, n.Prefix.(*name.Name).Parts)
		}

		// no reason to iterate into depth
		return false

	case *stmt.Class:
		if n.Extends != nil {
			r.nsr.ResolveName(n.Extends.ClassName, "")
		}

		if n.Implements != nil {
			for _, interfaceName := range n.Implements.InterfaceNames {
				r.nsr.ResolveName(interfaceName, "")
			}
		}

		if n.ClassName != nil {
			r.nsr.AddNamespacedName(n, n.ClassName.(*node.Identifier).Value)
			r.addProvidedClass(r.nsr.ResolvedNames[n])
		}

	case *stmt.Interface:
		if n.Extends != nil {
			for _, interfaceName := range n.Extends.InterfaceNames {
				r.nsr.ResolveName(interfaceName, "")
			}
		}

		r.nsr.AddNamespacedName(n, n.InterfaceName.(*node.Identifier).Value)

	case *stmt.Trait:
		r.nsr.AddNamespacedName(n, n.TraitName.(*node.Identifier).Value)

	case *stmt.Function:
		r.nsr.AddNamespacedName(n, n.FunctionName.(*node.Identifier).Value)

		for _, parameter := range n.Params {
			r.nsr.ResolveType(parameter.(*node.Parameter).VariableType)
		}

		if n.ReturnType != nil {
			r.nsr.ResolveType(n.ReturnType)
		}

	case *stmt.ClassMethod:
		for _, parameter := range n.Params {
			r.nsr.ResolveType(parameter.(*node.Parameter).VariableType)
		}

		if n.ReturnType != nil {
			r.nsr.ResolveType(n.ReturnType)
		}

	case *expr.Closure:
		for _, parameter := range n.Params {
			r.nsr.ResolveType(parameter.(*node.Parameter).VariableType)
		}

		if n.ReturnType != nil {
			r.nsr.ResolveType(n.ReturnType)
		}

	case *stmt.ConstList:
		for _, constant := range n.Consts {
			r.nsr.AddNamespacedName(constant, constant.(*stmt.Constant).ConstantName.(*node.Identifier).Value)
		}

	case *expr.StaticCall:
		r.nsr.ResolveName(n.Class, "")

	case *expr.StaticPropertyFetch:
		r.nsr.ResolveName(n.Class, "")

	case *expr.ClassConstFetch:
		r.nsr.ResolveName(n.Class, "")

	case *expr.New:
		r.nsr.ResolveName(n.Class, "")

	case *expr.InstanceOf:
		r.nsr.ResolveName(n.Class, "")

	case *stmt.Catch:
		for _, t := range n.Types {
			r.nsr.ResolveName(t, "")
		}

	case *expr.FunctionCall:
		r.nsr.ResolveName(n.Function, "function")

	case *expr.ConstFetch:
		r.nsr.ResolveName(n.Constant, "const")

	case *stmt.TraitUse:
		for _, t := range n.Traits {
			r.nsr.ResolveName(t, "")
		}

		if n.TraitAdaptationList != nil {
			for _, a := range n.TraitAdaptationList.Adaptations {
				switch aa := a.(type) {
				case *stmt.TraitUsePrecedence:
					refTrait := aa.Ref.(*stmt.TraitMethodRef).Trait
					if refTrait != nil {
						r.nsr.ResolveName(refTrait, "")
					}
					for _, insteadOf := range aa.Insteadof {
						r.nsr.ResolveName(insteadOf, "")
					}

				case *stmt.TraitUseAlias:
					refTrait := aa.Ref.(*stmt.TraitMethodRef).Trait
					if refTrait != nil {
						r.nsr.ResolveName(refTrait, "")
					}
				}
			}
		}
	}

	return true
}

func (r *ImportsResolver) GetChildrenVisitor(Key string) walker.Visitor {
	return r
}

func (r *ImportsResolver) LeaveNode(w walker.Walkable) {
	r.nsr.LeaveNode(w)
}

func concatNameParts(parts ...[]node.Node) string {
	str := ""

	for _, p := range parts {
		for _, n := range p {
			if str == "" {
				str = n.(*name.NamePart).Value
			} else {
				str = str + "\\" + n.(*name.NamePart).Value
			}
		}
	}

	return str
}

func ResolveImports(path string) (*ImportsResolver, error) {
	src, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := src.Close(); err != nil {
			panic(err)
		}
	}()

	parser := php7.NewParser(src, path)
	parser.Parse()

	for _, e := range parser.GetErrors() {
		fmt.Println(e)
	}

	resolver := NewImportsResolver()
	rootNode := parser.GetRootNode()

	rootNode.Walk(resolver)

	return resolver, nil
}
