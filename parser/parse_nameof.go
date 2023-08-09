package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func ReplaceNameof(pkg string) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedExportFile |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedModule,
	}
	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		return
	}

	npkgs := make([]*packages.Package, 0, len(pkgs))
	for _, pkg := range pkgs {
		pkg := pkg
		syntax := make([]*ast.File, 0, len(pkg.Syntax))
		for _, file := range pkg.Syntax {
			file := file
			var (
				reachNameof  bool
				nameofCursor *astutil.Cursor
			)
			node := astutil.Apply(file, func(c *astutil.Cursor) bool {
				if reachNameof {
					ast.Print(token.NewFileSet(), c.Node())

					var (
						ne  ast.Expr
						err error
					)
					switch cn := c.Node().(type) {
					case *ast.Ident:
						ne, err = parser.ParseExpr(strconv.Quote(cn.Name))
					case *ast.UnaryExpr:
						if v, ok := cn.X.(*ast.Ident); ok {
							ne, err = parser.ParseExpr(strconv.Quote(v.Name))
						}
						if v, ok := cn.X.(*ast.CompositeLit); ok {
							if vv, ok := v.Type.(*ast.Ident); ok {
								ne, err = parser.ParseExpr(strconv.Quote(vv.Name))
							}
						}
					case *ast.CompositeLit:
						if v, ok := cn.Type.(*ast.Ident); ok {
							ne, err = parser.ParseExpr(strconv.Quote(v.Name))
						}

					case *ast.SelectorExpr:
						ne, err = parser.ParseExpr(strconv.Quote(cn.Sel.Name))
					}
					if err == nil && ne != nil {
						fmt.Println("==== replace before ====")
						ast.Print(token.NewFileSet(), nameofCursor.Node())
						nameofCursor.Replace(ne)
						fmt.Println("==== replace after ====")
					}

					reachNameof = false
					nameofCursor = nil
				}

				ci, ok := c.Node().(*ast.Ident)
				if ok && ci.Name == "nameof" && ci.Obj != nil && ci.Obj.Kind == ast.Fun {
					reachNameof = true
					nameofCursor = c
					// ast.Print(token.NewFileSet(), ci)
				}

				return true
			}, func(c *astutil.Cursor) bool {
				return true
			})
			if err := printer.Fprint(os.Stdout, token.NewFileSet(), node); err != nil {
				panic(err)
			}
			nfile := node.(*ast.File)
			syntax = append(syntax, nfile)
		}
		pkg.Syntax = syntax
		npkgs = append(npkgs, pkg)
	}

	_ = npkgs
	// if err := printer.Fprint(os.Stdout, token.NewFileSet(), npkgs); err != nil {
	// 	panic(err)
	// }
}
