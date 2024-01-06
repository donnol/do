package main

const (
	struct2structTmpl = `
	func To{{.ToType}}(in *{{.FromType}}) *{{.ToType}} {
		return &{{.ToType}}{
			{{range $k,$v:=.FieldPair}}
			{{$v.To}}: {{if $v.From}} in.{{$v.From}} {{else}} "" {{end}},
			{{end}}
		}
	}
	`
)
