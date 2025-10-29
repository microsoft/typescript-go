package customlint

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"slices"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var unexportedAPIAnalyzer = &analysis.Analyzer{
	Name: "unexportedapi",
	Doc:  "finds exporrted APIs referencing unexported identifiers",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: func(pass *analysis.Pass) (any, error) {
		return (&unexportedAPIPass{pass: pass}).run()
	},
}

type unexportedAPIPass struct {
	pass     *analysis.Pass
	inspect  *inspector.Inspector
	file     *ast.File
	currDecl ast.Node
}

func (u *unexportedAPIPass) run() (any, error) {
	inspect := u.pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.TypeSpec)(nil),
		(*ast.ValueSpec)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		u.currDecl = n
		switch n := n.(type) {
		case *ast.File:
			u.file = n
		case *ast.FuncDecl:
			u.checkFuncDecl(n)
		case *ast.TypeSpec:
			u.checkTypeSpec(n)
		case *ast.ValueSpec:
			u.checkValueSpec(n)
		}
	})

	return nil, nil
}

func (u *unexportedAPIPass) checkFuncDecl(fn *ast.FuncDecl) {
	if !fn.Name.IsExported() {
		return
	}
	u.checkExpr(fn.Type)
}

func (u *unexportedAPIPass) checkTypeSpec(ts *ast.TypeSpec) {
	if !ts.Name.IsExported() {
		return
	}

	if ts.TypeParams != nil {
		for _, param := range ts.TypeParams.List {
			if anyIdentExported(param.Names) {
				if u.checkField(param) {
					return
				}
			}
		}
	}

	u.checkExpr(ts.Type)
}

func (u *unexportedAPIPass) checkValueSpec(vs *ast.ValueSpec) {
	if !anyIdentExported(vs.Names) {
		return
	}

	if vs.Type != nil {
		u.checkExpr(vs.Type)
		return
	}

	for _, value := range vs.Values {
		if u.checkExpr(value) {
			return
		}
	}
}

func anyIdentExported(idents []*ast.Ident) bool {
	for _, ident := range idents {
		if ident.IsExported() {
			return true
		}
	}
	return false
}

func (u *unexportedAPIPass) checkFieldIgnoringNames(field *ast.Field) (stop bool) {
	if field.Type == nil {
		return false
	}

	return u.checkExpr(field.Type)
}

func (u *unexportedAPIPass) checkFieldIfNamesExported(field *ast.Field) (stop bool) {
	if anyIdentExported(field.Names) {
		return u.checkField(field)
	}
	return false
}

func (u *unexportedAPIPass) checkFieldsIgnoringNames(fields *ast.FieldList) (stop bool) {
	if fields == nil {
		return false
	}
	return slices.ContainsFunc(fields.List, u.checkFieldIgnoringNames)
}

func (u *unexportedAPIPass) checkField(field *ast.Field) (stop bool) {
	if field.Type == nil {
		return false
	}

	if u.checkExpr(field.Type) {
		return true
	}

	// TODO

	return false
}

func (u *unexportedAPIPass) checkExpr(expr ast.Expr) (stop bool) {
	if expr == nil {
		return false
	}

	switch expr := expr.(type) {
	case *ast.StructType:
		return slices.ContainsFunc(expr.Fields.List, u.checkFieldIfNamesExported)
	case *ast.StarExpr:
		return u.checkExpr(expr.X)
	case *ast.Ident:
		obj := u.pass.TypesInfo.Defs[expr]
		if obj == nil {
			return false
		}
		if !expr.IsExported() {
			if obj.Parent() == types.Universe {
				return false
			}
			u.pass.Reportf(u.currDecl.Pos(), "exported API %s references unexported identifier %s", u.file.Name.Name, expr.Name)
			return true
		}
		return u.checkType(obj.Type())
	case *ast.MapType:
		if u.checkExpr(expr.Key) {
			return true
		}
		if u.checkExpr(expr.Value) {
			return true
		}
		return false
	case *ast.ArrayType:
		if u.checkExpr(expr.Len) {
			return true
		}
		return u.checkExpr(expr.Elt)
	case *ast.SelectorExpr:
		if !expr.Sel.IsExported() {
			u.pass.Reportf(u.currDecl.Pos(), "exported API %s references unexported identifier %s", u.file.Name.Name, expr.Sel.Name)
			return true
		}
		return false
	case *ast.InterfaceType:
		return slices.ContainsFunc(expr.Methods.List, u.checkFieldIfNamesExported)
	case *ast.ChanType:
		return u.checkExpr(expr.Value)
	case *ast.FuncType:
		if u.checkFieldsIgnoringNames(expr.TypeParams) {
			return true
		}

		if u.checkFieldsIgnoringNames(expr.Params) {
			return true
		}

		if u.checkFieldsIgnoringNames(expr.Results) {
			return true
		}

		return false
	case *ast.Ellipsis:
		return u.checkExpr(expr.Elt)
	case *ast.CompositeLit:
		return u.checkExpr(expr.Type)
	case *ast.IndexListExpr:
		if u.checkExpr(expr.X) {
			return true
		}
		return slices.ContainsFunc(expr.Indices, u.checkExpr)
	case *ast.IndexExpr:
		if u.checkExpr(expr.X) {
			return true
		}
		if u.checkExpr(expr.Index) {
			return true
		}
		return false
	case *ast.UnaryExpr:
		return u.checkExpr(expr.X)
	case *ast.BinaryExpr:
		if u.checkExpr(expr.X) {
			return true
		}
		return u.checkExpr(expr.Y)
	case *ast.BasicLit:
		expr.Obj
	default:
		var buf bytes.Buffer
		format.Node(&buf, u.pass.Fset, expr)
		panic(fmt.Sprintf("%T, unhandled case %T: %s", u.currDecl, expr, buf.String()))
	}
}

func (u *unexportedAPIPass) checkType(typ types.Type) (stop bool) {
	return false
}
