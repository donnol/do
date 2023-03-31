package do

import (
	"reflect"
	"testing"
)

type User struct {
	Id   uint64
	Name string
}

type Book struct {
	Id     uint64
	Title  string
	Author uint64 // User.Id
}

type BookWithUser struct {
	Book
	UserName string
}

func UserBookMatcher(j Book, k User) bool {
	return j.Author == k.Id
}

func TestNestedJoin(t *testing.T) {
	r := NestedJoin([]Book{
		{Id: 1, Title: "hello", Author: 1},
		{Id: 2, Title: "world", Author: 1},
		{Id: 3, Title: "good", Author: 2},
		{Id: 4, Title: "job", Author: 2},
	}, []User{
		{Id: 1, Name: "jd"},
		{Id: 2, Name: "jc"},
	}, UserBookMatcher, func(j Book, k User) BookWithUser {
		return BookWithUser{
			Book:     j,
			UserName: k.Name,
		}
	})
	want := []BookWithUser{
		{Book{1, "hello", 1}, "jd"},
		{Book{2, "world", 1}, "jd"},
		{Book{3, "good", 2}, "jc"},
		{Book{4, "job", 2}, "jc"},
	}
	if !reflect.DeepEqual(r, want) {
		t.Errorf("bad case, %+v != %+v", r, want)
	}
}

func TestHashJoin(t *testing.T) {
	r := HashJoin([]Book{
		{Id: 1, Title: "hello", Author: 1},
		{Id: 2, Title: "world", Author: 1},
		{Id: 3, Title: "good", Author: 2},
		{Id: 4, Title: "job", Author: 2},
	}, []User{
		{Id: 1, Name: "jd"},
		{Id: 2, Name: "jc"},
	}, func(item Book) uint64 {
		return item.Author
	}, func(item User) uint64 {
		return item.Id
	}, func(j Book, k User) BookWithUser {
		return BookWithUser{
			Book:     j,
			UserName: k.Name,
		}
	})
	want := []BookWithUser{
		{Book{1, "hello", 1}, "jd"},
		{Book{2, "world", 1}, "jd"},
		{Book{3, "good", 2}, "jc"},
		{Book{4, "job", 2}, "jc"},
	}
	if !reflect.DeepEqual(r, want) {
		t.Errorf("bad case, %+v != %+v", r, want)
	}
}

func BenchmarkNestedJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NestedJoin([]Book{
			{Id: 1, Title: "hello", Author: 1},
			{Id: 2, Title: "world", Author: 1},
			{Id: 3, Title: "good", Author: 2},
			{Id: 4, Title: "job", Author: 2},
		}, []User{
			{Id: 1, Name: "jd"},
			{Id: 2, Name: "jc"},
		}, UserBookMatcher, func(j Book, k User) BookWithUser {
			return BookWithUser{
				Book:     j,
				UserName: k.Name,
			}
		})
	}
}

func BenchmarkNestedJoin100X(b *testing.B) {
	users := []User{
		{Id: 1, Name: "jd"},
		{Id: 2, Name: "jc"},
	}
	books := []Book{
		{Id: 1, Title: "hello", Author: 1},
		{Id: 2, Title: "world", Author: 1},
		{Id: 3, Title: "good", Author: 2},
		{Id: 4, Title: "job", Author: 2},
	}
	for i := 0; i < 100; i++ {
		users = append(users, User{
			Id: 2, Name: "jc",
		})
		books = append(books, Book{
			Id: 4, Title: "job", Author: 2,
		})
	}

	for i := 0; i < b.N; i++ {
		NestedJoin(books, users, UserBookMatcher, func(j Book, k User) BookWithUser {
			return BookWithUser{
				Book:     j,
				UserName: k.Name,
			}
		})
	}
}

func BenchmarkHashJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashJoin([]Book{
			{Id: 1, Title: "hello", Author: 1},
			{Id: 2, Title: "world", Author: 1},
			{Id: 3, Title: "good", Author: 2},
			{Id: 4, Title: "job", Author: 2},
		}, []User{
			{Id: 1, Name: "jd"},
			{Id: 2, Name: "jc"},
		}, func(item Book) uint64 {
			return item.Author
		}, func(item User) uint64 {
			return item.Id
		}, func(j Book, k User) BookWithUser {
			return BookWithUser{
				Book:     j,
				UserName: k.Name,
			}
		})
	}
}

var (
	books []Book
	users []User
)

func TestMain(m *testing.M) {
	users = []User{
		{Id: 1, Name: "jd"},
		{Id: 2, Name: "jc"},
	}
	books = []Book{
		{Id: 1, Title: "hello", Author: 1},
		{Id: 2, Title: "world", Author: 1},
		{Id: 3, Title: "good", Author: 2},
		{Id: 4, Title: "job", Author: 2},
	}
	for i := 0; i < 1_000_000; i++ {
		users = append(users, User{
			Id: 2, Name: "jc",
		})
		books = append(books, Book{
			Id: 4, Title: "job", Author: 2,
		})
	}

	m.Run()
}

func BenchmarkHashJoin1000000X(b *testing.B) {

	for i := 0; i < b.N; i++ {
		HashJoin(books, users, func(item Book) uint64 {
			return item.Author
		}, func(item User) uint64 {
			return item.Id
		}, func(j Book, k User) BookWithUser {
			return BookWithUser{
				Book:     j,
				UserName: k.Name,
			}
		})
	}
}
