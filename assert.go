package do

import "fmt"

type AssertHandler interface {
	Errorf(format string, args ...any)
}

func Assert[T comparable](handler AssertHandler, l, r T, msgAndArgs ...any) {
	if Equal(l, r) {
		return
	}

	handler.Errorf("Bad case, %v != %v, %s", l, r, messageFromMsgAndArgs(msgAndArgs...))
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
