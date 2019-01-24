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
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
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

func (r *ImportsResolver) clean() {
	r.FunctionsUsed = cleanResolved(r.FunctionsUsed)
	r.ClassesUsed = cleanResolved(r.ClassesUsed)
	r.FunctionsProvided = cleanResolved(r.FunctionsProvided)
	r.ClassesProvided = cleanResolved(r.ClassesProvided)
}

func cleanResolved(resolved []string) []string {
	resolved = slices.UniqueString(resolved)
	resolved = removeNativeTypes(resolved)
	resolved = slices.FilterString(resolved, IsEmpty)

	return resolved
}

func IsEmpty(s string) bool {
	return s != ""
}

func (r *ImportsResolver) addFunctionUsed(n node.Node) {
	fqn := r.nsr.ResolvedNames[n]
	r.FunctionsUsed = append(r.FunctionsUsed, fqn)
}

func (r *ImportsResolver) addClassUsed(n node.Node) {
	fqn := r.nsr.ResolvedNames[n]
	r.ClassesUsed = append(r.ClassesUsed, fqn)
}

func (r *ImportsResolver) addFunctionProvided(n node.Node) {
	fqn := r.nsr.ResolvedNames[n]
	r.FunctionsProvided = append(r.FunctionsProvided, fqn)
}

func (r *ImportsResolver) addClassProvided(n node.Node) {
	fqn := r.nsr.ResolvedNames[n]
	r.ClassesProvided = append(r.ClassesProvided, fqn)
}

func (r *ImportsResolver) uses(n node.Node) {
	fqn := r.nsr.ResolvedNames[n]

	if IsFunctionName(fqn) {
		r.addFunctionUsed(n)
		return
	}

	if IsClassName(fqn) {
		r.addClassUsed(n)
		return
	}
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
		for _, nn := range n.Uses {
			r.uses(nn)
		}

		// no reason to iterate into depth
		return false

	case *stmt.GroupUse:
		for _, nn := range n.UseList {
			r.uses(nn)
		}

		// no reason to iterate into depth
		return false

	case *stmt.Class:
		if n.Extends != nil {
			r.addClassUsed(n.Extends.ClassName)
		}

		if n.Implements != nil {
			for _, interfaceName := range n.Implements.InterfaceNames {
				r.addClassUsed(interfaceName)
			}
		}

		if n.ClassName != nil {
			r.addClassProvided(n)
		}

	case *stmt.Interface:
		if n.Extends != nil {
			for _, interfaceName := range n.Extends.InterfaceNames {
				r.addClassUsed(interfaceName)
			}
		}

		r.addClassProvided(n)

	case *stmt.Trait:
		r.addClassProvided(n)

	case *stmt.Function:
		r.addFunctionProvided(n)

		for _, parameter := range n.Params {
			r.addClassUsed(parameter)
		}

		if n.ReturnType != nil {
			r.addClassUsed(n.ReturnType)
		}

	case *stmt.ClassMethod:
		for _, parameter := range n.Params {
			r.addClassUsed(parameter)
		}

		if n.ReturnType != nil {
			r.addClassUsed(n.ReturnType)
		}

	case *expr.Closure:
		for _, parameter := range n.Params {
			r.addClassUsed(parameter)
		}

		if n.ReturnType != nil {
			r.addClassUsed(n.ReturnType)
		}

	case *stmt.ConstList:
		for _, constant := range n.Consts {
			r.addClassProvided(constant)
		}

	case *expr.StaticCall:
		r.addClassUsed(n.Class)

	case *expr.StaticPropertyFetch:
		r.addClassUsed(n.Class)

	case *expr.ClassConstFetch:
		r.addClassUsed(n.Class)

	case *expr.New:
		r.addClassUsed(n.Class)

	case *expr.InstanceOf:
		r.addClassUsed(n.Class)

	case *stmt.Catch:
		for _, t := range n.Types {
			r.addClassUsed(t)
		}

	case *expr.FunctionCall:
		r.addFunctionUsed(n.Function)

	case *expr.ConstFetch:
		r.addClassUsed(n.Constant)

	case *stmt.TraitUse:
		for _, t := range n.Traits {
			r.addClassUsed(t)
		}

		if n.TraitAdaptationList != nil {
			for _, a := range n.TraitAdaptationList.Adaptations {
				switch aa := a.(type) {
				case *stmt.TraitUsePrecedence:
					refTrait := aa.Ref.(*stmt.TraitMethodRef).Trait
					if refTrait != nil {
						r.addClassUsed(refTrait)
					}
					for _, insteadOf := range aa.Insteadof {
						r.addClassUsed(insteadOf)
					}

				case *stmt.TraitUseAlias:
					refTrait := aa.Ref.(*stmt.TraitMethodRef).Trait
					if refTrait != nil {
						r.addClassUsed(refTrait)
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
	// do nothing
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

	// Resolve fully-qualified names
	rootNode.Walk(resolver.nsr)

	// Resolve imports
	rootNode.Walk(resolver)
	resolver.clean()

	return resolver, nil
}
