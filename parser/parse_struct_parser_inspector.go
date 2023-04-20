package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/go/packages"
)

type Packages struct {
	Patterns []string
	Pkgs     []Package
}

func (pkgs Packages) LookupPkg(name string) (Package, bool) {
	pkg := Package{}
	for _, single := range pkgs.Pkgs {
		if single.Name == name {
			return single, true
		}
	}
	return pkg, false
}

func (pkg Package) NewGoFileWithSuffix(suffix string) (file string) {
	if pkg.Module == nil {
		fmt.Printf("pkg.Module is nil\n")
	}
	part := strings.ReplaceAll(pkg.PkgPath, pkg.Module.Path, "")
	log.Printf("pkg: %+v, module: %+v, %s\n", pkg.PkgPath, pkg.Module, part)

	dir := filepath.Join(pkg.Module.Dir, part)
	file = filepath.Join(dir, suffix+".go")

	return
}

func (pkg Package) SaveInterface(file string) error {
	var gocontent = "package " + pkg.Name + "\n"

	var content string
	for _, single := range pkg.Structs {
		content += single.MakeInterface() + "\n\n"
	}
	if content == "" {
		return nil
	}
	gocontent += content

	// TODO:检查是否重复

	if file == "" {
		file = pkg.NewGoFileWithSuffix("interface")
	}
	// 写入
	formatContent, err := Format(file, gocontent, false)
	if err != nil {
		return err
	}
	log.Printf("content: %s, file: %s\n", formatContent, file)

	if err = ioutil.WriteFile(file, []byte(formatContent), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (pkg Package) SaveMock(file string) error {
	var gocontent = "package " + pkg.Name + "\n"

	// 找出所有外部包引用，生成import
	// 因为是生成mock结构体，所以有包引用的都是参数和返回值
	imports := make(map[string]struct{}, 4)

	// debug.Printf("===test\n")
	var content string
	for _, single := range pkg.Interfaces {
		// debug.Printf("have type set: %+v, embeds: %d\n", single.Interface, single.Interface.NumEmbeddeds())
		if single.Interface.NumEmbeddeds() != 0 {
			log.Printf("have type set: %+v\n", single.Interface)
			continue
		}
		mock, imps := single.MakeMock()
		for imp := range imps {
			imports[imp] = struct{}{}
		}
		content += mock + "\n\n"
	}
	if content == "" {
		return nil
	}

	// 全局变量/函数
	globalContent := `
	var (
		_gen_customCtxMap = make(map[string]inject.CtxFunc)
	)

	func RegisterProxyMethod(pctx inject.ProxyContext, cf inject.CtxFunc) {
		_gen_customCtxMap[pctx.Uniq()] = cf
	}
	`
	content = globalContent + content

	// 导入
	var impcontent string
	for imp := range imports {
		if imp == "" {
			continue
		}
		impcontent += `"` + imp + `"` + "\n"
	}
	if impcontent != "" {
		impcontent = "import (\n" + impcontent + ")\n"
		// debug.Printf("import: %s\n", impcontent)
	}

	gocontent += impcontent

	// mock
	gocontent += content
	// debug.Printf("gocontent: %s\n", gocontent)

	// TODO:检查是否重复

	if file == "" {
		file = pkg.NewGoFileWithSuffix("mock")
	}
	// 写入
	formatContent, err := Format(file, gocontent, false)
	if err != nil {
		return fmt.Errorf("format failed: %w, content: \n%s", err, gocontent)
	}
	// debug.Printf("content: %s, file: %s\n", formatContent, file)

	if err = ioutil.WriteFile(file, []byte(formatContent), os.ModePerm); err != nil {
		return err
	}

	return nil
}

type Package struct {
	*packages.Package

	Funcs      []Func
	Structs    []Struct
	Interfaces []Interface
}

type Interface struct {
	*types.Interface

	PkgPath string
	PkgName string
	Name    string
	Methods []Method // 方法列表
}

var (
	proxyMethodTmpl = `
	{{.methodName}}: {{.funcSignature}} {
		_gen_begin := time.Now()

		{{.funcResult}}

		_gen_ctx := {{.mockType}}{{.funcName}}ProxyContext
		_gen_cf, _gen_ok := _gen_customCtxMap[_gen_ctx.Uniq()]
		if _gen_ok {
			_gen_params := []any{}
			{{.params}}
			_gen_res := _gen_cf(_gen_ctx, base.{{.funcName}}, _gen_params)
			{{.resultAssert}}
		} else {
			{{.funcResultList}} = base.{{.funcName}}({{.argNames}})
		}

		log.Printf("[ctx: %s]used time: %v\n", _gen_ctx.Uniq(), time.Since(_gen_begin))

		return {{.funcResultList}}
	},
	`

	proxyMethodParamsTmpl = `
	{{range $index, $ele := .args}}
		_gen_params = append(_gen_params, {{.Name}})
	{{end}}
	`
	proxyMethodResultTmpl = `
	 {{range $index, $ele := .reses}}
	 	var _gen_r{{$index}} {{.Typ}}
	 {{end}}
	`

	proxyMethodResultAssertTmpl = `
	{{range $index, $ele := .reses}}
			_gen_tmpr{{$index}}, _gen_exist := _gen_res[{{$index}}].({{.Typ}})
			if _gen_exist {
				_gen_r{{$index}} = _gen_tmpr{{$index}}
			}
	{{end}}
	`
)

type arg struct {
	Name     string
	Typ      string
	Variadic bool
}

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func (s Interface) MakeMock() (string, map[string]struct{}) {
	mockType := s.makeMockName()
	mockRecv := s.makeMockRecv()
	proxyFuncName := s.makeProxyFuncName()
	// debug.Printf("proxyfuncname:%s\n", proxyFuncName)

	proxyFunc := "func " + proxyFuncName + "(base " + s.Name + ") *" + mockType + "{" + `if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	return &` + mockType + `{`
	cc := fmt.Sprintf(`%sCommonProxyContext`, LcFirst(mockType))

	var is string
	var pc string
	var ms string
	var proxyMethod = new(bytes.Buffer)
	var imports = make(map[string]struct{}, 4)
	for _, m := range s.Methods {
		fieldName, fieldType, methodSig, returnStmt, call, args, reses, imps := s.processFunc(m)

		for imp := range imps {
			imports[imp] = struct{}{}
		}

		is += fmt.Sprintf("\n%s %s\n", fieldName, fieldType)

		pc += fmt.Sprintf(`%s%sProxyContext = func() (pctx inject.ProxyContext) { 
			pctx = %s
			pctx.MethodName = "%s"
			return
		} () 
		`, mockType, m.Name, cc, m.Name)

		ms += fmt.Sprintf("\nfunc (%s *%s) %s {\n %s %s.%s \n}\n", mockRecv, mockType, methodSig, returnStmt, mockRecv, call)

		assertBuf := new(bytes.Buffer)
		assertTmpl, err := template.New("proxyMethodResultAssert").Parse(proxyMethodResultAssertTmpl)
		if err != nil {
			panic(err)
		}
		assertTmpl.Execute(assertBuf, map[string]interface{}{
			"reses": reses,
		})
		paramBuf := new(bytes.Buffer)
		paramTmpl, err := template.New("proxyMethodParam").Parse(proxyMethodParamsTmpl)
		if err != nil {
			panic(err)
		}
		paramTmpl.Execute(paramBuf, map[string]interface{}{
			"args": args,
		})
		argNames := ""
		for i, arg := range args {
			argNames += arg.Name
			if i != len(args)-1 {
				argNames += ", "
			}
		}
		resBuf := new(bytes.Buffer)
		resTmpl, err := template.New("proxyMethodResult").Parse(proxyMethodResultTmpl)
		if err != nil {
			panic(err)
		}
		resTmpl.Execute(resBuf, map[string]interface{}{
			"reses": reses,
		})
		funcResultList := ""
		for i := range reses {
			funcResultList += "_gen_r" + strconv.Itoa(i)
			if i != len(reses)-1 {
				funcResultList += ", "
			}
		}
		tmpl, err := template.New("proxyMethod").Parse(proxyMethodTmpl)
		if err != nil {
			panic(err)
		}
		tmpl.Execute(proxyMethod, map[string]interface{}{
			"methodName":     fieldName,
			"funcSignature":  strings.Replace(methodSig, m.Name, "func", 1),
			"mockType":       mockType,
			"funcName":       m.Name,
			"funcResult":     resBuf.String(),
			"funcResultList": funcResultList,
			"argNames":       argNames,
			"params":         paramBuf.String(),
			"resultAssert":   assertBuf.String(),
		})
	}

	proxyFunc += proxyMethod.String() + "}}"
	is = mockPrefix(mockType, is)

	is += `var (_ ` + s.Name + ` = &` + mockType + "{}\n\n"
	is += fmt.Sprintf(`%s = inject.ProxyContext {
		PkgPath: "%s",
		InterfaceName: "%s",
	}
	`, cc, s.PkgPath, s.Name)
	is += pc + "\n_ =" + proxyFuncName + `)`
	is += "\n" + "\n\n" + proxyFunc + "\n"
	is += ms

	// debug.Printf("is: %s\n", is)

	return is, imports
}

const (
	sep         = ","
	leftParent  = "("
	rightParent = ")"
)

type qualifierParam struct {
	pkgPath string
}

var (
	// 包名
	// 包名有可能不等于包路径的最后一部分的（最后一个'/'后面的部分）
	pkgNameQualifier = func(qp qualifierParam) types.Qualifier {
		return func(pkg *types.Package) string {
			name := pkg.Name()

			// 如果是同一个包内的，省略包名
			if pkg.Path() == qp.pkgPath {
				return ""
			}

			return name
		}
	}
)

// func(ctx context.Context, m M) (err error) -> (ctx, m)
// func(context.Context,M) (error) -> (p0, p1)
func (s Interface) processFunc(m Method) (fieldName, fieldType, methodSig, returnStmt, call string, args []arg, reses []arg, imports map[string]struct{}) {

	imports = make(map[string]struct{}, 4) // 导入的包
	fieldName = m.Name + "Func"
	fieldType = m.Signature

	sigType := m.Origin.Type().(*types.Signature)
	if sigType.Variadic() {
		//  在这里获取完整签名字符串时，还是正常的：func(interface{}, string, ...interface{}) error
		typStr := types.TypeString(sigType, pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath}))
		_ = typStr
		// debug.Printf("typ: %+v, str: %s\n", sigType, typStr)
	}
	params := sigType.Params()
	for i := 0; i < params.Len(); i++ {
		pvar := params.At(i)
		name := pvar.Name()

		// 参数名可能为空，需要置默认值
		if name == "" {
			name = fmt.Sprintf("p%d", i)
		}

		// 参数类型的包路径信息
		pkgPath := getTypesPkgPath(pvar.Type())
		imports[pkgPath] = struct{}{}

		// 解析进来之后，不定参数类型变成了slice：[]interface{}
		typStr := types.TypeString(pvar.Type(), pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath}))

		// 处理最后一个是不定参数的情况
		var paramTypePrefix string
		var variadic bool
		if sigType.Variadic() && i == params.Len()-1 {
			paramTypePrefix = "..."
			variadic = true
			// debug.Printf("typ: %+v, str: %s, params: %v\n", pvar.Type(), typStr, params.String())
		}

		// FIXME:感觉不太好，怎么办呢？
		// 当是不定参数，typStr会从...interface{}变为[]interface{}，因此，需要再将它重新变回来
		if paramTypePrefix != "" && strings.Index(typStr, "[]") == 0 {
			typStr = typStr[2:]
		}
		methodSig += name + " " + paramTypePrefix + typStr + sep

		call += name + paramTypePrefix + sep

		args = append(args, arg{Name: name + paramTypePrefix, Typ: typStr, Variadic: variadic})
	}
	methodSig = strings.TrimRight(methodSig, sep)
	methodSig = m.Name + leftParent + methodSig + rightParent

	res := sigType.Results()
	returnStmt = "return"
	if res.Len() == 0 {
		returnStmt = " "
	}
	var resString string
	for i := 0; i < res.Len(); i++ {
		rvar := res.At(i)
		name := rvar.Name()

		// 返回类型的包路径信息
		pkgPath := getTypesPkgPath(rvar.Type())
		imports[pkgPath] = struct{}{}

		typ := types.TypeString(rvar.Type(), pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath}))
		resString += name + " " + typ + sep

		reses = append(reses, arg{
			Name: name,
			Typ:  typ,
		})
	}
	resString = strings.TrimRight(resString, sep)
	resString = leftParent + resString + rightParent
	methodSig = methodSig + resString

	// debug.Printf("methodSig: %v\n", methodSig)

	call = strings.TrimRight(call, sep)
	call = leftParent + call + rightParent
	call = fieldName + call

	return
}

func (s Interface) makeProxyFuncName() string {
	return "get" + s.Name + "Proxy"
}

func (s Interface) makeMockName() string {
	name := s.removeI()
	return name + "Mock"
}

func (s Interface) removeI() string {
	name := s.Name
	// 如果首个字符是I，则去掉
	index := strings.Index(name, "I")
	if index == 0 {
		name = name[1:]
	}
	return name
}

func (s Interface) makeMockRecv() string {
	return "mockRecv"
}

func mockPrefix(name, is string) string {
	return "type " + name + " struct{ " + is + "}\n"
}

// Struct 结构体
type Struct struct {
	// 如：github.com/pkg/errors
	PkgPath string `json:"pkgPath" toml:"pkg_path"` // 包路径

	// 如: errors
	PkgName string // 包名

	Field

	Fields  []Field  // 字段列表
	Methods []Method // 方法列表
}

// --- 测试方法

// 让它传入本包里的另外一个结构体
// 传入本项目其它包的结构体
func (s Struct) String(f Field, ip ImportPath) {
	fmt.Printf("%s\n", s.PkgName)
}

func (s Struct) TypeAlias(p IIIIIIIInfo, ip ImportPathAlias) {

}

func (s Struct) Demo(in types.Array) types.Basic {
	return types.Basic{}
}

func (s *Struct) PointerMethod(in types.Basic) types.Slice {
	return types.Slice{}
}

// --- 测试方法

// MMakeInterface 根据结构体的方法生成相应接口
func (s Struct) MakeInterface() string {
	methods := make([]*types.Func, 0, len(s.Methods))
	for _, m := range s.Methods {
		if !m.Origin.Exported() {
			continue
		}
		methods = append(methods, m.Origin)
		// fmt.Printf("method: %+v, %s\n", m.Origin, m.Signature)
	}

	if len(methods) == 0 {
		return ""
	}

	i := types.NewInterfaceType(methods, nil)
	i = i.Complete()
	is := types.TypeString(i.Underlying(), pkgNameQualifier(qualifierParam{pkgPath: s.PkgPath}))

	is = interfacePrefix(s.makeInterfaceName(), is)

	return is
}

func (s Struct) makeInterfaceName() string {
	// 考虑到结构体名称是非导出有后缀的，如：fileImpl
	// 1. 针对非导出，将首字母变大写
	// 2. 针对impl后缀，直接去掉
	name := s.Name
	name = cases.Title(language.SimplifiedChinese).String(name)
	index := strings.Index(name, "Impl")
	if index != -1 {
		name = name[:index]
	}
	return "I" + name
}

func interfacePrefix(name, is string) string {
	return "type " + name + " " + is
}

type Method = Func

type Func struct {
	Origin *types.Func

	PkgPath   string // 包路径
	Recv      string // 方法的receiver
	Name      string
	Signature string

	Calls []Func // 调用的函数/方法
}

func (f *Func) Set(fm map[string]Func, depth int) {
	l := 1
	setLowerCalls(f.Calls, fm, l, depth)
}

func setLowerCalls(calls []Func, fm map[string]Func, l, depth int) {
	if l > depth {
		return
	}
	for i, call := range calls {
		var key = call.Name
		if call.Recv != "" {
			key = call.Recv + "." + call.Name
		}
		if len(call.Calls) == 0 {
			calls[i].Calls = fm[key].Calls
			nl := l + 1
			setLowerCalls(calls[i].Calls, fm, nl, depth)
		}
	}
}

// PrintCallGraph 打印调用图，用ignore忽略包，用depth指定深度
func (f Func) PrintCallGraph(ignore []string, depth int) {
	ip := &ImportPath{}
	curPath, err := ip.GetByCurrentDir()
	if err != nil {
		panic(err) // 怎么知道这些内置函数是内置函数呢？
	}
	fmt.Printf("root module path: %s\n", curPath)

	fmt.Printf("root: %s(%s)\n", f.Name, f.PkgPath)
	l := 1

	printCallGraph(f.Calls, ignore, l, depth)
}

func printCallGraph(calls []Func, ignores []string, l, depth int) {
	for _, one := range calls {
		if l > depth {
			break
		}

		// 判断是否需要跳过
		pkgPath := one.PkgPath
		needIgnore := false
		for _, ignore := range ignores {
			if pkgPath != "" && ignore == pkgPath {
				needIgnore = true
				break
			}
		}
		if needIgnore {
			continue
		}

		fmt.Printf("%s -> %s(%s)\n", getIdent(l), one.Name, one.PkgPath)

		if len(one.Calls) > 0 {
			nl := l + 1
			printCallGraph(one.Calls, ignores, nl, depth)
		}
	}
}

const (
	ident = "	"
)

func getIdent(l int) string {
	s := ""
	for i := 0; i < l; i++ {
		if i == l-1 {
			s += "   " + strconv.Itoa(l)
		} else {
			s += ident
		}
	}
	return s
}

// Field 字段
type Field struct {
	Id        string // 唯一标志
	Name      string // 名称
	Anonymous bool   // 是否匿名

	TypesType types.Type // 原始类型
	Type      string     // 类型，包含包导入路径

	Tag         string        `json:"tag"` // 结构体字段的tag
	TagBasicLit *ast.BasicLit // ast的tag类型

	Doc     string // 文档
	Comment string // 注释
}

// IIIIIIIInfo 别名测试
type IIIIIIIInfo = Field // 别名测试注释

type ImportPathAlias = ImportPath

type PkgInfo struct {
	dir     string
	pkgName string
}

func (i PkgInfo) GetDir() string {
	return i.dir
}

func (i PkgInfo) GetPkgName() string {
	return i.pkgName
}

type FileResult struct {
	structMap    map[string]Struct    // 名称 -> 结构体
	methodMap    map[string][]Method  // 名称 -> 方法列表
	interfaceMap map[string]Interface // 名称 -> 接口
	funcMap      map[string]Func      // 名称 -> 方法
}

func MakeFileResult() FileResult {
	return FileResult{
		structMap:    make(map[string]Struct),
		methodMap:    make(map[string][]Method),
		interfaceMap: make(map[string]Interface),
		funcMap:      make(map[string]Func),
	}
}

type DeclResult struct {
	structMap    map[string]Struct
	methodMap    map[string][]Method
	interfaceMap map[string]Interface // 名称 -> 接口
	funcMap      map[string]Func      // 名称 -> 方法
}

func MakeDeclResult() DeclResult {
	return DeclResult{
		structMap:    make(map[string]Struct),
		methodMap:    make(map[string][]Method),
		interfaceMap: make(map[string]Interface),
		funcMap:      make(map[string]Func),
	}
}

type SpecResult struct {
	structMap    map[string]Struct    // 名称 -> 结构体
	interfaceMap map[string]Interface // 名称 -> 接口
	funcMap      map[string]Func      // 名称 -> 方法
}

func MakeSpecResult() SpecResult {
	return SpecResult{
		structMap:    make(map[string]Struct),
		interfaceMap: make(map[string]Interface),
		funcMap:      make(map[string]Func),
	}
}

type ExprResult struct {
	Fields  []Field
	pkgPath string
	funcMap map[string]Func // 名称 -> 方法
}

func MakeExprResult() ExprResult {
	return ExprResult{
		Fields:  make([]Field, 0),
		funcMap: make(map[string]Func),
	}
}

func (er ExprResult) Merge(oer ExprResult) (ner ExprResult) {
	ner = er

	if ner.pkgPath == "" && oer.pkgPath != "" {
		ner.pkgPath = oer.pkgPath
	}
	ner.Fields = append(ner.Fields, oer.Fields...)
	for k, v := range oer.funcMap {
		ner.funcMap[k] = v
	}

	return
}

type StmtResult struct {
	pkgPath string
	funcMap map[string]Func // 名称 -> 方法
}

func MakeStmtResult() StmtResult {
	return StmtResult{
		funcMap: make(map[string]Func),
	}
}

func (er StmtResult) Merge(oer StmtResult) (ner StmtResult) {
	ner = er

	if ner.pkgPath == "" && oer.pkgPath != "" {
		ner.pkgPath = oer.pkgPath
	}
	for k, v := range oer.funcMap {
		ner.funcMap[k] = v
	}

	return
}

func (er StmtResult) MergeExprResult(oer ExprResult) (ner StmtResult) {
	ner = er

	if ner.pkgPath == "" && oer.pkgPath != "" {
		ner.pkgPath = oer.pkgPath
	}
	for k, v := range oer.funcMap {
		ner.funcMap[k] = v
	}

	return
}

type FieldResult struct {
	RecvName string

	Fields []Field
}

func MakeFieldResult() FieldResult {
	return FieldResult{
		Fields: make([]Field, 0),
	}
}

type TokenResult struct {
}

type Inspector struct {
	parser *Parser

	pkg *packages.Package
}

type InspectOption struct {
	Parser *Parser
}

func NewInspector(opt InspectOption) *Inspector {
	return &Inspector{
		parser: opt.Parser,
	}
}

func (ins *Inspector) InspectPkg(pkg *packages.Package) Package {
	if pkg == nil {
		panic("input pkg is nil")
	}
	ins.pkg = pkg

	// 解析*ast.File信息
	structMap := make(map[string]Struct)
	methodsMap := make(map[string][]Method)
	interfaceMap := make(map[string]Interface)
	funcMap := make(map[string]Func)
	for i, astFile := range pkg.Syntax {
		// 替换import path
		if ins.parser.replaceImportPath {
			fileName := pkg.CompiledGoFiles[i]
			// debug.Printf("%v\n", pkg.CompiledGoFiles)
			if err := ins.parser.replaceFileImportPath(fileName, astFile); err != nil {
				panic(fmt.Errorf("replaceFileImportPath failed: %+v", err))
			}
			continue
		}

		fileResult := ins.InspectFile(astFile)

		for k, v := range fileResult.structMap {
			structMap[k] = v
		}
		for k, v := range fileResult.methodMap {
			methodsMap[k] = append(methodsMap[k], v...)
		}
		for k, v := range fileResult.interfaceMap {
			interfaceMap[k] = v
		}
		for k, v := range fileResult.funcMap {
			funcMap[k] = v
		}
	}

	structNames := make([]string, 0, len(structMap))
	for structName := range structMap {
		structNames = append(structNames, structName)
	}
	sort.Slice(structNames, func(i, j int) bool {
		return structNames[i] < structNames[j]
	})
	structs := make([]Struct, 0, len(structMap))
	for _, structName := range structNames {
		single := structMap[structName]
		methods := methodsMap[structName]
		single.Methods = methods
		structs = append(structs, single)
	}

	interNames := make([]string, 0, len(interfaceMap))
	for interName := range interfaceMap {
		interNames = append(interNames, interName)
	}
	sort.Slice(interNames, func(i, j int) bool {
		return interNames[i] < interNames[j]
	})
	inters := make([]Interface, 0, len(interfaceMap))
	for _, interName := range interNames {
		single := interfaceMap[interName]
		inters = append(inters, single)
	}

	funcNames := make([]string, 0, len(funcMap))
	for funcName := range funcMap {
		funcNames = append(funcNames, funcName)
	}
	sort.Slice(funcNames, func(i, j int) bool {
		return funcNames[i] < funcNames[j]
	})
	funcs := make([]Func, 0, len(funcMap))
	for _, funcName := range funcNames {
		single := funcMap[funcName]
		funcs = append(funcs, single)
	}

	return Package{
		Package:    pkg,
		Structs:    structs,
		Interfaces: inters,
		Funcs:      funcs,
	}
}

func (ins *Inspector) InspectFile(file *ast.File) (result FileResult) {
	if file == nil {
		return
	}
	result = MakeFileResult()

	structMap := make(map[string]Struct)
	methodsMap := make(map[string][]Method)
	interfaceMap := make(map[string]Interface)
	funcMap := make(map[string]Func)
	for _, decl := range file.Decls {
		declResult := ins.inspectDecl(decl, "")
		for k, v := range declResult.structMap {
			structMap[k] = v
		}
		for k, v := range declResult.methodMap {
			methodsMap[k] = append(methodsMap[k], v...)
		}
		for k, v := range declResult.interfaceMap {
			interfaceMap[k] = v
		}
		for k, v := range declResult.funcMap {
			funcMap[k] = v
		}
	}
	result.structMap = structMap
	result.methodMap = methodsMap
	result.interfaceMap = interfaceMap
	result.funcMap = funcMap

	return
}

func (ins *Inspector) inspectDecl(decl ast.Decl, from string) (result DeclResult) {
	if decl == nil {
		return
	}
	result = MakeDeclResult()

	switch declValue := decl.(type) {
	case *ast.BadDecl:
		panic(fmt.Errorf("BadDecl: %+v", declValue))

	case *ast.FuncDecl:
		// spew.Dump(declValue)
		// debug.Printf("FundDecl name: %s, %s\n", declValue.Name, declValue.Doc.Text())

		funcType := &types.Func{}
		obj := ins.pkg.TypesInfo.Defs[declValue.Name]
		switch objTyp := obj.Type().(type) {
		case *types.Signature:
			// debug.Printf("objTyp sig: %+v, %s\n", objTyp, toString(objTyp))
			funcType = types.NewFunc(declValue.Type.Func, ins.pkg.Types, obj.Name(), objTyp)
		}
		method := Method{
			Origin:    funcType,
			PkgPath:   obj.Pkg().Path(),
			Name:      obj.Name(),
			Signature: toString(obj.Type()),
		}
		from = method.Name

		ins.inspectExpr(declValue.Type, from)               // 函数签名
		stmtResult := ins.inspectStmt(declValue.Body, from) // 函数体
		for _, oneFunc := range stmtResult.funcMap {
			method.Calls = append(method.Calls, oneFunc)
		}

		// debug.Printf(from+"method: %+v\n", method)

		// method receiver: func (x *X) XXX()里的(x *X)部分
		var recvName string
		if declValue.Recv != nil { // 方法
			// debug.Printf("FundDecl recv: %v\n", declValue.Recv.List)

			fieldResult := ins.inspectFields(declValue.Recv, from)
			recvName = fieldResult.RecvName
			method.Recv = recvName

			result.methodMap[recvName] = append(result.methodMap[recvName], method)
		}

		// 函数和方法
		result.funcMap[obj.Name()] = method

	case *ast.GenDecl:
		switch declValue.Tok {
		case token.IMPORT:
		case token.CONST:
		case token.VAR:
		case token.TYPE:
			for _, spec := range declValue.Specs {
				specResult := ins.inspectSpec(spec, from)
				for k, v := range specResult.structMap {
					result.structMap[k] = v
				}
				for k, v := range specResult.interfaceMap {
					result.interfaceMap[k] = v
				}
			}
		}
	}

	return
}

func (ins *Inspector) inspectSpec(spec ast.Spec, from string) (result SpecResult) {
	if spec == nil {
		return
	}
	result = MakeSpecResult()

	switch specValue := spec.(type) {
	case *ast.ImportSpec:
		// debug.Printf("ImportSpec, name: %v, path: %v\n", specValue.Name, specValue.Path)

	case *ast.ValueSpec:
		// debug.Printf("ValueSpec, name: %+v, type: %+v, value: %+v\n", specValue.Names, specValue.Type, specValue.Values)

	case *ast.TypeSpec:
		// 这里拿到类型信息: 名称，注释，文档
		// debug.Printf("TypeSpec name: %s, type: %+v, comment: %s, doc: %s\n", specValue.Name, specValue.Type, specValue.Comment.Text(), specValue.Doc.Text())

		switch specValue.Type.(type) {
		case *ast.InterfaceType:
			exprResult := ins.inspectExpr(specValue.Type, from)
			_ = exprResult
			// debug.Printf("interface type name: %s, exprValue: %+v, type: %+v, result: %+v\n", specValue.Name, specValue, specValue.Type, exprResult)

			interType := ins.pkg.TypesInfo.TypeOf(specValue.Type)
			r := parseTypesType(interType, parseTypesTypeOption{pkgPath: ins.pkg.PkgPath})
			methods := r.methods

			inter := Interface{
				Interface: ins.pkg.TypesInfo.Types[specValue.Type].Type.(*types.Interface),
				Name:      specValue.Name.Name,
				PkgPath:   ins.pkg.PkgPath,
				PkgName:   ins.pkg.Name,
				Methods:   methods,
			}
			mock, imports := inter.MakeMock()
			_, _ = mock, imports
			// debug.Printf("mock: %s, imports: %v\n", mock, imports)
			result.interfaceMap[specValue.Name.Name] = inter

		default:
			structOne := Struct{
				PkgPath: ins.pkg.PkgPath,
				PkgName: ins.pkg.Name,
				Field: Field{
					Id:        ins.pkg.TypesInfo.Types[specValue.Type].Type.String(),
					Name:      specValue.Name.Name,
					TypesType: ins.pkg.TypesInfo.Types[specValue.Type].Type,
					Type:      toString(specValue.Type),
					Doc:       specValue.Doc.Text(),
					Comment:   specValue.Comment.Text(),
				},
			}

			// 再拿field
			exprResult := ins.inspectExpr(specValue.Type, from)
			structOne.Fields = exprResult.Fields
			result.structMap[specValue.Name.Name] = structOne
		}
	}

	return
}

func (ins *Inspector) inspectExpr(expr ast.Expr, from string) (result ExprResult) {
	if expr == nil {
		return
	}
	result = MakeExprResult()

	switch exprValue := expr.(type) {
	case *ast.StructType:
		fieldResult := ins.inspectFields(exprValue.Fields, from)
		result.Fields = fieldResult.Fields

	case *ast.StarExpr: // *T
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)

	case *ast.TypeAssertExpr: // X.(*T)
		ins.inspectExpr(exprValue.X, from)
		ins.inspectExpr(exprValue.Type, from)

	case *ast.ArrayType: // [L]T
		ins.inspectExpr(exprValue.Len, from)
		ins.inspectExpr(exprValue.Elt, from)

	case *ast.BadExpr:
		panic(fmt.Errorf("BadExpr: %+v", exprValue))

	case *ast.IndexListExpr:
		ins.inspectExpr(exprValue.X, from)
		for _, indice := range exprValue.Indices {
			ins.inspectExpr(indice, from)
		}

	case *ast.SelectorExpr: // X.M
		// debug.Printf("SelectorExpr value: %v, typesString: %s\n", exprValue, toString(exprValue))

		exprResult := ins.inspectExpr(exprValue.X, from) // 也会进到下面的*ast.CallExpr分支
		result = result.Merge(exprResult)

		pkgID, ok := exprValue.X.(*ast.Ident)
		if ok {
			if so, ok := ins.pkg.TypesInfo.Uses[pkgID].(*types.PkgName); ok {
				pkgPath := so.Imported().Path()
				// debug.Printf(from+"SelectorExpr pkgPath: %#v\n", pkgPath)
				result.pkgPath = pkgPath
			}
		}

		// debug.Printf(from+"SelectorExpr value: %#v, result: %#v\n", exprValue, result)

	case *ast.SliceExpr: // []T, slice[1:3:5]
		ins.inspectExpr(exprValue.X, from)
		ins.inspectExpr(exprValue.Low, from)
		ins.inspectExpr(exprValue.High, from)
		ins.inspectExpr(exprValue.Max, from)

	case *ast.BasicLit: // 33 40.0 0x1f

	case *ast.BinaryExpr: // X+Y X-Y X*Y X/Y X%Y
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)
		exprResult = ins.inspectExpr(exprValue.Y, from)
		result = result.Merge(exprResult)
		// debug.Printf(from+"BinaryExpr: %+v\n", result)

	case *ast.CallExpr: // M(1, 2)
		// debug.Printf(from+"funcMap 1: %#v, %+v\n", exprValue.Fun, result)
		exprResult := ins.inspectExpr(exprValue.Fun, from)
		// debug.Printf(from+"funcMap mid: %#v, %+v\n", exprValue.Fun, exprResult)

		result.funcMap[toString(exprValue.Fun)] = Func{
			PkgPath: exprResult.pkgPath,
			Name:    toString(exprValue.Fun),
		}

		result = result.Merge(exprResult)
		// debug.Printf(from+"funcMap 2: %#v, %+v\n", exprValue.Fun, result)

		for _, arg := range exprValue.Args {
			// debug.Printf("CallExpr: %+v, %+v\n", exprValue.Fun, arg)
			exprResult := ins.inspectExpr(arg, from)
			result = result.Merge(exprResult)
		}
		// debug.Printf(from+"funcMap: %+v\n", result)

	case *ast.ChanType: // chan T, <-chan T, chan<- T
		exprResult := ins.inspectExpr(exprValue.Value, from)
		result = result.Merge(exprResult)

	case *ast.CompositeLit: // T{Name: Value}
		ins.inspectExpr(exprValue.Type, from)
		for _, elt := range exprValue.Elts {
			exprResult := ins.inspectExpr(elt, from)
			result = result.Merge(exprResult)
		}

	case *ast.Ellipsis: // ...int, [...]Arr
		ins.inspectExpr(exprValue.Elt, from)

	case *ast.FuncLit:
		ins.inspectExpr(exprValue.Type, from)
		ins.inspectStmt(exprValue.Body, from)

	case *ast.FuncType:
		ins.inspectFields(exprValue.Params, from)
		ins.inspectFields(exprValue.Results, from)

	case *ast.Ident:

		// if exprValue != nil {
		// debug.Printf(from+"Ident, name: %s, obj: %+v\n", exprValue.Name, exprValue.Obj)
		// } else {
		// debug.Printf(from+"Ident is nil: %+v\n", expr)
		// }

		obj, ok := ins.pkg.TypesInfo.Uses[exprValue]
		if ok {
			if obj.Pkg() != nil {
				_ = obj.Pkg().Path() // 变量的包路径

				// 变量类型的包路径
				var varTypePkgPath string
				if ptr, ok := obj.Type().(*types.Pointer); ok {
					// FIXME:改用parseTypesType统一处理types.Type信息
					switch ptrElem := ptr.Elem().(type) {
					case *types.Named:
						varTypePkgPath = ptrElem.Obj().Pkg().Path()
						// debug.Printf(from+"Ident obj: %#v, ptr: %#v, pkgPath: %#v\n", obj.Type(), ptr, varTypePkgPath)
					}
				}
				result.pkgPath = varTypePkgPath
			}
		}

		// debug.Printf(from+"Ident value: %#v, result: %#v\n", exprValue, result)

	case *ast.IndexExpr: // s[1], arr[1]
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)
		exprResult = ins.inspectExpr(exprValue.Index, from)
		result = result.Merge(exprResult)

	case *ast.InterfaceType: // interface { A(); B() }
		fieldResult := ins.inspectFields(exprValue.Methods, from)
		result.Fields = fieldResult.Fields

	case *ast.KeyValueExpr: // key:value
		ins.inspectExpr(exprValue.Key, from)
		exprResult := ins.inspectExpr(exprValue.Value, from)
		result = result.Merge(exprResult)

	case *ast.MapType: // map[string]T
		exprResult := ins.inspectExpr(exprValue.Key, from)
		result = result.Merge(exprResult)
		exprResult = ins.inspectExpr(exprValue.Value, from)
		result = result.Merge(exprResult)

	case *ast.ParenExpr: // (1==1)
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)

	case *ast.UnaryExpr: // *a
		exprResult := ins.inspectExpr(exprValue.X, from)
		result = result.Merge(exprResult)

	}

	return
}

func (ins *Inspector) inspectStmt(stmt ast.Stmt, from string) (result StmtResult) {
	if stmt == nil {
		return
	}
	result = MakeStmtResult()

	switch stmtValue := stmt.(type) {
	case *ast.AssignStmt: // a, b := 1, 2
		for _, lhs := range stmtValue.Lhs {
			ins.inspectExpr(lhs, from)
		}
		for _, rhs := range stmtValue.Rhs {
			exprResult := ins.inspectExpr(rhs, from)
			result = result.MergeExprResult(exprResult)
		}

	case *ast.SelectStmt: // select { }
		stmtResult := ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)

	case *ast.SendStmt: // c <- 1
		ins.inspectExpr(stmtValue.Chan, from)
		exprResult := ins.inspectExpr(stmtValue.Value, from)
		result = result.MergeExprResult(exprResult)

	case *ast.SwitchStmt: // switch { }
		stmtResult := ins.inspectStmt(stmtValue.Init, from)
		result = result.Merge(stmtResult)
		exprResult := ins.inspectExpr(stmtValue.Tag, from)
		result = result.MergeExprResult(exprResult)
		stmtResult = ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)

	case *ast.BadStmt:
		panic(fmt.Errorf("BadStmt: %+v", stmtValue))

	case *ast.BlockStmt:
		if stmtValue != nil {
			for _, single := range stmtValue.List {
				// debug.Printf(from+"block stmt: %+v\n", single)
				res := ins.inspectStmt(single, from)
				result = result.Merge(res)
			}
		}
		// debug.Printf(from+"block funcMap: %+v\n", result.funcMap)

	case *ast.BranchStmt:
		exprResult := ins.inspectExpr(stmtValue.Label, from)
		result = result.MergeExprResult(exprResult)

	case *ast.CaseClause:
		for _, one := range stmtValue.List {
			exprResult := ins.inspectExpr(one, from)
			result = result.MergeExprResult(exprResult)
		}
		for _, one := range stmtValue.Body {
			stmtResult := ins.inspectStmt(one, from)
			result = result.Merge(stmtResult)
		}

	case *ast.CommClause:
		stmtResult := ins.inspectStmt(stmtValue.Comm, from)
		result = result.Merge(stmtResult)
		for _, one := range stmtValue.Body {
			stmtResult := ins.inspectStmt(one, from)
			result = result.Merge(stmtResult)
		}

	case *ast.DeclStmt:
		ins.inspectDecl(stmtValue.Decl, from)

	case *ast.DeferStmt:
		exprResult := ins.inspectExpr(stmtValue.Call, from)
		result = result.MergeExprResult(exprResult)

	case *ast.EmptyStmt:

	case *ast.ExprStmt:
		// debug.Printf(from+"expr stmt: %+v\n", stmtValue.X)
		exprResult := ins.inspectExpr(stmtValue.X, from)
		result = result.MergeExprResult(exprResult)
		// debug.Printf(from+"expr funcMap: %+v\n", result.funcMap)

	case *ast.ForStmt: // for i:=0; i< l; i++ { }
		ins.inspectStmt(stmtValue.Init, from)
		exprResult := ins.inspectExpr(stmtValue.Cond, from)
		result = result.MergeExprResult(exprResult)
		ins.inspectStmt(stmtValue.Post, from)
		stmtResult := ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)

	case *ast.GoStmt:
		exprResult := ins.inspectExpr(stmtValue.Call, from)
		result = result.MergeExprResult(exprResult)

	case *ast.IfStmt:
		stmtResult := ins.inspectStmt(stmtValue.Init, from)
		result = result.Merge(stmtResult)
		exprResult := ins.inspectExpr(stmtValue.Cond, from)
		result = result.MergeExprResult(exprResult)
		stmtResult = ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)
		stmtResult = ins.inspectStmt(stmtValue.Else, from)
		result = result.Merge(stmtResult)

	case *ast.IncDecStmt:
		exprResult := ins.inspectExpr(stmtValue.X, from)
		result = result.MergeExprResult(exprResult)

	case *ast.LabeledStmt:
		exprResult := ins.inspectExpr(stmtValue.Label, from)
		result = result.MergeExprResult(exprResult)
		ins.inspectStmt(stmtValue.Stmt, from)

	case *ast.RangeStmt: // for key, value := range slice { }
		ins.inspectExpr(stmtValue.Key, from)
		ins.inspectExpr(stmtValue.Value, from)
		exprResult := ins.inspectExpr(stmtValue.X, from)
		result = result.MergeExprResult(exprResult)
		stmtResult := ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)

	case *ast.ReturnStmt:
		for _, one := range stmtValue.Results {
			exprResult := ins.inspectExpr(one, from)
			result = result.MergeExprResult(exprResult)
			// debug.Printf(from+"return stmt: %#v, %+v\n", one, result.funcMap)
		}

	case *ast.TypeSwitchStmt: // switch x := m(); a := x.(type) { }
		stmtResult := ins.inspectStmt(stmtValue.Init, from)
		result = result.Merge(stmtResult)
		stmtResult = ins.inspectStmt(stmtValue.Assign, from)
		result = result.Merge(stmtResult)
		stmtResult = ins.inspectStmt(stmtValue.Body, from)
		result = result.Merge(stmtResult)
	}

	return
}

func (ins *Inspector) inspectFields(fields *ast.FieldList, from string) (result FieldResult) {
	if fields == nil {
		return
	}
	result = MakeFieldResult()

	var _ *ast.Field // 是一个Node，但不是一个Expr

	for _, field := range fields.List {
		// 拿field的名称，类型，tag，注释，文档
		// debug.Printf("StructType field name: %v, type: %+v, tag: %v, comment: %s, doc: %s\n", field.Names, field.Type, field.Tag, field.Comment.Text(), field.Doc.Text())

		// 获取receiver name
		fieldTyp := field.Type
		if singleTyp, ok := field.Type.(*ast.StarExpr); ok {
			fieldTyp = singleTyp.X
		}
		result.RecvName = toString(fieldTyp)

		ins.inspectExpr(field.Type, from)

		name := ""
		anonymous := false
		if len(field.Names) != 0 {
			for _, s := range field.Names {
				name += s.Name
			}
		} else {
			// 匿名结构体
			name = toString(field.Type)
			anonymous = true
		}

		tag := ""
		if field.Tag != nil {
			tag = field.Tag.Value
		}
		result.Fields = append(result.Fields, Field{
			Id:          name,
			Name:        name,
			Anonymous:   anonymous,
			TypesType:   ins.pkg.TypesInfo.TypeOf(field.Type),
			Type:        toString(field.Type),
			Tag:         tag,
			TagBasicLit: field.Tag,
			Doc:         field.Doc.Text(),
			Comment:     field.Comment.Text(),
		})
	}

	return
}

func toString(v any) string {
	qualifier := pkgNameQualifier(qualifierParam{})

	switch vv := v.(type) {
	case ast.Expr:
		return types.ExprString(vv)
	case types.Type:
		return types.TypeString(vv, qualifier)
	case types.Object:
		return types.ObjectString(vv, qualifier)
	case *types.Selection:
		return types.SelectionString(vv, qualifier)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getTypesPkgPath(t types.Type) string {
	// debug.Printf("pvar type: %s\n", t)

	pkgPath := ""
	switch v := t.(type) {
	case *types.Named:
		if v.Obj().Pkg() != nil {
			pkgPath = v.Obj().Pkg().Path()
			// debug.Printf("path: %s\n", pkgPath)
		}
	}

	return pkgPath
}

type parseTypesTypeOption struct {
	_       string
	pkgPath string
}

func parseTypesType(t types.Type, opt parseTypesTypeOption) (r struct {
	methods []Method
}) {
	switch tv := t.(type) {
	case *types.Interface:
		methods := make([]Method, 0, tv.NumMethods())
		for i := 0; i < tv.NumMethods(); i++ {
			met := tv.Method(i)
			methods = append(methods, Method{
				Origin:    met,
				Name:      met.Name(),
				Signature: types.TypeString(met.Type(), pkgNameQualifier(qualifierParam{pkgPath: opt.pkgPath})),
			})
		}
		// debug.Printf("| parseTypesType | interface methods: %+v\n", methods)
		r.methods = methods

	case *types.Signature:
		// debug.Printf("=== signature: %+v, %+v, %+v\n", tv, tv.Params(), tv.Results())

	case *types.Pointer:
		// debug.Printf("=== pointer: %+v, %+v\n", tv, tv.Elem())
		parseTypesType(tv.Elem(), opt)

	case *types.Named:
		methods := []Method{}
		for i := 0; i < tv.NumMethods(); i++ {
			met := tv.Method(i)
			methods = append(methods, Method{
				Origin:    met,
				Signature: types.TypeString(met.Type(), pkgNameQualifier(qualifierParam{pkgPath: opt.pkgPath})),
			})
		}
		_ = methods
		// debug.Printf("=== named: %+v, is alias: %v, pkgPath: %v, methods: %+v\n", tv, tv.Obj().IsAlias(), tv.Obj().Pkg().Path(), methods)
		// if tv.Obj().IsAlias() {
		// debug.Printf("===============================: %+v\n", tv)
		// }

	case *types.Struct:
		fields := []Field{}
		for i := 0; i < tv.NumFields(); i++ {
			field := tv.Field(i)

			tmpField := Field{
				Id:   field.Id(),
				Name: field.Name(),
				Type: types.TypeString(field.Type(), pkgNameQualifier(qualifierParam{pkgPath: opt.pkgPath})),
			}
			fields = append(fields, tmpField)
		}
		_ = fields
		// debug.Printf("=== struct: %+v, fields: %+v\n", tv, fields)

	case *types.Slice:
		// debug.Printf("| parseTypesType | elem: %+v\n", tv.Elem())
		parseTypesType(tv.Elem(), opt)

	case *types.Array:
		// debug.Printf("| parseTypesType | elem: %+v\n", tv.Elem())
		parseTypesType(tv.Elem(), opt)

	case *types.Basic:
		// debug.Printf("| parseTypesType | elem: %+v\n", tv.Info())

	case *types.Chan:
		// debug.Printf("| parseTypesType | elem: %+v\n", tv.Elem())
		parseTypesType(tv.Elem(), opt)

	case *types.Map:
		// debug.Printf("| parseTypesType | key: %+v, value: %+v\n", tv.Key(), tv.Elem())
		parseTypesType(tv.Key(), opt)
		parseTypesType(tv.Elem(), opt)

	case *types.Tuple:
		// debug.Printf("| parseTypesType | len: %+v\n", tv.Len())

	default:
		fmt.Printf("| parseTypesType | tv: %+v\n", tv)
	}

	return
}
