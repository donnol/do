package do

import "testing"

type tableWriterEntity struct {
	Name string
	Age  int
}

type tableWriter struct {
	headers []string
	rows    [][]string
}

func (t *tableWriter) SetHeader(headers []string) {
	t.headers = headers
}

func (t *tableWriter) SetRows(rows [][]string) {
	t.rows = rows
}

func TestTableWrite(t *testing.T) {
	type args struct {
		w    TableWriter
		list []tableWriterEntity
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				w: &tableWriter{},
				list: []tableWriterEntity{
					{
						Name: "jd",
						Age:  18,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TableWrite(tt.args.w, tt.args.list)

			AssertSlice(t, tt.args.w.(*tableWriter).headers, []string{"Name", "Age"})
			AssertSlice(t, tt.args.w.(*tableWriter).rows[0], []string{"jd", "18"})
		})
	}
}
