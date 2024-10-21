package do

import (
	"encoding/json"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type (
	Entity[U any] struct {
		Id   U      `json:"id"`   // id
		Name string `json:"name"` // 名称
	}
	TreeEntity[U any] struct {
		Entity[U]

		Childs []Entity[U] `json:"childs"` // 子元素
	}
)

type (
	Pager struct {
		Page     int `json:"page" form:"page"`         // page No.
		PageSize int `json:"pageSize" form:"pageSize"` // page size
	}
	PageCond[T, C any] struct {
		Table T

		Cond C
	}

	PageResult[T any] struct {
		Total int64 `json:"total"` // total

		ListResult[T]
	}
	ListResult[T any] struct {
		List []T `json:"list"` // list
	}
)

func PageCondFrom[T, C any](t T, c C) (p PageCond[T, C]) {
	return PageCond[T, C]{
		Table: t,
		Cond:  c,
	}
}

func (p PageCond[T, C]) Unpack() (T, C) {
	return p.Table, p.Cond
}

func (p Pager) Limit() int {
	return p.PageSize
}

func (p Pager) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type (
	IdType[T any] struct {
		Id T `json:"id" form:"id"` // id
	}
	IdsType[T any] struct {
		Ids []T `json:"ids" form:"ids"` // ids
	}
)

// Id to/from string when json encode/decode
type Id uint64

var (
	_ json.Marshaler   = (*Id)(nil)
	_ json.Unmarshaler = (*Id)(nil)
)

func (id Id) MarshalJSON() ([]byte, error) {
	s := strconv.FormatUint(uint64(id), 10)
	s = strconv.Quote(s)

	return StringToBytes(s), nil
}

func (id *Id) UnmarshalJSON(data []byte) error {
	s := BytesToString(data)
	s, err := strconv.Unquote(s)
	if err != nil {
		return fmt.Errorf("[id] data unquote failed: %v of %s", err, data)
	}

	if s == "" {
		return nil
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fmt.Errorf("[id] data parse uint failed: %v of %s", err, data)
	}
	*id = Id(i)

	return nil
}

type Password string

// Encrypt 使用bcrypt算法将明文密码哈希得到hash字符串
// bcrypt算法在对同一个密码哈希多次会得出不同结果，极大的保证了用户密码的安全
func (p Password) Encrypt() (pp string, err error) {
	r, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	return string(r), nil
}

// Compare 使用bcrypt算法判断密码是否与传入hash值匹配
func (p Password) Compare(hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
}

// Password 实现String方法，打印时自动替换为*
func (p Password) String() string {
	return "*"
}
