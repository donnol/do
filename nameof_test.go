package do

import (
	"go/ast"
	"go/parser"
	"testing"

	"golang.org/x/tools/go/ast/astutil"
)

func Test_nameof(t *testing.T) {
	t.SkipNow()

	type User struct {
		Name string
	}
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "struct",
			args: args{
				v: User{},
			},
			want: "User",
		},
		{
			name: "struct field",
			args: args{
				v: (User{}).Name,
			},
			want: "Name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nameof(tt.args.v); got != tt.want {
				t.Errorf("nameof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Types_name(t *testing.T) {
	expr, err := parser.ParseExpr(`name == u.Name`)
	if err != nil {
		t.Fatal(err)
	}

	// 0  *ast.BinaryExpr {
	// 	1  .  X: *ast.Ident {
	// 	2  .  .  NamePos: -
	// 	3  .  .  Name: "name"
	// 	4  .  }
	// 	5  .  OpPos: -
	// 	6  .  Op: ==
	// 	7  .  Y: *ast.SelectorExpr {
	// 	8  .  .  X: *ast.Ident {
	// 	9  .  .  .  NamePos: -
	//    10  .  .  .  Name: "u"
	//    11  .  .  }
	//    12  .  .  Sel: *ast.Ident {
	//    13  .  .  .  NamePos: -
	//    14  .  .  .  Name: "Name"
	//    15  .  .  }
	//    16  .  }
	//    17  }
	// err = ast.Print(token.NewFileSet(), expr)
	// if err != nil {
	// 	t.Error(err)
	// }

	be := expr.(*ast.BinaryExpr)
	bes := be.Y.(*ast.SelectorExpr)
	if bes.Sel.Name != "Name" { // got the Name
		t.Errorf("bad case: %s != %s", bes.Sel.Name, "Name")
	}

	// Replace the ast.Node
	node := astutil.Apply(expr, func(c *astutil.Cursor) bool {
		n := Must1(parser.ParseExpr(`age == u.Age`))
		c.Replace(n)
		return true
	}, func(c *astutil.Cursor) bool {
		return false
	})

	// 	0  *ast.BinaryExpr {
	// 	1  .  X: *ast.Ident {
	// 	2  .  .  NamePos: -
	// 	3  .  .  Name: "age"
	// 	4  .  }
	// 	5  .  OpPos: -
	// 	6  .  Op: ==
	// 	7  .  Y: *ast.SelectorExpr {
	// 	8  .  .  X: *ast.Ident {
	// 	9  .  .  .  NamePos: -
	//    10  .  .  .  Name: "u"
	//    11  .  .  }
	//    12  .  .  Sel: *ast.Ident {
	//    13  .  .  .  NamePos: -
	//    14  .  .  .  Name: "Age"
	//    15  .  .  }
	//    16  .  }
	//    17  }
	// err = ast.Print(token.NewFileSet(), node)
	// if err != nil {
	// 	t.Error(err)
	// }

	{
		be := node.(*ast.BinaryExpr)
		bes := be.Y.(*ast.SelectorExpr)
		if bes.Sel.Name != "Age" { // got the Name
			t.Errorf("bad case: %s != %s", bes.Sel.Name, "Age")
		}
	}
}
