package do

import "fmt"

type (
	SQLQuery string
	SQLArgs  []any

	SQLQueryFunc func() (query SQLQuery, args SQLArgs)
)

// Append add sql to query, if the q is a format string of fmt, fmt.Sprintf will be used to concat the sql; the q will be the format, the left sql will be arguments.
func (q SQLQuery) Append(ss ...SQLQuery) (r SQLQuery) {
	if isFormatString(q.Raw()) {
		args := make([]any, 0, len(ss))
		for _, s := range ss {
			args = append(args, s)
		}
		return SQLQuery(fmt.Sprintf(q.Raw(), args...))
	}
	r += q
	for _, s := range ss {
		r += s
	}
	return
}

func (q SQLQuery) Raw() string {
	return string(q)
}

func isFormatString(format string) bool {
	ok, _ := checkPrintf(format)
	return ok
}

func (a SQLArgs) Append(s SQLArgs) SQLArgs {
	return append(a, s...)
}
func (a SQLArgs) Raw() []any {
	return a
}

// SQLProcess process sql and args one by one, add more parts based on first sql
func SQLProcess(funcs ...SQLQueryFunc) (query SQLQuery, args SQLArgs) {
	if len(funcs) == 0 {
		return
	}
	query, args = funcs[0]()

	qs := make([]SQLQuery, 0, len(funcs))
	as := make(SQLArgs, 0, len(funcs))
	for i, f := range funcs {
		if i == 0 {
			continue
		}

		q, a := f()
		qs = append(qs, q)
		as = append(as, a...)
	}
	query = query.Append(qs...)
	args = args.Append(as)
	return
}
