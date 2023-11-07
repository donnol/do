package sqlparser

const (
	enumTmpl = `
type _{{.StructName}}Enum struct {
	{{range $key,$value := .EnumFields}}
		{{$value.FieldName}} struct {
		{{range $ikey,$ivalue := $value.EnumFieldValues}} 
			E_{{$ivalue.EnumValue}} do.Enum[{{$ivalue.FieldType}}] 	{{end}}
	} // {{$value.FieldComment}}	{{end}}
}

func ({{.StructName}}) EnumHelper() _{{.StructName}}Enum {
	e := _{{.StructName}}Enum{}
	{{range $key,$value := .EnumFields}}
		{{range $ikey,$ivalue := $value.EnumFieldValues}} 
			e.{{$value.FieldName}}.E_{{$ivalue.EnumValue}} = do.Enum[{{$ivalue.FieldType}}]{Name: "{{$ivalue.EnumName}}", Value: {{$ivalue.EnumValueProcess}}} 	{{end}}	{{end}}
	return e
}

var _ = func() struct{} {
	e := {{.StructName}}{}.EnumHelper()
	{{range $key,$value := .EnumFields}}
		{{range $ikey,$ivalue := $value.EnumFieldValues}}
			if e.{{$value.FieldName}}.E_{{$ivalue.EnumValue}}.Value != {{$ivalue.EnumValueProcess}} || e.{{$value.FieldName}}.E_{{$ivalue.EnumValue}}.Name != "{{$ivalue.EnumName}}" {
				panic("invalid enum")
			} 	{{end}}	{{end}}
	return struct{}{}
}
	`
)

type StructForTmpl struct {
	StructName    string
	StructComment string
	Fields        []StructField
	EnumFields    []EnumField
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
}
