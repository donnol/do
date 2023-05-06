package do

import (
	"bytes"
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

func (impl *tracerImpl) New(pctx ProxyContext) Tracer {
	i++
	return &tracerImpl{
		no: i,
	}
}

func (impl *tracerImpl) Begin() {
	impl.begin = time.Now()
}

func (impl *tracerImpl) Stop() {
	log.Printf("[%s] NO.%d: used time %v\n", impl.pctx, impl.no, time.Since(impl.begin))
}

func TestProxyTracer(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	RegisterProxyTracer(&tracerImpl{})
	RegisterProxyTracer(&tracerImpl{})

	tpctx := ProxyContext{
		PkgPath:       "test",
		InterfaceName: "test",
		MethodName:    "test",
	}
	func() {
		stop := ProxyTraceBegin(tpctx)
		defer stop()
	}()

	func() {
		stop := ProxyTraceBegin(tpctx)
		defer stop()
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
		}
	}
}
