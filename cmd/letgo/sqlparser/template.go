package sqlparser

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/donnol/do"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

const (
	tableTmpl = `
	// {{.StructName}} {{.StructComment}}
	type {{.StructName}} struct {
		{{- range $k,$v := .Fields}}
		{{$v.FieldName}} {{$v.FieldType}} {{$v.FieldTag}} // {{$v.FieldComment -}}
		{{- end}}
	}
	
	func ({{.StructName}}) TableName() string {
		return "{{.TableName}}"
	}
	
	func (s {{.StructName}}) Columns() []string {
		return s.NameHelper().Columns()
	}
	
	func (s {{.StructName}}) Values() []any {
		return []any{
			{{- range $k,$v := .Fields}}
			s.{{.FieldName -}},
			{{- end}}
		}
	}
	
	func (s *{{.StructName}}) ValuePtrs() []any {
		return []any{
			{{- range $k,$v := .Fields}}
			&s.{{.FieldName -}},
			{{- end}}
		}
	}
	
	func (s {{.StructName}}) Exists() bool {
		return s.Id != 0
	}
	
	type _{{.StructName}}NameHelper struct {
		{{- range $k,$v := .Fields}}
		{{.FieldName}} string // field: {{.DBField -}}
		{{- end}}
	}
	
	// FuzzWrap make v become %v%
	func (_{{.StructName}}NameHelper) FuzzWrap(v string) string {
		return "%" + v + "%"
	}
	
	func (_{{.StructName}}NameHelper) Columns() []string {
		return []string{
			{{- range $k,$v := .Fields}}
			"{{.DBField -}}",
			{{- end}}
		}
	}
	
	func ({{.StructName}}) NameHelper() _{{.StructName}}NameHelper {
		return _{{.StructName}}NameHelper{
			{{- range $k,$v := .Fields}}
			{{.FieldName}}: "{{.DBField -}}",
			{{- end}}
		}
	}
	
	{{if .HaveEnum -}}
	type _{{.StructName}}Enum struct {
		{{range $key,$value := .EnumFields}}
			{{$value.FieldName}} struct {
			{{range $ikey,$ivalue := $value.EnumFieldValues}} 
				E_{{$ivalue.EnumValue}} do.Enum[{{$ivalue.FieldType}}] // {{$ivalue.EnumComment -}}
			{{- end}}
		} // {{$value.FieldComment -}}
		{{- end}}
	}
	
	func ({{.StructName}}) EnumHelper() _{{.StructName}}Enum {
		e := _{{.StructName}}Enum{}
		{{range $key,$value := .EnumFields}}
			{{range $ikey,$ivalue := $value.EnumFieldValues}} 
				e.{{$value.FieldName}}.E_{{$ivalue.EnumValue}} = do.Enum[{{$ivalue.FieldType}}]{Name: "{{$ivalue.EnumName}}", Value: {{$ivalue.EnumValueProcess}}} 	
			{{- end}}
		{{- end}}
		return e
	}
	
	var _ = func() struct{} {
		e := {{.StructName}}{}.EnumHelper()
		{{range $key,$value := .EnumFields}}
			{{range $ikey,$ivalue := $value.EnumFieldValues}}
				if e.{{$value.FieldName}}.E_{{$ivalue.EnumValue}}.Value != {{$ivalue.EnumValueProcess}} || e.{{$value.FieldName}}.E_{{$ivalue.EnumValue}}.Name != "{{$ivalue.EnumName}}" {
					panic("invalid enum")
				} 	
			{{- end}}
		{{- end}}
		return struct{}{}
	}()
	{{- end}}
`
)

type StructForTmpl struct {
	PkgName       string
	Imports       []string
	TableName     string
	StructName    string
	StructComment string
	Fields        []StructField
	EnumFields    []EnumField
	HaveEnum      bool
}

type StructField struct {
	FieldName    string
	FieldType    string
	FieldTag     string
	FieldComment string
	DBField      string
}

type EnumField struct {
	StructField
	EnumFieldValues []EnumFieldValue
}

type EnumFieldValue struct {
	StructField
	EnumName         string
	EnumValue        string
	EnumValueProcess string // 如果是字符串类型，需要在EnumValue基础上添加双引号
	EnumComment      string
}

func FromStruct(s *Struct, opt Option) StructForTmpl {
	fields := make([]StructField, 0, len(s.Fields))
	efields := make([]EnumField, 0, len(s.Fields))
	for _, field := range s.Fields {
		fieldName := field.Name
		if len(opt.IgnoreField) > 0 {
			if lo.IndexOf(opt.IgnoreField, fieldName) > -1 {
				continue
			}
		}
		if opt.FieldNameMapper != nil {
			fieldName = opt.FieldNameMapper(fieldName)
		}
		// 与TableName()方法重名时，添加后缀；因为用了tag来与数据库字段对应，所以影响不大
		if fieldName == "TableName" {
			fieldName += "Field"
		}

		fieldType := field.Type
		if opt.FieldTypeMapper != nil {
			fieldType = opt.FieldTypeMapper(fieldType)
		}
		// 根据comment里的 type(do.Id) 获得类型
		{
			const bs = "type("
			bi := strings.Index(field.Comment, bs)
			ei := strings.Index(field.Comment, ")")
			if bi != -1 && ei != -1 && bi < ei {
				es := field.Comment[bi+len(bs) : ei]
				fieldType = es
				field.Comment = strings.ReplaceAll(field.Comment, field.Comment[bi:ei+1], "")
			}
		}

		fieldTag := field.Tag
		if opt.FieldTagMapper != nil {
			fieldTag = opt.FieldTagMapper(field.Name, field.Type)
		}

		structField := StructField{
			FieldName:    fieldName,
			FieldType:    fieldType,
			FieldTag:     fieldTag,
			FieldComment: field.Comment,
			DBField:      field.DBField,
		}
		fields = append(fields, structField)

		if len(field.Enums) > 0 {
			s.HaveEnum = true

			fieldEnum := EnumField{
				StructField: structField,
			}
			for _, e := range field.Enums {
				ev := e.EnumValue
				if fieldType == "string" {
					ev = fmt.Sprintf("%q", ev)
				} else {
					_, err1 := strconv.ParseInt(ev, 10, 64)
					_, err2 := strconv.ParseUint(ev, 10, 64)
					if err1 != nil && err2 != nil {
						ev = fmt.Sprintf("%q", ev)
					}
				}
				efv := EnumFieldValue{
					StructField:      structField,
					EnumName:         e.EnumName,
					EnumValue:        e.EnumValue,
					EnumValueProcess: ev,
					EnumComment:      e.EnumValue + " " + e.EnumName,
				}
				fieldEnum.EnumFieldValues = append(fieldEnum.EnumFieldValues, efv)
			}
			efields = append(efields, fieldEnum)
		}
	}
	return StructForTmpl{
		PkgName:       s.PkgName,
		TableName:     s.TableName,
		StructName:    s.Name,
		StructComment: s.Comment,
		Fields:        fields,
		EnumFields:    efields,
		HaveEnum:      s.HaveEnum,
	}
}

const (
	ResultToJSONObjectTmpl = `json_object(
		{{- range $k,$v := .Fields}}
		'{{$v.JSONName}}', {{$v.ColumnName -}} {{- if $v.NoComma}} {{- else}},{{- end}}
		{{- end}}
	)`
)

type ResultToJSONObject struct {
	Fields []ResultToJSONObjectField
}

type ResultToJSONObjectField struct {
	JSONName   string
	ColumnName string
	NoComma    bool
}

func FromStructForTmpl(s *StructForTmpl) *ResultToJSONObject {
	fields := make([]ResultToJSONObjectField, 0, len(s.Fields))
	for i, field := range s.Fields {
		jname := strcase.ToCamel(field.DBField)
		jname = strings.ToLower(string(jname[0])) + jname[1:]
		fields = append(fields, ResultToJSONObjectField{
			JSONName:   jname,
			ColumnName: field.DBField,
			NoComma:    i == len(s.Fields)-1,
		})
	}
	return &ResultToJSONObject{
		Fields: fields,
	}
}

const (
	SqlInsertTmpl = `
	INSERT {{if .Ignore}}IGNORE{{end}} INTO {{.Table}} (
		{{- range $k,$v:=.Fields -}}
			{{$v.Name}}{{if $v.NoDot}}{{else}}, {{end -}}
		{{- end -}}
	) VALUES 
	{{- range $vk, $vv:=.Values}}
	(
		{{- range $ik, $iv := $vv.FieldValues -}}
			{{$iv.FieldValue}}{{if $iv.NoDot}}{{else}}, {{end -}}
		{{- end -}}
	){{if $vv.NoDot}};{{else}}, {{end -}}
	{{end -}}
	`
)

type InsertParam struct {
	Ignore bool
	Table  string
	Fields []InsertField
	Values []InsertValue
}

type InsertField struct {
	Name  string
	NoDot bool
}

type InsertValue struct {
	FieldValues []InsertFieldValue
	NoDot       bool
}

type InsertFieldValue struct {
	FieldValue string
	NoDot      bool
}

var (
	insertValues = make(map[string]map[string][]string)
)

func InsertParamFromStruct(s *Struct, opt *Option) *InsertParam {
	opt.fillByDefault()

	fields := make([]InsertField, 0, len(s.Fields))
	values := make([]InsertValue, 0, len(s.Fields))

	for i, field := range s.Fields {
		fields = append(fields, InsertField{
			Name:  field.DBField,
			NoDot: i == len(s.Fields)-1,
		})
	}

	for i := 0; i < opt.Amount; i++ {
		fieldValues := make([]InsertFieldValue, 0, len(s.Fields))
		for j, field := range s.Fields {
			fieldType := field.Type
			if opt.FieldTypeMapper != nil {
				fieldType = opt.FieldTypeMapper(fieldType)
			}
			// 当字段关联其它表时，输入已生成记录的id列表，从中任选一个 -- ref(table.id)
			choose := []string{}
			if iv, ok := insertValues[field.Ref.Table]; ok {
				vv, ok := iv[field.Ref.Field]
				if ok {
					choose = vv
				}
			}
			var value string
			if len(choose) != 0 {
				value = choose[rand.Intn(len(choose))]
			} else {
				value = valueByType(fieldType)
			}
			fieldValues = append(fieldValues, InsertFieldValue{
				FieldValue: value,
				NoDot:      j == len(s.Fields)-1,
			})

			if _, ok := insertValues[s.TableName]; !ok {
				insertValues[s.TableName] = make(map[string][]string)
			}
			insertValues[s.TableName][field.DBField] = append(insertValues[s.TableName][field.DBField], value)
		}

		values = append(values, InsertValue{
			FieldValues: fieldValues,
			NoDot:       i == opt.Amount-1,
		})
	}

	return &InsertParam{
		Ignore: true,
		Table:  s.TableName,
		Fields: fields,
		Values: values,
	}
}

func valueByType(typ string) string {
	switch typ {
	case "bool":
		if gofakeit.Bool() {
			return "1"
		} else {
			return "0"
		}
	case "int":
		return strconv.FormatInt(int64(gofakeit.Int32()), 10)
	case "int64":
		return strconv.FormatInt(gofakeit.Int64(), 10)
	case "uint":
		return strconv.FormatUint(uint64(gofakeit.Uint32()), 10)
	case "uint64":
		return strconv.FormatUint(gofakeit.Uint64(), 10)
	case "uint16":
		return strconv.FormatUint(uint64(gofakeit.Uint16()), 10)
	case "int16":
		return strconv.FormatInt(int64(gofakeit.Int16()), 10)
	case "uint8":
		return strconv.FormatUint(uint64(gofakeit.Uint8()), 10)
	case "int8":
		return strconv.FormatInt(int64(gofakeit.Int8()), 10)
	case "float64":
		return strconv.FormatFloat(gofakeit.Float64(), 'e', 2, 64)
	case "float32":
		return strconv.FormatFloat(float64(gofakeit.Float32()), 'e', 2, 64)
	case "string", "[]byte":
		return withQuote(gofakeit.Name())
	case "json.RawMessage":
		j := do.Must1(gofakeit.JSON(nil))
		return withQuote(string(j))
	case "time.Time":
		return withQuote(gofakeit.Date().Format(do.DateTimeFormat))
	}
	return withQuote(gofakeit.Name())
}

func withQuote(s string) string {
	return `'` + s + `'`
}
