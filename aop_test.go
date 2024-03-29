package do

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

type tracerImpl struct {
	pctx  ProxyContext
	no    int
	begin time.Time
}

var i int

func (impl *tracerImpl) New(pctx ProxyContext, extras ...any) Tracer {
	i++
	return &tracerImpl{
		pctx: pctx,
		no:   i,
	}
}

func (impl *tracerImpl) Begin() {
	impl.begin = time.Now()
}

func (impl *tracerImpl) Stop(args ...any) {
	log.Output(3, fmt.Sprintf("[%s] NO.%d: used time %v, args: %+v\n", impl.pctx, impl.no, time.Since(impl.begin), args))
}

func TestProxyTracer(t *testing.T) {
	gtracers = make(tracers, 0)

	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	log.SetFlags(log.Llongfile | log.LstdFlags)

	RegisterProxyTracer(&tracerImpl{})
	RegisterProxyTracer(&tracerImpl{})

	tpctx := ProxyContext{
		PkgPath:       "testpkg",
		InterfaceName: "testinter",
		MethodName:    "testmethod",
	}
	func() {
		stop := ProxyTraceBegin(tpctx)
		defer stop("test result")
	}()

	func() {
		stop := ProxyTraceBegin(tpctx)
		defer stop("test result")
	}()

	// assert
	{
		output := buf.String()
		parts := strings.Split(output, "\n")
		for index, part := range parts {
			var want string
			switch index {
			case 0:
				want = "NO.2"
			case 1:
				want = "NO.1"
			case 2:
				want = "NO.4"
			case 3:
				want = "NO.3"
			}
			if !strings.Contains(part, want) {
				t.Errorf("bad case: %s don't contains %s", part, want)
			}
			if part != "" {
				wantr := "args: [test result]"
				if !strings.Contains(part, wantr) {
					t.Errorf("bad case: %s don't contains %s", part, wantr)
				}
			}
		}
	}
}
