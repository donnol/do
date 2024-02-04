package do

import (
	"testing"
)

func TestEscapeStruct(t *testing.T) {
	type Embed struct {
		Addr string
	}
	type E struct {
		Embed
		Result string
		Range  string
		Age    int
	}
	type M struct {
		Embed
		Name string
		Es   []E
		Ea   [1]E
	}
	type args struct {
		v       any
		escaper func(field any) (any, error)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				v:       &M{Name: "< 0.1", Embed: Embed{Addr: "< 0.2"}, Es: []E{{Embed: Embed{Addr: "<0.3"}, Result: "< 0.1", Range: "0.1~0.2", Age: 10}}, Ea: [1]E{{Result: "<0.6"}}},
				escaper: XMLEscaper,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EscapeStruct(tt.args.v, tt.args.escaper); (err != nil) != tt.wantErr {
				t.Errorf("EscapeStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			Assert(t, tt.args.v.(*M).Embed.Addr, "&lt; 0.2")
			Assert(t, tt.args.v.(*M).Name, "&lt; 0.1")
			Assert(t, tt.args.v.(*M).Es[0].Result, "&lt; 0.1")
			Assert(t, tt.args.v.(*M).Ea[0].Result, "&lt;0.6")
		})
	}
}
