package sqlparser

import (
	"fmt"
	"strconv"
	"strings"

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
