package sqlparser

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iancoleman/strcase"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/samber/lo"
)

type (
	ParseOption struct {
		tablePrefix string
	}
	ParseSetter func(*ParseOption)
)

func WithTablePrefix(tablePrefix string) ParseSetter {
	return func(po *ParseOption) {
		po.tablePrefix = tablePrefix
	}
}

func ParseCreateSQL(sql string) *Struct {
	s := &Struct{}

	node, err := parse(sql)
	if err != nil {
		log.Fatal(err)
	}
	(*node).Accept(s)

	return s
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &stmtNodes, nil
}

func ParseCreateSQLBatch(sql string, opts ...ParseSetter) []*Struct {
	opt := &ParseOption{}
	for _, setter := range opts {
		setter(opt)
	}

	r := make([]*Struct, 0)

	nodes, err := parseBatch(sql)
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range nodes {
		s := &Struct{
			tablePrefix: opt.tablePrefix,
		}
		node.Accept(s)
		if s.TableName == "" {
			continue
		}
		r = append(r, s)
	}

	return r
}

func parseBatch(sql string) ([]ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.ParseSQL(sql)
	if err != nil {
		return nil, err
	}
	return stmtNodes, err
}

type colX struct {
	colNames []string
}

func (v *colX) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.ColumnName); ok {
		v.colNames = append(v.colNames, name.Name.O)
	}
	return in, false
}

func (v *colX) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func extract(rootNode *ast.StmtNode) []string {
	v := &colX{}
	(*rootNode).Accept(v)
	return v.colNames
}

type Struct struct {
	tablePrefix string

	PkgName   string
	Name      string
	TableName string
	Comment   string
	Fields    []Field

	HaveEnum bool
}
type Field struct {
	Name    string
	Type    string
	Tag     string
	DBField string
	Comment string
	Enums   []Enum
}
type Enum struct {
	EnumName  string
	EnumValue string
}

func (v *Struct) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.CreateTableStmt:
		v.TableName = node.Table.Name.O
		for _, opt := range node.Options {
			if opt.Tp == ast.TableOptionComment {
				v.Comment = opt.StrValue
			}
		}
		for _, col := range node.Cols {
			field := Field{
				Name: col.Name.Name.O,
			}
			field.DBField = col.Name.Name.L

			field.Type = col.Tp.InfoSchemaStr()
			field.Type = processFieldType(field.Type)

			for _, opt := range col.Options {
				switch opt.Tp {
				case ast.ColumnOptionPrimaryKey:
				case ast.ColumnOptionNotNull:
				case ast.ColumnOptionAutoIncrement:
				case ast.ColumnOptionDefaultValue:
				case ast.ColumnOptionUniqKey:
				case ast.ColumnOptionNull:
				case ast.ColumnOptionOnUpdate:
				case ast.ColumnOptionFulltext:
				case ast.ColumnOptionComment:
					field.Comment = opt.Expr.(ast.ValueExpr).GetDatumString()
					// (查找"enum("和")"之间内容，内容可以用中文或英文分号分割)
					const bs = "enum("
					bi := strings.Index(field.Comment, bs)
					ei := strings.Index(field.Comment, ")")
					if bi != -1 && ei != -1 && bi < ei {
						es := field.Comment[bi+len(bs) : ei]
						es = strings.ReplaceAll(es, "；", ";")
						parts := strings.Split(es, ";")
						for _, part := range parts {
							if part == "" {
								continue
							}
							vparts := strings.Split(part, " ")
							if len(vparts) < 2 {
								continue
							}
							field.Enums = append(field.Enums, Enum{
								EnumName:  vparts[1],
								EnumValue: vparts[0],
							})
						}
						// fmt.Printf("%+v\n", field.Enums)
					}
				case ast.ColumnOptionGenerated:
				case ast.ColumnOptionReference:
				case ast.ColumnOptionCollate:
				case ast.ColumnOptionCheck:
				case ast.ColumnOptionColumnFormat:
				case ast.ColumnOptionStorage:
				case ast.ColumnOptionAutoRandom:
				}
			}
			v.Fields = append(v.Fields, field)
		}
	}
	return in, false
}

func (v *Struct) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func processFieldType(fieldType string) string {
	// 去掉括号及其内部内容
	re := regexp.MustCompile(`([(\d)(,)])`)
	fieldType = re.ReplaceAllString(fieldType, "")
	return fieldType
}

type Option struct {
	StructNameMapper       func(string) string         // 名称映射
	IgnoreField            []string                    // 忽略字段
	FieldNameMapper        func(string) string         // 字段名称映射
	FieldTypeMapper        func(string) string         // 字段类型映射
	FieldTagMapper         func(string, string) string // 可根据名称和类型自行决定字段tag
	RandomFieldValueByType func(string) any            // 根据字段类型生成随机值
	TrimTablePrefix        string                      // 表名前缀，若有则在转为结构体名前去掉
}

var (
	doption = Option{
		StructNameMapper: func(k string) string {
			// snake to camel
			return strcase.ToCamel(k)
		},
		FieldNameMapper: func(k string) string {
			return strcase.ToCamel(k)
		},
		FieldTypeMapper: func(k string) string {
			var (
				// type mapper
				m = map[string]string{
					"BIT":               "bool",
					"TEXT":              "string",
					"BLOB":              "[]byte",
					"DATE":              "time.Time",
					"DATETIME":          "time.Time",
					"DECIMAL":           "float64",
					"DOUBLE":            "float64",
					"ENUM":              "",
					"FLOAT":             "float32",
					"GEOMETRY":          "",
					"MEDIUMINT":         "int",
					"JSON":              "json.RawMessage",
					"UNSIGNED INT":      "uint",
					"INT UNSIGNED":      "uint",
					"INT":               "int",
					"LONGTEXT":          "string",
					"LONGBLOB":          "[]byte",
					"UNSIGNED BIGINT":   "uint64",
					"BIGINT UNSIGNED":   "uint64",
					"BIGINT":            "int64",
					"MEDIUMTEXT":        "string",
					"MEDIUMBLOB":        "[]byte",
					"NULL":              "",
					"SET":               "",
					"UNSIGNED SMALLINT": "uint16",
					"SMALLINT UNSIGNED": "uint16",
					"SMALLINT":          "int16",
					"BINARY":            "[]byte",
					"CHAR":              "string",
					"TIME":              "time.Time",
					"TIMESTAMP":         "time.Time",
					"UNSIGNED TINYINT":  "uint8",
					"TINYINT UNSIGNED":  "uint8",
					"TINYINT":           "int8",
					"TINYTEXT":          "string",
					"TINYBLOB":          "[]byte",
					"VARBINARY":         "[]byte",
					"VARCHAR":           "string",
					"YEAR":              "time.Time",
				}
			)
			if v, ok := m[strings.ToUpper(k)]; ok {
				return v
			}
			return k
		},
		FieldTagMapper: func(name string, typ string) string {
			if name == "" {
				return ""
			}
			jname := strcase.ToCamel(name)
			jname = strings.ToLower(string(jname[0])) + jname[1:]
			return fmt.Sprintf("`json:\"%s\" db:\"%s\"`", jname, name)
		},
		RandomFieldValueByType: func(s string) any {
			var v any
			switch s {
			case "bool":
				v = gofakeit.Bool()
			case "string":
				v = strconv.Quote(gofakeit.HexUint256())
			case "[]byte":
				v = strconv.Quote(gofakeit.HexUint128())
			case "time.Time":
				t := gofakeit.Date()
				gofakeit.DateRange(
					time.Date(1970, 8, 8, 0, 0, 0, 0, time.Local),
					time.Date(2047, 1, 3, 0, 0, 0, 0, time.Local),
				)
				v = strconv.Quote(t.Format("2006-01-02 15:04:05"))
			case "float64":
				v = gofakeit.Float32()
			case "float32":
				v = gofakeit.Float64()
			case "json.RawMessage":
				t, _ := gofakeit.JSON(&gofakeit.JSONOptions{})
				v = strconv.Quote(string(t))
			case "uint":
				v = gofakeit.Uint32()
			case "int":
				v = gofakeit.Int32()
			case "uint64":
				v = gofakeit.Uint64()
			case "int64":
				v = gofakeit.Int64()
			case "uint16":
				v = gofakeit.Uint16()
			case "int16":
				v = gofakeit.Int16()
			case "uint8":
				v = gofakeit.Uint8()
			case "int8":
				v = gofakeit.Int8()
			}
			return v
		},
	}
)

func (opt *Option) fillByDefault() {
	if opt.StructNameMapper == nil {
		opt.StructNameMapper = doption.StructNameMapper
	}
	if opt.FieldNameMapper == nil {
		opt.FieldNameMapper = doption.FieldNameMapper
	}
	if opt.FieldTypeMapper == nil {
		opt.FieldTypeMapper = doption.FieldTypeMapper
	}
	if opt.FieldTagMapper == nil {
		opt.FieldTagMapper = doption.FieldTagMapper
	}
	if opt.RandomFieldValueByType == nil {
		opt.RandomFieldValueByType = doption.RandomFieldValueByType
	}
}

func (s *Struct) Gen(w io.Writer, opt Option) error {
	(&opt).fillByDefault()

	name := s.TableName
	// 去掉前缀
	if opt.TrimTablePrefix != "" {
		name = strings.TrimPrefix(name, opt.TrimTablePrefix)
	}
	if opt.StructNameMapper != nil {
		name = opt.StructNameMapper(name)
	}
	s.Name = name

	{
		temp, err := template.New("sql2struct").Parse(tableTmpl)
		if err != nil {
			return err
		}
		if err := temp.Execute(w, FromStruct(s, opt)); err != nil {
			return err
		}
	}

	return nil
}

var (
	insertHeadTmpl  = "INSERT IGNORE INTO `{{.TableName}}` "
	insertFieldTmpl = "({{ range $index, $element := .Fields }}" + "\n" +
		"`{{$element.Name}}`{{if ne $index $.FieldLastIndex }},{{end}}{{end}}" + "\n" + ") VALUES "
	insertValueTmpl = "({{ range $index, $element := .Values }}" + "\n" +
		"{{$element}}{{if ne $index $.FieldLastIndex }},{{end}}{{ end }}" + "\n" + ")"
)

func (s *Struct) GenData(w io.Writer, n int64, opt Option) error {
	(&opt).fillByDefault()

	if n <= 0 {
		n = 1
	}

	// gen random value by field type
	fields := make([]any, 0, len(s.Fields))
	values := make([][]any, 0, n)
	for i := int64(0); i < n; i++ {

		v := make([]any, 0, len(s.Fields))
		for _, field := range s.Fields {
			if len(opt.IgnoreField) > 0 {
				if lo.IndexOf(opt.IgnoreField, field.Name) > -1 {
					continue
				}
			}

			if i == 0 {
				fields = append(fields, field)
			}

			fieldType := field.Type
			if opt.FieldTypeMapper != nil {
				fieldType = opt.FieldTypeMapper(fieldType)
			}
			var value any = "[NULL]"
			if opt.RandomFieldValueByType != nil {
				value = opt.RandomFieldValueByType(fieldType)
			}
			v = append(v, value)
		}

		values = append(values, v)
	}

	{
		temp, err := template.New("insertHead").Parse(insertHeadTmpl)
		if err != nil {
			return err
		}
		if err := temp.Execute(w, map[string]any{
			"TableName": s.TableName,
		}); err != nil {
			return err
		}
	}

	{
		temp, err := template.New("insertField").Parse(insertFieldTmpl)
		if err != nil {
			return err
		}
		if err := temp.Execute(w, map[string]any{
			"Fields":         fields,
			"FieldLastIndex": len(fields) - 1,
		}); err != nil {
			return err
		}
	}

	{
		temp, err := template.New("insertValue").Parse(insertValueTmpl)
		if err != nil {
			return err
		}
		for i := int64(0); i < n; i++ {
			if err := temp.Execute(w, map[string]any{
				"Values":         values[i],
				"FieldLastIndex": len(fields) - 1,
			}); err != nil {
				return err
			}
			if i == n-1 {
				w.Write([]byte(";\n"))
			} else {
				w.Write([]byte(", "))
			}
		}
	}

	return nil
}
