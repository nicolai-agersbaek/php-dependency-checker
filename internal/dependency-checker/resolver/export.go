package resolver

import (
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/name"
	"github.com/z7zmey/php-parser/node/stmt"
	"github.com/z7zmey/php-parser/visitor"
	"github.com/z7zmey/php-parser/walker"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
)

type ExportResolver struct {
	visitor.NamespaceResolver
	Exports *Names
}

func NewExportResolver() *ExportResolver {
	return &ExportResolver{
		*visitor.NewNamespaceResolver(),
		NewNames(),
	}
}

func (r *ExportResolver) Clean() {
	r.Exports.Clean()
}

func (r *ExportResolver) addExport(n node.Node) {
	r.Exports.Add(r.resolveName(n))
}

func (r *ExportResolver) resolveName(nn node.Node) string {
	var nameParts []node.Node

	switch n := nn.(type) {
	case *stmt.Use:
		nameParts = n.Use.(*name.Name).Parts
	default:
		return r.ResolvedNames[n]
	}

	return ConcatNameParts(nameParts)
}

func (r *ExportResolver) EnterNode(w walker.Walkable) bool {
	switch n := w.(type) {
	case *stmt.Namespace:
		if n.NamespaceName == nil {
			r.Namespace = visitor.NewNamespace("")
		} else {
			NSParts := n.NamespaceName.(*name.Name).Parts
			nsName := ConcatNameParts(NSParts)
			r.Namespace = visitor.NewNamespace(nsName)
			r.Exports.AddNs(nsName)
		}

	case *stmt.Class:
		if n.ClassName != nil {
			r.AddNamespacedName(n, n.ClassName.(*node.Identifier).Value)
			r.addExport(n)
		}

	case *stmt.Interface:
		r.AddNamespacedName(n, n.InterfaceName.(*node.Identifier).Value)
		r.addExport(n)

	case *stmt.Trait:
		r.AddNamespacedName(n, n.TraitName.(*node.Identifier).Value)
		r.addExport(n)

	case *stmt.Function:
		r.AddNamespacedName(n, n.FunctionName.(*node.Identifier).Value)
		r.addExport(n)

	case *stmt.ConstList:
		for _, constant := range n.Consts {
			r.AddNamespacedName(constant, constant.(*stmt.Constant).ConstantName.(*node.Identifier).Value)
			r.addExport(constant)
		}
	}

	return true
}

// GetChildrenVisitor is invoked at every node parameter that contains children nodes
func (r *ExportResolver) GetChildrenVisitor(key string) walker.Visitor {
	return r
}

// LeaveNode is invoked after node process
func (r *ExportResolver) LeaveNode(w walker.Walkable) {
	switch n := w.(type) {
	case *stmt.Namespace:
		if n.Stmts != nil {
			r.Namespace = visitor.NewNamespace("")
		}
	}
}
