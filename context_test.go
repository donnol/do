package do

import (
	"context"
	"database/sql"
	"testing"
)

func TestContextHelper(t *testing.T) {
	{
		type nameck struct{}
		ctx := context.Background()
		helper := NewContextHelper[nameck, string]()
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
			helper := NewContextHelper[nameck, string]()
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
			helper := NewContextHelper[novalueck, string]()
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
		helper := NewContextHelper[dbconnck, *sql.Conn]()
		ctx = helper.WithValue(ctx, &sql.Conn{})
		v := helper.MustValue(ctx)
		if v == nil {
			t.Fatalf("bad case: %+v is nil", v)
		}
	}
}
