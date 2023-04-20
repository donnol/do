package parser

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
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
			packages.NeedExportFile |
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
