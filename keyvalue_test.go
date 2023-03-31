package do

import (
	"reflect"
	"testing"
)

func TestKeyValueBy(t *testing.T) {
	t.Parallel()

	result1 := KeyValueBy([]string{"a", "aa", "aaa"}, func(str string) (string, int) {
		return str, len(str)
	})

	want := map[string]int{"a": 1, "aa": 2, "aaa": 3}
	if !reflect.DeepEqual(result1, want) {
		t.Errorf("bad case, %+v != %+v", result1, want)
	}
}

func TestKeyBy(t *testing.T) {
	t.Parallel()

	result1 := KeyBy([]string{"a", "aa", "aaa"}, func(str string) string {
		return str
	})

	want := map[string]string{"a": "a", "aa": "aa", "aaa": "aaa"}
	if !reflect.DeepEqual(result1, want) {
		t.Errorf("bad case, %+v != %+v", result1, want)
	}
}

func BenchmarkKeyValueBy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		KeyValueBy([]string{"a", "aa", "aaa"}, func(str string) (string, int) {
			return str, len(str)
		})
	}
}
