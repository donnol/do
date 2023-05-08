package do

import (
	"reflect"
	"sort"
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

func TestKeys(t *testing.T) {
	t.Parallel()

	result1 := Keys(map[int]string{
		1: "j",
		2: "k",
	})
	sort.Slice(result1, func(i, j int) bool {
		return result1[i] < result1[j]
	})

	want := []int{1, 2}
	if !reflect.DeepEqual(result1, want) {
		t.Errorf("bad case, %+v != %+v", result1, want)
	}
}

func TestValuess(t *testing.T) {
	t.Parallel()

	result1 := Values(map[int]string{
		1: "j",
		2: "k",
	})
	sort.Slice(result1, func(i, j int) bool {
		return result1[i] < result1[j]
	})

	want := []string{"j", "k"}
	if !reflect.DeepEqual(result1, want) {
		t.Errorf("bad case, %+v != %+v", result1, want)
	}
}

func BenchmarkKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Keys(map[int]string{
			1: "j",
			2: "k",
		})
	}
}

func BenchmarkValues(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Values(map[int]string{
			1: "j",
			2: "k",
		})
	}
}

func TestMergeKeyValue(t *testing.T) {
	type args[K comparable, V any] struct {
		m1 map[K]V
		m2 map[K]V
	}
	tests := []struct {
		name string
		args args[string, int]
		want map[string]int
	}{
		// TODO: Add test cases.
		{
			name: "string-int",
			args: args[string, int]{
				m1: map[string]int{
					"a1": 1,
					"a2": 2,
					"a3": 3,
				},
				m2: map[string]int{
					"a1": 11,
					"a4": 4,
					"a5": 5,
					"a6": 6,
				},
			},
			want: map[string]int{
				"a1": 11,
				"a2": 2,
				"a3": 3,
				"a4": 4,
				"a5": 5,
				"a6": 6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeKeyValue(tt.args.m1, tt.args.m2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeKeyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
