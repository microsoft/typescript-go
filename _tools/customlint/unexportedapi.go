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

// allowEmbeddedUnexportedWithoutExportedMembers controls whether embedded unexported types
// are allowed if they don't expose any exported fields or methods.
// When true: embedded unexported types are OK if they have no exported members
// When false: all embedded unexported types are flagged
const allowEmbeddedUnexportedWithoutExportedMembers = true

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

	// If this is a method on an unexported type, skip it
	// Unexported types and their methods (even if exported) are not part of the public API
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		recvType := fn.Recv.List[0].Type
		// Unwrap pointer receiver if needed
		if star, ok := recvType.(*ast.StarExpr); ok {
			recvType = star.X
		}
		// Check if the receiver type is unexported
		if ident, ok := recvType.(*ast.Ident); ok && !ident.IsExported() {
			return
		}
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

	// If there's an explicit type annotation, check it
	if vs.Type != nil {
		u.checkExpr(vs.Type)
		return
	}

	// If there's no explicit type, we need to check the inferred type, not the initialization expression
	// The initialization expression is an implementation detail and not part of the API
	// For example: var Foo = unexportedFunc() where unexportedFunc returns an exported type is OK
	for _, name := range vs.Names {
		if !name.IsExported() {
			continue
		}
		obj := u.pass.TypesInfo.Defs[name]
		if obj != nil {
			if u.checkType(obj.Type()) {
				return
			}
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
	// For embedded fields (no names), handle specially
	if len(field.Names) == 0 {
		return u.checkEmbeddedField(field)
	}
	// For named fields, only check if at least one name is exported
	if anyIdentExported(field.Names) {
		return u.checkField(field)
	}
	return false
}

func (u *unexportedAPIPass) checkEmbeddedField(field *ast.Field) (stop bool) {
	if field.Type == nil {
		return false
	}

	// Get the type of the embedded field
	typ := u.pass.TypesInfo.TypeOf(field.Type)
	if typ == nil {
		// Fallback to regular checking if we can't get type info
		return u.checkField(field)
	}

	// For embedded fields, don't check the type name itself.
	// Instead, walk through the embedded type's exported members and check those.
	// This way, embedding an unexported type is OK as long as its exported members
	// don't reference other unexported types.

	// Dereference pointers
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	// Check exported fields in structs
	if structType, ok := typ.Underlying().(*types.Struct); ok {
		for field := range structType.Fields() {
			field := field
			if field.Exported() {
				if u.checkType(field.Type()) {
					return true
				}
			}
		}
	}

	// Check exported methods on the type
	if named, ok := typ.(*types.Named); ok {
		for method := range named.Methods() {
			method := method
			if method.Exported() {
				if u.checkType(method.Type()) {
					return true
				}
			}
		}
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
		// First check Defs (for defining occurrences), then Uses (for referring occurrences)
		obj := u.pass.TypesInfo.Defs[expr]
		if obj == nil {
			obj = u.pass.TypesInfo.Uses[expr]
		}
		if obj == nil {
			return false
		}
		if !expr.IsExported() {
			if obj.Parent() == types.Universe {
				return false
			}
			// Only report if the unexported identifier is from the same package
			if obj.Pkg() != nil && obj.Pkg() == u.pass.Pkg {
				u.pass.Reportf(expr.Pos(), "exported API references unexported identifier %s", expr.Name)
				return true
			}
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
		return false
	case *ast.CallExpr:
		// For call expressions, check the function being called
		// We don't check arguments since those are values, not types in the API
		return u.checkExpr(expr.Fun)
	case *ast.FuncLit:
		// Function literals - check the function type
		return u.checkExpr(expr.Type)
	case *ast.ParenExpr:
		return u.checkExpr(expr.X)
	default:
		var buf bytes.Buffer
		format.Node(&buf, u.pass.Fset, expr)
		panic(fmt.Sprintf("%T, unhandled case %T: %s", u.currDecl, expr, buf.String()))
	}
}

func (u *unexportedAPIPass) hasExportedMembers(typ types.Type) bool {
	if typ == nil {
		return false
	}

	// Dereference pointers
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	// Check for exported fields in structs
	if structType, ok := typ.Underlying().(*types.Struct); ok {
		for field := range structType.Fields() {
			if field.Exported() {
				return true
			}
		}
	}

	// Check for exported methods on the type
	if named, ok := typ.(*types.Named); ok {
		for method := range named.Methods() {
			if method.Exported() {
				return true
			}
		}
	}

	return false
}

func (u *unexportedAPIPass) checkType(typ types.Type) (stop bool) {
	if typ == nil {
		return false
	}

	switch typ := typ.(type) {
	case *types.Named:
		// Check if the named type itself is unexported
		obj := typ.Obj()
		if obj != nil && !obj.Exported() && obj.Pkg() == u.pass.Pkg {
			u.pass.Reportf(u.currDecl.Pos(), "exported API references unexported type %s", obj.Name())
			return true
		}
		// Check type arguments if any (for generics)
		if typ.TypeArgs() != nil {
			for t := range typ.TypeArgs().Types() {
				if u.checkType(t) {
					return true
				}
			}
		}
		return false
	case *types.Pointer:
		return u.checkType(typ.Elem())
	case *types.Slice:
		return u.checkType(typ.Elem())
	case *types.Array:
		return u.checkType(typ.Elem())
	case *types.Map:
		if u.checkType(typ.Key()) {
			return true
		}
		return u.checkType(typ.Elem())
	case *types.Chan:
		return u.checkType(typ.Elem())
	case *types.Signature:
		// Check parameters
		if typ.Params() != nil {
			for v := range typ.Params().Variables() {
				if u.checkType(v.Type()) {
					return true
				}
			}
		}
		// Check results
		if typ.Results() != nil {
			for v := range typ.Results().Variables() {
				if u.checkType(v.Type()) {
					return true
				}
			}
		}
		return false
	case *types.Struct:
		// Check all fields
		for field := range typ.Fields() {
			field := field
			// Only check exported fields
			if field.Exported() {
				if u.checkType(field.Type()) {
					return true
				}
			}
		}
		return false
	case *types.Interface:
		// Check all methods
		for method := range typ.Methods() {
			method := method
			// Only check exported methods
			if method.Exported() {
				if u.checkType(method.Type()) {
					return true
				}
			}
		}
		return false
	case *types.Basic, *types.TypeParam:
		// Basic types and type parameters are always OK
		return false
	default:
		// For any unhandled type, be conservative and don't report
		return false
	}
}
