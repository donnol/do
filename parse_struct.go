package do

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/donnol/do/parser"
)

type Field struct {
	reflect.StructField        // 内嵌反射结构体字段类型
	Comment             string // 注释
	Struct              Struct // 字段的类型是其它结构体
}

type Struct struct {
	Name        string       // 名字
	Comment     string       // 注释
	Description string       // 描述
	Type        reflect.Type // 反射类型
	Fields      []Field      // 结构体字段
}

func MakeStruct() Struct {
	return Struct{
		Fields: make([]Field, 0),
	}
}

func ResolveStruct(value any) (Struct, error) {
	s := MakeStruct()

	var refType reflect.Type
	if v, ok := value.(reflect.Type); ok {
		refType = v
	} else {
		refType = reflect.TypeOf(value)
	}
	s.Type = refType

	if refType == nil {
		return s, fmt.Errorf("nil refType")
	}

	if refType.Kind() == reflect.Ptr { // 指针
		refType = refType.Elem()
	}
	if refType.Kind() != reflect.Struct {
		return s, fmt.Errorf("bad value type , type is %v", refType.Kind())
	}
	structName := refType.PkgPath() + "." + refType.Name()
	s.Name = structName

	if refType.NumField() == 0 { // 空结构体
		return s, nil
	}
	if err := collectStructComment(refType, &s); err != nil {
		return s, err
	}

	return s, nil
}

// GetFields return all field in struct, include anonymous fields
func (s Struct) GetFields() []Field {
	return getFields(s)
}

func getFields(s Struct) []Field {
	var fields = make([]Field, 0)
	for _, f := range s.Fields {
		if f.Anonymous {
			fields = append(fields, getFields(f.Struct)...)
		} else {
			fields = append(fields, f)
		}
	}

	return fields
}

func uniqKey(rt reflect.Type) string {
	return rt.PkgPath() + "|" + rt.Name()
}

// collectStructComment collect struct comment
func collectStructComment(refType reflect.Type, s *Struct) error {
	// 解析-获取结构体注释
	var r map[string]string
	var f map[string]string
	var err error
	if r, f, err = resolve(s.Name); err != nil {
		return fmt.Errorf("resolve output failed, error is %v", err)
	}
	s.Comment = r[commentKey]
	s.Description = r[descriptionKey]

	// 内嵌结构体
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)

		sf := Field{
			StructField: field,
			Comment:     f[field.Name],
		}

		fieldType := field.Type
		if field.Anonymous { // 匿名
			// 忽略匿名接口
			if fieldType.Kind() != reflect.Interface {
				sf.Struct, err = ResolveStruct(fieldType)
				if err != nil {
					return err
				}
			}
		}
		// 非匿名结构体类型
		if fieldType.Kind() == reflect.Ptr ||
			fieldType.Kind() == reflect.Slice ||
			fieldType.Kind() == reflect.Map ||
			fieldType.Kind() == reflect.Chan ||
			fieldType.Kind() == reflect.Array {
			fieldType = fieldType.Elem()
		}
		// 忽略time.Time
		if fieldType.Kind() == reflect.Struct && fieldType != reflect.TypeOf((*time.Time)(nil)).Elem() {
			// 字段类型元素包含本类型
			isSelfType := uniqKey(fieldType) == uniqKey(refType)
			if !isSelfType {
				sf.Struct, err = ResolveStruct(fieldType)
				if err != nil {
					return err
				}
			}
		}

		s.Fields = append(s.Fields, sf)
	}

	return nil
}

const (
	structStart    = "type"
	structEnd      = "}"
	fieldSep       = " "
	commentSep     = "//"
	commentKey     = "comment"
	descriptionKey = "description"
)

var (
	structCommentCache = make(map[string]StructCommentEntity)
)

type StructCommentEntity struct {
	StructName    string
	StructComment map[string]string
	FieldComment  map[string]string
}

func resolve(structName string) (map[string]string, map[string]string, error) {
	return resolveWithParser(structName)
}

func resolveWithParser(structName string) (map[string]string, map[string]string, error) {
	var structCommentMap = make(map[string]string)
	var fieldCommentMap = make(map[string]string)

	if ent, ok := structCommentCache[structName]; ok {
		return ent.StructComment, ent.FieldComment, nil
	}

	ip := &parser.ImportPath{}
	path, err := ip.GetByCurrentDir()
	if err != nil {
		return structCommentMap, fieldCommentMap, err
	}

	name := structName
	dotIndex := strings.LastIndex(structName, ".")
	if dotIndex != -1 {
		if dotIndex != 0 {
			path = structName[:dotIndex]
		}
		name = structName[dotIndex+1:]
	}

	parserObj := parser.NewParser(parser.Option{})
	pkg, err := parserObj.ParseByGoPackages(path)
	if err != nil {
		return structCommentMap, fieldCommentMap, err
	}
	structs := make([]parser.Struct, 0, len(pkg.Pkgs))
	for _, pkg := range pkg.Pkgs {
		structs = append(structs, pkg.Structs...)
	}

	var exist bool
	for _, oneStruct := range structs {
		// 缓存
		var tmpStructCommentMap = make(map[string]string)
		var tmpFieldCommentMap = make(map[string]string)
		tmpStructCommentMap[commentKey] = strings.TrimSpace(oneStruct.Comment)
		tmpStructCommentMap[descriptionKey] = strings.TrimSpace(oneStruct.Doc)
		for _, field := range oneStruct.Fields {
			tmpFieldCommentMap[field.Name] = strings.TrimSpace(field.Comment)
		}
		structName := path + "." + oneStruct.Name
		structCommentCache[structName] = StructCommentEntity{
			StructName:    structName,
			StructComment: tmpStructCommentMap,
			FieldComment:  tmpFieldCommentMap,
		}

		if oneStruct.Name != name {
			continue
		}
		exist = true
		structCommentMap, fieldCommentMap = tmpStructCommentMap, tmpFieldCommentMap
	}
	_ = exist
	// if !exist {
	// 	log.Printf("Can't find comment info of %s", structName)
	// }

	return structCommentMap, fieldCommentMap, nil
}
