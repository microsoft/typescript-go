package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_go_file.exe <path-to-go-file>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Check if file exists and has .go extension
	if !strings.HasSuffix(filename, ".go") {
		log.Fatalf("File must have .go extension: %s", filename)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", filename)
	}

	// Create a new token file set
	fset := token.NewFileSet()

	// Parse the Go source file
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
	}

	// Modify the AST to add print statements
	modifyFunctions(node)

	writeModifiedFile(fset, node, filename)
	fmt.Printf("Modified file written to: %s\n", filename)
}

type visitor struct{}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	addPrints(&node)
	return v
}

// modifyFunctions adds print statements to all functions in the AST
func modifyFunctions(node *ast.File) {
	// Add fmt import if not already present
	addFmtImport(node)

	v := &visitor{}
	ast.Walk(v, node)
}

// addFmtImport adds the fmt import if it's not already present
func addFmtImport(node *ast.File) {
	// Check if fmt is already imported
	for _, imp := range node.Imports {
		if strings.Trim(imp.Path.Value, `"`) == "fmt" {
			return // fmt already imported
		}
	}

	// Create new fmt import
	fmtImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"fmt"`,
		},
	}

	// Add to imports
	node.Imports = append(node.Imports, fmtImport)

	// // Also add to declarations if needed
	// found := false
	// for _, decl := range node.Decls {
	// 	if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
	// 		genDecl.Specs = append(genDecl.Specs, fmtImport)
	// 		found = true
	// 		break
	// 	}
	// }

	// // If no import declaration exists, create one
	// if !found {
	// 	importDecl := &ast.GenDecl{
	// 		Tok:   token.IMPORT,
	// 		Specs: []ast.Spec{fmtImport},
	// 	}
	// 	// Insert at the beginning of declarations (after package)
	// 	newDecls := make([]ast.Decl, 0, len(node.Decls)+1)
	// 	newDecls = append(newDecls, importDecl)
	// 	newDecls = append(newDecls, node.Decls...)
	// 	node.Decls = newDecls
	// }
}

func addPrints(n *ast.Node) {
	switch n := (*n).(type) {
	case *ast.FuncDecl:
		addPrintBody(n.Body, n.Name.Name)
	case *ast.AssignStmt:
		for i, expr := range n.Rhs {
			lhs := n.Lhs[i]
			if ident, ok := lhs.(*ast.Ident); ok {
				if fnLit, ok := expr.(*ast.FuncLit); ok {
					addPrintBody(fnLit.Body, ident.Name)
				}
			}
		}
	default:
		// do nothing
	}
}

func addPrintBody(fnBody *ast.BlockStmt, fnName string) {
	if fnBody == nil || len(fnBody.List) == 0 {
		return
	}

	// Insert entering print statement at the beginning of the function body
	newStmts := make([]ast.Stmt, 0, len(fnBody.List)+2)
	enterStmt := createPrintStmt(fnName, true /*enter*/)
	exitStmt := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.FuncLit{
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						createPrintStmt(fnName, false /*enter*/),
					},
				},
			},
		},
	}
	newStmts = append(newStmts, enterStmt)
	newStmts = append(newStmts, exitStmt)
	newStmts = append(newStmts, fnBody.List...)
	fnBody.List = newStmts
}

// createPrintStmt creates a print statement for function enter/exit
func createPrintStmt(fnName string, enter bool) *ast.ExprStmt {
	var printStr string
	if enter {
		printStr = `"> %s\n"`
	} else {
		printStr = `"< %s\n"`
	}
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "fmt"},
				Sel: &ast.Ident{Name: "Printf"},
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: printStr,
				},
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"` + fnName + `"`,
				},
			},
		},
	}
}

// writeModifiedFile writes the modified AST to a new file
func writeModifiedFile(fset *token.FileSet, node *ast.File, filename string) {
	// For debugging: print the modified AST to stdout
	// printer.Fprint(os.Stdout, fset, node)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer file.Close()

	// Format and write the modified AST
	if err := format.Node(file, fset, node); err != nil {
		log.Fatalf("Error formatting code: %v", err)
	}
}
