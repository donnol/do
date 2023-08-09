package do

import "fmt"

// nameof return v's name
// use `go generate` to run `letgo nameof` on project root directory, it will generate the name info of every v passed to nameof.
func nameof(v any) string {
	// How to know which type v has?
	k := fmt.Sprintf("%v", v)

	// We can't know the name of a variable through reflect.

	// 除非替换代码：在遍历代码时是可以拿到变量名及其类型名的，拿到后直接替换源码
	// 但是这样做，侵入性太强，不太好。
	// 编译的时候修改ast呢？这样就不会修改源文件。
	// 编译前执行的钩子
	// 或者在Makefile里添加上执行文本替换的指令，该指令执行时复制一份源码然后替换，最后再编译，编译好之后删除复制的源码
	// tmpbuild=/tmp/build/
	// build:
	//     cp -r . $(tmpbuild) && cd $(tmpbuild) && letgo nameof
	//     go build $(tmpbuild)
	//     rm -rf $(tmpbuild)

	return nameofStore[k]
}

// map's key is struct name + field name
var nameofStore = make(map[string]string)
