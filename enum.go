package do

type Enum[T any] struct {
	Name  string `json:"name"`
	Value T      `json:"value"`
}
