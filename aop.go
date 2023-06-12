package do

import (
	"context"
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

var (
	globalProxyCtxMap = NewProxyCtxMap()
)

// RegisterProxyMethod 注册代理方法，根据包名+接口名+方法名唯一对应一个方法；在有了泛型后还要加上类型参数，唯一键变为包名+接口名+方法名+TP1,TP2,...
func RegisterProxyMethod(pctx ProxyContext, cf ProxyCtxFunc, typeParams ...string) {
	globalProxyCtxMap.Set(pctx, cf, typeParams...)
}

func GlobalProxyCtxMap() *ProxyCtxFuncStore {
	return globalProxyCtxMap
}

type Tracer interface {
	New(pctx ProxyContext, extras ...any) Tracer // 新建Tracer，每个方法调用均新建一个
	Begin()
	Stop()
}

type tracers []Tracer

var (
	gtracers = tracers{
		&TimeTracer{},
	}
)

func RegisterProxyTracer(tracers ...Tracer) {
	gtracers = append(gtracers, tracers...)
}

// ProxyTraceBegin LIFO
func ProxyTraceBegin(pctx ProxyContext, extras ...any) (stop func()) {
	stops := make([]func(), 0, len(gtracers))
	for _, tc := range gtracers {
		o := tc.New(pctx, extras...)
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
	pctx    ProxyContext
	traceId TraceId
	begin   time.Time
}

type (
	TraceKey struct{}
	TraceId  string
)

func (impl *TimeTracer) New(pctx ProxyContext, extras ...any) Tracer {
	traceId := parseExtra(extras...)

	return &TimeTracer{
		pctx:    pctx,
		traceId: traceId,
	}
}

func (impl *TimeTracer) Begin() {
	impl.begin = time.Now()
}

func (impl *TimeTracer) Stop() {
	log.Output(3, fmt.Sprintf("[%s] |%s| used time %v\n", impl.pctx, impl.traceId, time.Since(impl.begin)))
}

func parseExtra(extras ...any) (traceId TraceId) {
	if len(extras) == 0 {
		return
	}

	ctx, ok := extras[0].(context.Context)
	if ok {
		v := ctx.Value(TraceKey{})
		if vv, ok := v.(string); ok {
			traceId = TraceId(vv)
		} else if vv, ok := v.(TraceId); ok {
			traceId = vv
		}
	}

	return
}
