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

type fromUser struct {
	Name string
}

func (u fromUser) To() toUser {
	return toUser{
		Name: u.Name + "1",
	}
}

type fromUser2 struct {
	Name string
}

func (u fromUser2) To() *toUser {
	return &toUser{
		Name: u.Name + "1",
	}
}

type fromUser3 struct {
	Name string
}

func (u fromUser3) To(v fromUser3) toUser {
	return toUser{
		Name: v.Name + "1",
	}
}

type toUser struct {
	Name string
}

func NewToUser() *toUser {
	return &toUser{}
}

func (u *toUser) From(v fromUser) {
	u.Name = v.Name + "1"
}

type toUser2 struct {
	Name string
}

func (u toUser2) From(v fromUser) toUser2 {
	u.Name = v.Name + "1"
	return u
}

func TestMapFrom(t *testing.T) {
	got := MapFrom([]fromUser{{Name: "jd"}, {Name: "jc"}}, NewToUser)
	AssertSlicePtr(t, got, []*toUser{{"jd1"}, {"jc1"}})

	AssertSlice(t, MapSlicePtr(got), []toUser{{"jd1"}, {"jc1"}})
}

func TestMapFrom2(t *testing.T) {
	got := MapFrom2[fromUser, toUser2]([]fromUser{{Name: "jd"}, {Name: "jc"}})
	AssertSlice(t, got, []toUser2{{"jd1"}, {"jc1"}})
}

func TestMapTo(t *testing.T) {
	got := MapTo[fromUser, toUser]([]fromUser{{Name: "jd"}, {Name: "jc"}})
	AssertSlice(t, got, []toUser{{"jd1"}, {"jc1"}})
	{
		got := MapTo[fromUser2, *toUser]([]fromUser2{{Name: "jd"}, {Name: "jc"}})
		AssertSlicePtr(t, got, []*toUser{{"jd1"}, {"jc1"}})
	}
}

func TestMapTo2(t *testing.T) {
	got := MapTo2[fromUser3, toUser]([]fromUser3{{Name: "jd"}, {Name: "jc"}})
	AssertSlice(t, got, []toUser{{"jd1"}, {"jc1"}})
}

func BenchmarkMapFrom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapFrom([]fromUser{{Name: "jd"}, {Name: "jc"}}, NewToUser)
	}
}

func BenchmarkMapFrom2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapFrom2[fromUser, toUser2]([]fromUser{{Name: "jd"}, {Name: "jc"}})
	}
}

func BenchmarkMapTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapTo[fromUser, toUser]([]fromUser{{Name: "jd"}, {Name: "jc"}})
	}
}

func BenchmarkMapTo2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapTo[fromUser2, *toUser]([]fromUser2{{Name: "jd"}, {Name: "jc"}})
	}
}

func BenchmarkMapTo3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapTo2[fromUser3, toUser]([]fromUser3{{Name: "jd"}, {Name: "jc"}})
	}
}
