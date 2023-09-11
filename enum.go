package do

type Enum[T ~int | ~string] struct {
	Name  string `json:"name"`
	Value T      `json:"value"`
}
