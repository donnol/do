package do

func FuzzWrap(v string) string {
	return "%" + v + "%"
}
