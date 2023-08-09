package parser

import "testing"

func TestReplaceNameof(t *testing.T) {
	type args struct {
		pkg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "example",
			args: args{
				pkg: "../example/nameof",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReplaceNameof(tt.args.pkg)
		})
	}
}
