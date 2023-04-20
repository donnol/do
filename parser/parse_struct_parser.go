package parser

import (
	"bytes"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

// 为什么要这样写，因为：https://github.com/golang/go/issues/27477
type (
	// Parser 解析器
	// 解析指定的包导入路径，获取go源码信息
	Parser struct {
		filter func(os.FileInfo) bool // 过滤器

		fset *token.FileSet

		useSourceImporter bool // 使用源码importer

		replaceImportPath bool // 替换导入路径
		fromPath          string
		toPath            string
		output            io.Writer

		op Op // 操作，如生成接口，生成实现等

		replaceCallExpr bool

		PkgInfo
	}
)

func NewParser(opt Option) *Parser {
	return &Parser{
		op:                opt.Op,
		filter:            opt.Filter,
		useSourceImporter: opt.UseSourceImporter,
		replaceImportPath: opt.ReplaceImportPath,
		fromPath:          opt.FromPath,
		toPath:            opt.ToPath,
		output:            opt.Output,
		replaceCallExpr:   opt.ReplaceCallExpr,
	}
}

func (p *Parser) GetPkgInfo() PkgInfo {
	return p.PkgInfo
}

// ParseByGoPackages 使用x/tools/go/packages解析指定导入路径
func (p *Parser) ParseByGoPackages(patterns ...string) (result Packages, err error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedExportsFile |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedModule,
	}
	// pattern可以是文件目录，也可以是包导入路径，如：'~/a/b/c', 'bytes', 'github.com/donnol/tools'...
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return
	}

	result.Patterns = patterns
	result.Pkgs = make([]Package, 0, len(pkgs))
	inspector := NewInspector(InspectOption{
		Parser: p,
	})
	for _, pkg := range pkgs {
		p.fset = pkg.Fset
		tmpPkg := inspector.InspectPkg(pkg)

		result.Pkgs = append(result.Pkgs, tmpPkg)
	}

	return
}

func (p *Parser) GetStandardPackages() []string {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	standardPackages := make([]string, 0, len(pkgs))
	for _, p := range pkgs {
		standardPackages = append(standardPackages, p.PkgPath)
	}

	return standardPackages
}

func (p *Parser) parseDir(fset *token.FileSet, fullDir string) (pkgs map[string]*ast.Package, err error) {
	const (
		testSuffix = "_test"
	)

	// 解析目录
	pkgs, err = parser.ParseDir(fset, fullDir, func(fi os.FileInfo) bool {
		li := strings.LastIndex(fi.Name(), filepath.Ext(fi.Name()))

		// 跳过test文件
		testi := strings.LastIndex(fi.Name(), testSuffix)
		if testi != -1 && li-testi == len(testSuffix) {
			return false
		}

		return true
	}, parser.ParseComments)
	if err != nil {
		return
	}

	return
}

func (p *Parser) typesCheck(path string, files []*ast.File) (info *types.Info, err error) {
	imp := importer.Default()
	if p.useSourceImporter {
		fset := token.NewFileSet()
		imp = importer.ForCompiler(fset, "source", nil)
	}

	// 获取类型信息
	conf := types.Config{
		IgnoreFuncBodies: true,

		// 默认是用go install安装后生成的.a文件，可以选择使用source，但是会慢很多
		Importer: imp,

		Error: func(err error) {
			log.Printf("Check Failed: %+v\n", err)
		},
		DisableUnusedImportCheck: true,
	}
	info = &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	// conf.Check的path参数传入的是包名，而不是导入路径
	pkg, err := conf.Check(path, p.fset, files, info)
	if err != nil {
		return
	}
	p.methodSet(pkg)

	return
}

// 根据types.Type找到method set，但是怎么将它转为interface呢？
func (p *Parser) methodSet(pkg *types.Package) {
	if pkg.Scope() == nil {
		return
	}
	for _, name := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(name)
		if obj == nil {
			continue
		}
		typ := obj.Type()
		for _, t := range []types.Type{typ, types.NewPointer(typ)} {
			// fmt.Printf("Method set of %s:\n", t)
			mset := types.NewMethodSet(t)
			for i := 0; i < mset.Len(); i++ {
				sel := mset.At(i)
				_ = sel
				// fmt.Println("sel: ", sel, "type:", sel.Type(), reflect.TypeOf(sel.Type().Underlying()), "obj:", sel.Obj())
			}
			// fmt.Println()
		}
	}
}

func (p *Parser) replaceFileImportPath(fileName string, file *ast.File) error {
	var err error

	// 替换import path
	for _, fi := range file.Imports {
		path := strings.Trim(fi.Path.Value, `"`)

		if strings.HasPrefix(path, p.fromPath) {
			topath := strings.Replace(path, p.fromPath, p.toPath, 1)

			rewrote := astutil.RewriteImport(p.fset, file, path, topath)

			log.Printf("From %s to %s, rewrote: %v\n", p.fromPath, p.toPath, rewrote)
		}
	}

	// 获取file的ast内容并格式化
	buf := bytes.NewBuffer([]byte{})
	printer.Fprint(buf, p.fset, file)
	content, err := Format(fileName, buf.String(), true)
	if err != nil {
		return err
	}

	// 将内容输出到原文件
	output, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = output.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}
