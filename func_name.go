package do

import (
	"runtime"
	"strconv"
)

func FuncName(skip int, withFileInfo bool) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	fun := runtime.FuncForPC(pc)
	if fun == nil {
		return ""
	}
	if withFileInfo {
		return file + ":" + strconv.Itoa(line) + " " + fun.Name()
	}
	return fun.Name()
}
