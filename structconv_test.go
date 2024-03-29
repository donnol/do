package do_test

import (
	"reflect"
	"testing"

	"github.com/donnol/do"
)

// 虽然通过名字来赋值，省下了功夫，但现实中字段对应往往不是简单的通过名字就能实现，还有其它因素，所以实用性并不太高
type (
	UserReq struct {
		Phone string
	}

	UserResp struct {
		Name string
		Age  uint
	}

	UserTable struct {
		Id    uint
		Name  string
		Age   uint
		Phone string
	}

	Article struct {
		Name     string
		UserName string
	}
)

type (
	UserEmbedReq struct {
		UserReq

		Addr []string
	}
	UserEmbed struct {
		UserTable

		Addr []string
	}
	UserEmbedTwice struct {
		UserEmbed

		Like []string
	}
)

func TestConvByFieldName(t *testing.T) {
	from := UserReq{
		Phone: "12345678901",
	}
	to := &UserTable{}
	do.ConvByName(from, to)

	if to.Phone != from.Phone {
		t.Fatalf("converse failed: %s != %s\n", to.Phone, from.Phone)
	}

	{
		from := UserEmbedReq{
			UserReq: UserReq{
				Phone: "123",
			},
		}
		to := &UserEmbed{}
		do.ConvByName(from, to)
		if to.Phone == "" || to.Phone != from.Phone {
			t.Fatalf("converse failed: %s != %s\n", to.Phone, from.Phone)
		}
	}

	{
		from := UserReq{
			Phone: "123",
		}
		to := &UserEmbed{}
		do.ConvByName(from, to)
		if to.Phone == "" || to.Phone != from.Phone {
			t.Fatalf("converse failed: %s != %s\n", to.Phone, from.Phone)
		}
	}

	{
		from := UserReq{
			Phone: "123",
		}
		to := &UserEmbedTwice{}
		do.ConvByName(from, to)
		if to.Phone == "" || to.Phone != from.Phone {
			t.Fatalf("converse failed: %s != %s\n", to.Phone, from.Phone)
		}
	}

	to.Id = 1
	to.Name = "jd"
	to.Age = 18

	to2 := &UserResp{}
	do.ConvByName(to, to2)

	if to2.Name == "" || to2.Name != to.Name {
		t.Fatalf("converse failed: %s != %s\n", to2.Name, to.Name)
	}
	if to2.Age == 0 || to2.Age != to.Age {
		t.Fatalf("converse failed: %d != %d\n", to2.Age, to.Age)
	}

	{
		from := UserEmbedReq{
			Addr: []string{"luowei", "haidao"},
		}
		to := &UserEmbed{}
		do.ConvByName(from, to)
		if len(to.Addr) != len(from.Addr) {
			t.Fatalf("conv failed of length %v != %v", len(to.Addr), len(from.Addr))
		}
		if !reflect.DeepEqual(to.Addr, from.Addr) {
			t.Fatalf("conv failed of not equal %v != %v", to.Addr, from.Addr)
		}
	}
}

func TestConvListByFieldName(t *testing.T) {
	type args struct {
		from []UserReq
		to   []*UserTable
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				from: []UserReq{
					{
						Phone: "12345678901",
					},
					{
						Phone: "234",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.to = do.MakeSlice[UserTable](len(tt.args.from))

			do.ConvSliceByName(tt.args.from, tt.args.to)

			for i := range tt.args.from {
				if tt.args.to[i].Phone == "" || tt.args.to[i].Phone != tt.args.from[i].Phone {
					t.Fatalf("converse failed: %s != %s\n", tt.args.to[i].Phone, tt.args.from[i].Phone)
				}
			}
		})
	}
}
