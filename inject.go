package do

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

// Ioc 控制反转，Inversion of Control
type Ioc struct {
	enableUnexportedFieldSetValue bool // 开启对非导出字段的赋值
	print                         bool

	providerMap map[reflect.Type]typeInfo
	cache       map[reflect.Type]reflect.Value
}

type IocOption struct {
	EnableUnexportedFieldSetValue bool // 开启对非导出字段的赋值
	Print                         bool // 打印初始话的对象
}

type typeInfo struct {
	depType  []reflect.Type
	provider reflect.Value
}

func NewIoc(
	opt *IocOption,
) *Ioc {
	if opt == nil {
		opt = &IocOption{}
	}
	return &Ioc{
		enableUnexportedFieldSetValue: opt.EnableUnexportedFieldSetValue,
		print:                         opt.Print,
		providerMap:                   make(map[reflect.Type]typeInfo),
		cache:                         make(map[reflect.Type]reflect.Value),
	}
}

// RegisterProvider 注册provider函数，形如`func New(fielda TypeA, fieldb TypeB) (T)`
func (ioc *Ioc) RegisterProvider(v any) (err error) {
	refValue := reflect.ValueOf(v)
	refType := refValue.Type()
	if refType.Kind() != reflect.Func {
		return fmt.Errorf("please input func")
	}

	// 分析函数的参数和返回值
	ti := typeInfo{
		depType:  make([]reflect.Type, 0, refType.NumIn()),
		provider: refValue,
	}
	for i := 0; i < refType.NumIn(); i++ {
		in := refType.In(i)
		ti.depType = append(ti.depType, in)
	}
	// 返回：instance
	min := 1
	if refType.NumOut() == 0 {
		return fmt.Errorf("can't find result in func")
	}
	if refType.NumOut() < min {
		return fmt.Errorf("too little result in func, min is %d", min)
	}
	for i := 0; i < refType.NumOut(); i++ {
		out := refType.Out(i)
		ioc.providerMap[out] = ti
	}

	return
}

// Inject 为对象注入依赖，传入结构体指针，根据结构体的字段类型找到对应的provider，执行后将获得的值赋予字段
//
// 如果provider需要参数，则根据参数类型继续找寻相应的provider，直至初始化完成
func (ioc *Ioc) Inject(v any) (err error) {
	refValue := reflect.ValueOf(v)
	refType := refValue.Type()
	if refType.Kind() != reflect.Ptr {
		return fmt.Errorf("v is not a pointer")
	}
	eleValue := refValue.Elem()
	eleType := eleValue.Type()
	if eleType.Kind() != reflect.Struct {
		return fmt.Errorf("v is not a struct")
	}

	ioc.levelPrint(0, "will inject object of type(%s)\n", eleType.Name())

	// 遍历field
	for i := 0; i < eleValue.NumField(); i++ {
		sf := eleType.Field(i)
		field := eleValue.Field(i)
		fieldType := field.Type()

		var level = 1
		ioc.levelPrint(level, "will set field %s(%s)\n", sf.Name, fieldType.Name())

		// 根据类型查找值
		var value reflect.Value
		value, err = ioc.find(fieldType, level+1)
		if err != nil {
			return
		}

		// 给字段赋值
		if ioc.enableUnexportedFieldSetValue {
			rf := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
			rf.Set(value)
		} else {
			field.Set(value)
		}

		ioc.levelPrint(level, "finish set field %s(%s)\n", sf.Name, fieldType.Name())
	}

	ioc.levelPrint(0, "finish inject object of type(%s)\n", eleType.Name())

	return
}

func (ioc *Ioc) levelPrint(l int, format string, args ...any) {
	if ioc.print {
		prefix := "[inject] "
		level := ""
		for i := 0; i < l; i++ {
			level += "    "
		}
		log.Printf(prefix+level+format, args...)
	}
}

var (
	emptyStruct         = reflect.TypeOf((*struct{})(nil))
	emptyStructValue    = reflect.New(emptyStruct.Elem()).Elem()
	emptyStructPtrValue = reflect.New(emptyStruct).Elem()

	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

func (ioc *Ioc) find(typ reflect.Type, level int) (r reflect.Value, err error) {
	ioc.levelPrint(level, "will get field type %s's value\n", typ.Name())

	value, ok := ioc.cache[typ]
	if ok {
		ioc.levelPrint(level, "finish get field type %s's value from cache\n", typ.Name())
		return value, nil
	}

	// 在provider里寻找初始化函数
	provider, ok := ioc.providerMap[typ]
	if !ok {
		// 检查类型是否是struct{}
		if typ.ConvertibleTo(emptyStruct.Elem()) {
			ioc.levelPrint(level, "finish get field type %s's value from empty struct\n", typ.Name())
			return emptyStructValue, nil
		}
		if typ.ConvertibleTo(emptyStruct) {
			ioc.levelPrint(level, "finish get field type %s's value from empty struct pointer\n", typ.Name())
			return emptyStructPtrValue, nil
		}
		return r, fmt.Errorf("can't find provider of %+v", typ)
	}

	// 调用
	in := make([]reflect.Value, 0, len(provider.depType))
	for _, dep := range provider.depType {
		// 在已有provider里找指定类型
		var value reflect.Value
		if value, ok = ioc.cache[dep]; !ok {
			value, err = ioc.find(dep, level+1)
			if err != nil {
				return r, err
			}
			ioc.cache[dep] = value
		}

		in = append(in, value)
	}
	newValues := provider.provider.Call(in)
	if len(newValues) == 0 {
		return r, fmt.Errorf("can't get new value by provider")
	}

	// 返回值里，第一个必须是实例，最后一个必须是error，中间的忽略
	newValue := newValues[0]
	ioc.cache[typ] = newValue

	if len(newValues) > 1 {
		lastValue := newValues[len(newValues)-1]
		last := lastValue.Interface()
		if lastValue.Type().Implements(errorType) {
			if v, ok := last.(error); ok {
				if v != nil {
					return r, fmt.Errorf("call failed, err is %+v", v)
				}
			}
		} else {
			return r, fmt.Errorf("last return value is not error, is %+v", last)
		}
	}

	ioc.levelPrint(level, "finish get field type %s's value from provider `%s`\n", typ.Name(), provider.provider.Type().String())

	return newValue, nil
}
