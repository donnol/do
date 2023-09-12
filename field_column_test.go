package do

import (
	"reflect"
	"testing"
)

type UserEmbed struct {
	User
	Phone string
}

func Test_fieldsByColumnName(t *testing.T) {
	type args struct {
		t           any
		validName   map[string]struct{}
		fieldMapper func(string) string
	}

	u := &User{}
	id := &u.Id
	name := &u.Name

	ue := &UserEmbed{}
	ueid := &ue.Id
	uename := &ue.Name
	uephone := &ue.Phone

	tests := []struct {
		name           string
		args           args
		wantNameValues map[string]any
	}{
		// TODO: Add test cases.
		{
			name: "user",
			args: args{
				t: u,
				validName: map[string]struct{}{
					"id":   {},
					"name": {},
				},
				fieldMapper: nil,
			},
			wantNameValues: map[string]any{
				"id":   id,
				"name": name,
			},
		},
		{
			name: "user embed",
			args: args{
				t: ue,
				validName: map[string]struct{}{
					"id":    {},
					"name":  {},
					"phone": {},
				},
				fieldMapper: nil,
			},
			wantNameValues: map[string]any{
				"id":    ueid,
				"name":  uename,
				"phone": uephone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNameValues := fieldsByColumnName(tt.args.t, tt.args.validName, tt.args.fieldMapper); !reflect.DeepEqual(gotNameValues, tt.wantNameValues) {
				t.Errorf("fieldsByColumnName() = %v, want %v", gotNameValues, tt.wantNameValues)
			}
		})
	}
}

func TestObjectFields(t *testing.T) {
	genObj := ObjectAndFieldsHelper[UserEmbed]()
	genObj(nil)
	genObj(nil)
}

type ValuePtrs struct {
	id   int
	name string
	addr string
}

func (v *ValuePtrs) ValuePtrs() []any {
	return []any{
		&v.id,
		&v.name,
		&v.addr,
	}
}

func TestObjectFieldsValuePtrs(t *testing.T) {
	genobj := ObjectAndFieldsHelper[ValuePtrs]()
	_, fields := genobj(nil)
	Assert(t, len(fields), 3)
	Assert(t, fields[0] != nil, true)
	Assert(t, fields[1] != nil, true)
	Assert(t, fields[2] != nil, true)

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("no panic occur")
			}
		}()
		// will panic because *ValuePtrs is not a interface{ ValuePtrs() []any } and colTypes is nil
		ObjectAndFieldsHelper[*ValuePtrs]()(nil)
	}()
}
