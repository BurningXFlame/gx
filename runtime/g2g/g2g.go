/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

const pkg = "github.com/burningxflame/gx/runtime/gls"

type g2G struct {
	replaced  bool
	lastFor   *ast.ForStmt
	lastRange *ast.RangeStmt
}

func (g *g2G) clear() {
	g.replaced = false
	g.lastFor = nil
	g.lastRange = nil
}

func (g *g2G) process(pa string) error {
	fset := token.NewFileSet()
	tree, err := parser.ParseFile(fset, pa, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	_ = astutil.Apply(tree, g.apply, nil)
	if !g.replaced {
		return nil
	}

	astutil.AddImport(fset, tree, pkg)

	out, err := os.OpenFile(pa, os.O_RDWR|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer out.Close()

	err = format.Node(out, fset, tree)
	if err != nil {
		return err
	}

	return nil
}

func (g *g2G) apply(cur *astutil.Cursor) bool {
	var goSt *ast.GoStmt

	switch node := cur.Node().(type) {
	case *ast.GoStmt:
		goSt = node
	case *ast.ForStmt:
		g.lastFor = node
		return true
	case *ast.RangeStmt:
		g.lastRange = node
		return true
	default:
		return true
	}

	if len(goSt.Call.Args) > 0 {
		if g.lastFor != nil && g.lastFor.Body == cur.Parent() {
			forClosure(cur, g.lastFor)
		} else if g.lastRange != nil && g.lastRange.Body == cur.Parent() {
			rangeClosure(cur, g.lastRange)
		}
	}

	var expr ast.Expr
	if fn, ok := goSt.Call.Fun.(*ast.FuncLit); ok && len(fn.Type.Params.List) == 0 {
		expr = goSt.Call.Fun
	} else {
		expr = &ast.FuncLit{
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: goSt.Call,
					},
				},
			},
		}
	}

	cur.Replace(&ast.ExprStmt{
		X: &ast.CallExpr{
			Fun:  ast.NewIdent("gls.Go"),
			Args: []ast.Expr{expr},
		},
	})

	g.replaced = true
	return true
}

func forClosure(cur *astutil.Cursor, forSt *ast.ForStmt) {
	var le []ast.Expr

	switch st := forSt.Post.(type) {
	case *ast.IncDecStmt:
		if id, ok := st.X.(*ast.Ident); ok {
			le = append(le, ast.NewIdent(id.Name))
		}

	case *ast.AssignStmt:
		for _, e := range st.Lhs {
			if id, ok := e.(*ast.Ident); ok {
				le = append(le, ast.NewIdent(id.Name))
			}
		}

	default:
		return
	}

	if len(le) == 0 {
		return
	}

	cur.InsertBefore(&ast.AssignStmt{
		Lhs: le,
		Tok: token.DEFINE,
		Rhs: le,
	})
}

func rangeClosure(cur *astutil.Cursor, rangeSt *ast.RangeStmt) {
	var le []ast.Expr

	if id, ok := rangeSt.Key.(*ast.Ident); ok && id.Name != "_" {
		le = append(le, ast.NewIdent(id.Name))
	}

	if id, ok := rangeSt.Value.(*ast.Ident); ok && id.Name != "_" {
		le = append(le, ast.NewIdent(id.Name))
	}

	if len(le) == 0 {
		return
	}

	cur.InsertBefore(&ast.AssignStmt{
		Lhs: le,
		Tok: token.DEFINE,
		Rhs: le,
	})
}
