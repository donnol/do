package do

import (
	"context"
	"database/sql"
	"testing"
	"unsafe"
)

func TestContextHelper(t *testing.T) {
	{
		type nameck struct{}
		ctx := context.Background()
		helper := ContextHelper[nameck, string]{}
		{
			ksize, ssize := unsafe.Sizeof(nameck{}), unsafe.Sizeof(helper)
			if ksize != 0 || ssize != 0 {
				t.Fatalf("bad size of var, ksize: %d, ssize: %d", ksize, ssize)
			}
		}
		ctx = helper.WithValue(ctx, "1")
		v := helper.MustValue(ctx)
		if v != "1" {
			t.Fatalf("bad case: v != %s", "1")
		}

		// replace value
		{
			ctx = helper.WithValue(ctx, "2")
			v := helper.MustValue(ctx)
			if v != "2" {
				t.Fatalf("bad case: v != %s", "2")
			}
		}

		// other helper will replace value too, because the k always is `nameck{}`
		{
			// ctx := context.Background()
			helper := ContextHelper[nameck, string]{}
			ctx = helper.WithValue(ctx, "3")
			v := helper.MustValue(ctx)
			if v != "3" {
				t.Fatalf("bad case: v != %s", "3")
			}
		}

		// can't find
		func() {
			defer func() {
				if v := recover(); v != nil {
					if v.(error).Error() != "context can't find value of do.novalueck{}" {
						t.Fatalf("bad case of value not exist: %s", v)
					}
				}
			}()
			type novalueck struct{}
			helper := ContextHelper[novalueck, string]{}
			v, ok := helper.Value(ctx)
			if ok || v != "" {
				t.Fatalf("bad case, invalid v or ok")
			}
			helper.MustValue(ctx)
		}()
	}

	{
		type dbconnck struct{}
		ctx := context.Background()
		helper := ContextHelper[dbconnck, *sql.Conn]{}
		ctx = helper.WithValue(ctx, &sql.Conn{})
		v := helper.MustValue(ctx)
		if v == nil {
			t.Fatalf("bad case: %+v is nil", v)
		}
	}
}
