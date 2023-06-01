package main

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/donnol/do"
)

func main() {
	// 初始化，建立接口依赖关系
	s := NewBookService(NewUserStore(), NewBookStore())

	// 运行，执行业务逻辑
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		// 每个业务执行前生成唯一id
		// 使用ctx传递traceId
		traceId := strconv.Itoa(rand.Int())
		atx := context.WithValue(ctx, do.TraceKey{}, traceId)

		book := s.WithUser(atx, uint64(i))
		if book.UserId != 1 || book.UserName != "jd" {
			panic("bad user info")
		}
	}
}

type bookServiceImpl struct {
	us UserStore
	bs BookStore
}

func NewBookService(
	us UserStore,
	bs BookStore,
) BookService {
	b, _ := GetBookServiceProxy(&bookServiceImpl{
		us: us,
		bs: bs,
	})
	return b
}

func (impl *bookServiceImpl) WithUser(ctx context.Context, id uint64) Book {
	book := impl.bs.ById(ctx, id)

	user := impl.us.ById(ctx, book.UserId)
	book.UserName = user.Name

	return book
}

type userStoreImpl struct {
}

func NewUserStore() UserStore {
	s, _ := GetUserStoreProxy(&userStoreImpl{})
	return s
}

func (impl *userStoreImpl) ById(ctx context.Context, id uint64) User {
	return User{
		Id:   id,
		Name: "jd",
	}
}

type bookStoreImpl struct {
}

func NewBookStore() BookStore {
	s, _ := GetBookStoreProxy(&bookStoreImpl{})
	return s
}

func (impl *bookStoreImpl) ById(ctx context.Context, id uint64) Book {
	return Book{
		Id:     id,
		Name:   "Good or bad",
		UserId: 1,
	}
}

func (impl *bookStoreImpl) Hook() {

}
