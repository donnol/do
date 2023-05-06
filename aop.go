package do

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// 每个包、每个接口、每个方法唯一对应一个方法
type ProxyContext struct {
	PkgPath       string
	InterfaceName string
	MethodName    string
}

func (pctx ProxyContext) String() string {
	return fmt.Sprintf(pctx.bracket("PkgPath: %s InterfaceName: %s MethodName: %s"), pctx.PkgPath, pctx.InterfaceName, pctx.MethodName)
}

func (pctx ProxyContext) Uniq() string {
	return pctx.PkgPath + "|" + pctx.InterfaceName + "|" + pctx.MethodName
}

func (pctx ProxyContext) IsEmpty() bool {
	return pctx.PkgPath == "" && pctx.InterfaceName == "" && pctx.MethodName == ""
}

func (pctx ProxyContext) Logf(format string, args ...any) {
	pctx.logf(pctx.String()+": "+format, args...)
}

func (pctx ProxyContext) LogShortf(format string, args ...any) {
	pctx.logf(pctx.bracket(pctx.MethodName)+": "+format, args...)
}

func (pctx ProxyContext) logf(format string, args ...any) {
	err := log.Output(3, fmt.Sprintf(format, args...))
	if err != nil {
		fmt.Printf("Output failed: %+v\n", err)
	}
}

func (pctx ProxyContext) bracket(s string) string {
	return "[" + s + "]"
}

// ProxyCtxFunc
// 对于method: func(string, int) (int, error)
// f := method.(func(string, int) (int, error))
// a1 := args[0].(string)
// a2 := args[1].(int)
// r1, r2 := f(a1, a2)
// res = append(res, r1, r2)
type ProxyCtxFunc func(ctx ProxyContext, method any, args []any) (res []any)

type ProxyCtxFuncStore struct {
	m map[string]ProxyCtxFunc
}

func NewProxyCtxMap() *ProxyCtxFuncStore {
	return &ProxyCtxFuncStore{
		m: make(map[string]ProxyCtxFunc),
	}
}

func (m *ProxyCtxFuncStore) Lookup(pctx ProxyContext, typeParams ...string) (ProxyCtxFunc, bool) {
	key := pctxUniqWithTypeParams(pctx.Uniq(), typeParams...)
	v, ok := m.m[key]
	return v, ok
}

func (m *ProxyCtxFuncStore) Set(pctx ProxyContext, f ProxyCtxFunc, typeParams ...string) {
	key := pctxUniqWithTypeParams(pctx.Uniq(), typeParams...)
	m.m[key] = f
}

func pctxUniqWithTypeParams(uniq string, typeParams ...string) string {
	if len(typeParams) > 0 {
		uniq = uniq + "|" + strings.Join(typeParams, ",")
	}
	return uniq
}

type Tracer interface {
	New(ProxyContext) Tracer // 新建Tracer，每个方法调用均新建一个
	Begin()
	Stop()
}

type tracers []Tracer

var (
	gtracers tracers
)

func RegisterProxyTracer(tracers ...Tracer) {
	gtracers = append(gtracers, tracers...)
}

// ProxyTraceBegin LIFO
func ProxyTraceBegin(pctx ProxyContext) (stop func()) {
	stops := make([]func(), 0, len(gtracers))
	for _, tc := range gtracers {
		o := tc.New(pctx)
		o.Begin()
		stops = append(stops, o.Stop)
	}
	stop = func() {
		for i := len(stops) - 1; i >= 0; i-- {
			so := stops[i]
			so()
		}
	}
	return
}

type TimeTracer struct {
	pctx  ProxyContext
	begin time.Time
}

func (impl *TimeTracer) New(pctx ProxyContext) Tracer {
	return &TimeTracer{
		pctx: pctx,
	}
}

func (impl *TimeTracer) Begin() {
	impl.begin = time.Now()
}

func (impl *TimeTracer) Stop() {
	log.Printf("[%s] used time %v\n", impl.pctx, time.Since(impl.begin))
}
