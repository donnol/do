package do

import (
	"fmt"
	"runtime"
)

type AssertHandler interface {
	Errorf(format string, args ...any)
}

func Assert[T comparable](handler AssertHandler, l, r T, msgAndArgs ...any) {
	if Equal(l, r) {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		handler.Errorf("[%s:%d] Bad case, %v != %v, %s", file, line, l, r, messageFromMsgAndArgs(msgAndArgs...))
	} else {
		handler.Errorf("Bad case, %v != %v, %s", l, r, messageFromMsgAndArgs(msgAndArgs...))
	}
}

func AssertSlice[T comparable](handler AssertHandler, lslice, rslice []T, msgAndArgs ...any) {
	ll, rl := len(lslice), len(rslice)
	if ll != rl {
		handler.Errorf("Bad case, left length(%d) != right length(%d)", ll, rl)
		return
	}

	for i := 0; i < ll; i++ {
		if !Equal(lslice[i], rslice[i]) {
			handler.Errorf("Bad case, No.%d: %v != %v, %s", i, lslice[i], rslice[i], messageFromMsgAndArgs(msgAndArgs...))
		}
	}
}

func AssertSlicePtr[T comparable](handler AssertHandler, lslice, rslice []*T, msgAndArgs ...any) {
	ll, rl := len(lslice), len(rslice)
	if ll != rl {
		handler.Errorf("Bad case, left length(%d) != right length(%d)", ll, rl)
		return
	}

	for i := 0; i < ll; i++ {
		if !Equal(*lslice[i], *rslice[i]) {
			handler.Errorf("Bad case, No.%d: %v != %v, %s", i, lslice[i], rslice[i], messageFromMsgAndArgs(msgAndArgs...))
		}
	}
}

func messageFromMsgAndArgs(msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}

	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}

	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}

	return ""
}
