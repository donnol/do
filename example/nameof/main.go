package main

func main() {
	// 变量
	var a int
	if nameof(a) != "a" {
		panic("bad name of a")
	}
	if nameof(&a) != "a" {
		panic("bad name of a")
	}

	// 结构体
	if nameof(User{}) != "User" {
		panic("bad name of User")
	}
	if nameof(&User{}) != "User" {
		panic("bad name of User")
	}

	// 结构体字段
	var u User
	if nameof(u.Name) != "Name" {
		panic("bad name of u.Name")
	}
	if nameof(u.age) != "age" {
		panic("bad age of u.age")
	}
}

type User struct {
	Name string
	age  int
}

// nameof is a marker func
func nameof(v any) string {
	return ""
}
