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
	"path/filepath"
)

const NamespaceSeparator = "\\"

type ImportsResolver struct {
	visitor.NamespaceResolver
	Imports *Names
	Exports *Names
}

func NewImportsResolver() *ImportsResolver {
	return &ImportsResolver{
		*visitor.NewNamespaceResolver(),
		NewNames(),
		NewNames(),
	}
}

func (r *ImportsResolver) clean() {
	r.Imports.clean()
	r.Exports.clean()
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

func (r *ImportsResolver) addImport(n node.Node) {
	r.Imports.Add(r.resolveName(n))
}

func (r *ImportsResolver) addExport(n node.Node) {
	r.Exports.Add(r.resolveName(n))
}

func (r *ImportsResolver) resolveName(nn node.Node) string {
	var nameParts []node.Node

	switch n := nn.(type) {
	case *stmt.Use:
		nameParts = n.Use.(*name.Name).Parts
	default:
		return r.ResolvedNames[n]
	}

	return concatNameParts(nameParts)
}

func (r *ImportsResolver) EnterNode(w walker.Walkable) bool {
	switch n := w.(type) {
	case *stmt.Namespace:
		if n.NamespaceName == nil {
			r.Namespace = visitor.NewNamespace("")
		} else {
			NSParts := n.NamespaceName.(*name.Name).Parts
			r.Namespace = visitor.NewNamespace(concatNameParts(NSParts))
		}

	case *stmt.UseList:
		useType := ""
		if n.UseType != nil {
			useType = n.UseType.(*node.Identifier).Value
		}

		for _, nn := range n.Uses {
			r.AddAlias(useType, nn, nil)
			r.addImport(nn)
		}

		// no reason to iterate into depth
		return false

	case *stmt.GroupUse:
		useType := ""
		if n.UseType != nil {
			useType = n.UseType.(*node.Identifier).Value
		}

		for _, nn := range n.UseList {
			r.AddAlias(useType, nn, n.Prefix.(*name.Name).Parts)
			r.addImport(nn)
		}

		// no reason to iterate into depth
		return false

	case *stmt.Class:
		if n.Extends != nil {
			r.ResolveName(n.Extends.ClassName, "")
			r.addImport(n.Extends.ClassName)
		}

		if n.Implements != nil {
			for _, interfaceName := range n.Implements.InterfaceNames {
				r.ResolveName(interfaceName, "")
				r.addImport(interfaceName)
			}
		}

		if n.ClassName != nil {
			r.AddNamespacedName(n, n.ClassName.(*node.Identifier).Value)
			r.addExport(n)
		}

	case *stmt.Interface:
		if n.Extends != nil {
			for _, interfaceName := range n.Extends.InterfaceNames {
				r.ResolveName(interfaceName, "")
				r.addImport(interfaceName)
			}
		}

		r.AddNamespacedName(n, n.InterfaceName.(*node.Identifier).Value)
		r.addExport(n)

	case *stmt.Trait:
		r.AddNamespacedName(n, n.TraitName.(*node.Identifier).Value)
		r.addExport(n)

	case *stmt.Function:
		r.AddNamespacedName(n, n.FunctionName.(*node.Identifier).Value)
		r.addExport(n)

		for _, parameter := range n.Params {
			r.ResolveType(parameter.(*node.Parameter).VariableType)
			r.addImport(parameter)
		}

		if n.ReturnType != nil {
			r.ResolveType(n.ReturnType)
			r.addImport(n.ReturnType)
		}

	case *stmt.ClassMethod:
		for _, parameter := range n.Params {
			r.ResolveType(parameter.(*node.Parameter).VariableType)
			r.addImport(parameter)
		}

		if n.ReturnType != nil {
			r.ResolveType(n.ReturnType)
			r.addImport(n.ReturnType)
		}

	case *expr.Closure:
		for _, parameter := range n.Params {
			r.ResolveType(parameter.(*node.Parameter).VariableType)
			r.addImport(parameter)
		}

		if n.ReturnType != nil {
			r.ResolveType(n.ReturnType)
			r.addImport(n.ReturnType)
		}

	case *stmt.ConstList:
		for _, constant := range n.Consts {
			r.AddNamespacedName(constant, constant.(*stmt.Constant).ConstantName.(*node.Identifier).Value)
			r.addExport(constant)
		}

	case *expr.StaticCall:
		r.ResolveName(n.Class, "")
		r.addImport(n.Class)

	case *expr.StaticPropertyFetch:
		r.ResolveName(n.Class, "")
		r.addImport(n.Class)

	case *expr.ClassConstFetch:
		r.ResolveName(n.Class, "")
		r.addImport(n.Class)

	case *expr.New:
		r.ResolveName(n.Class, "")
		r.addImport(n.Class)

	case *expr.InstanceOf:
		r.ResolveName(n.Class, "")
		r.addImport(n.Class)

	case *stmt.Catch:
		for _, t := range n.Types {
			r.ResolveName(t, "")
			r.addImport(t)
		}

	case *expr.FunctionCall:
		r.ResolveName(n.Function, "function")
		r.addImport(n.Function)

	case *expr.ConstFetch:
		r.ResolveName(n.Constant, "const")
		r.addImport(n.Constant)

	case *stmt.TraitUse:
		for _, t := range n.Traits {
			r.ResolveName(t, "")
			r.addImport(t)
		}

		if n.TraitAdaptationList != nil {
			for _, a := range n.TraitAdaptationList.Adaptations {
				switch aa := a.(type) {
				case *stmt.TraitUsePrecedence:
					refTrait := aa.Ref.(*stmt.TraitMethodRef).Trait
					if refTrait != nil {
						r.ResolveName(refTrait, "")
						r.addImport(refTrait)
					}
					for _, insteadOf := range aa.Insteadof {
						r.ResolveName(insteadOf, "")
						r.addImport(insteadOf)
					}

				case *stmt.TraitUseAlias:
					refTrait := aa.Ref.(*stmt.TraitMethodRef).Trait
					if refTrait != nil {
						r.ResolveName(refTrait, "")
						r.addImport(refTrait)
					}
				}
			}
		}
	}

	return true
}

// GetChildrenVisitor is invoked at every node parameter that contains children nodes
func (r *ImportsResolver) GetChildrenVisitor(key string) walker.Visitor {
	return r
}

// LeaveNode is invoked after node process
func (r *ImportsResolver) LeaveNode(w walker.Walkable) {
	switch n := w.(type) {
	case *stmt.Namespace:
		if n.Stmts != nil {
			r.Namespace = visitor.NewNamespace("")
		}
	}
}

func concatNameParts(parts ...[]node.Node) string {
	str := ""

	for _, p := range parts {
		for _, n := range p {
			if str == "" {
				str = n.(*name.NamePart).Value
			} else {
				str = str + NamespaceSeparator + n.(*name.NamePart).Value
			}
		}
	}

	return str
}

func ResolveDirImports(dir string) (*Names, *Names, error) {
	I, E := make([]*Names, 0), make([]*Names, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			imports, exports, err := ResolveFileImports(path)

			if err != nil {
				return err
			}

			I = append(I, imports)
			E = append(E, exports)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return mergeNames(I), mergeNames(E), nil
}

func mergeNames(names []*Names) *Names {
	merged := NewNames()

	for _, n := range names {
		merged = merged.Merge(n)
	}

	return merged
}

func ResolveFileImports(path string) (*Names, *Names, error) {
	src, err := os.Open(path)

	if err != nil {
		return nil, nil, err
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

	// Resolve imports
	rootNode.Walk(resolver)
	resolver.clean()

	return resolver.Imports, resolver.Exports, nil
}
