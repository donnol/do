package do

import (
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	m := NewMap[string, int](10)

	m.Insert("k1", 1)
	r := m.Get("k1")
	if r != 1 {
		t.Errorf("bad case: %v != %v", r, 1)
	}

	m.Insert("k2", 2)
	r = m.Get("k2")
	if r != 2 {
		t.Errorf("bad case: %v != %v", r, 2)
	}

	m.Range(func(key string, value int) {
		if key == "k1" && value != 1 {
			t.Errorf("bad case: %v != %v", r, 1)
		}
		if key == "k2" && value != 2 {
			t.Errorf("bad case: %v != %v", r, 2)
		}
	})

	m.Remove("k1")
	v, ok := m.Lookup("k1")
	if ok {
		t.Errorf("bad case: exist k1 value: %v", v)
	}
}

func TestMapAny(t *testing.T) {
	m := NewMap[string, any](10)

	m.Insert("k1", 1)
	r := m.Get("k1").(int)
	if r != 1 {
		t.Errorf("bad case: %v != %v", r, 1)
	}
	{
		r := ValueAs[string, int](m, "k1")
		if r != 1 {
			t.Errorf("bad case: %v != %v", r, 1)
		}
	}

	m.Insert("k2", "2")
	r1 := m.Get("k2").(string)
	if r1 != "2" {
		t.Errorf("bad case: %v != %v", r, 2)
	}
	{
		r1 := ValueAs[string, string](m, "k2")
		if r1 != "2" {
			t.Errorf("bad case: %v != %v", r, 2)
		}
	}

	m.Range(func(key string, value any) {
		if key == "k1" && value != 1 {
			t.Errorf("bad case: %v != %v", r, 1)
		}
		if key == "k2" && value != "2" {
			t.Errorf("bad case: %v != %v", r, 2)
		}
	})

	m.Remove("k1")
	v, ok := m.Lookup("k1")
	if ok {
		t.Errorf("bad case: exist k1 value: %v", v)
	}
}

var _ = func() bool {
	watchInterval = 2 * time.Second
	return true
}()

func TestMapWatch(t *testing.T) {

	m := NewMap[string, int](10)

	m.Insert("k1", 1, &KeyOption{Timeout: 1 * time.Second})
	r := m.Get("k1") // get right now successfully
	if r != 1 {
		t.Errorf("bad case: %v != %v", r, 1)
	}

	time.Sleep(1 * time.Second)
	_, ok := m.Lookup("k1") // get after timeout failed
	if ok {
		t.Errorf("bad case, k1 is expired")
	}

	m.Insert("k2", 2, &KeyOption{Timeout: 1 * time.Second})

	time.Sleep(2 * time.Second) // watch interval exceed, delete expired keys

	r = m.Get("k2")
	if r != 0 {
		t.Errorf("bad case: %v != %v", r, 2)
	}

	m.Range(func(key string, value int) {
		if key == "k1" {
			t.Errorf("bad case: %v != %v", r, 1)
		}
		if key == "k2" {
			t.Errorf("bad case: %v != %v", r, 2)
		}
	})

	v, ok := m.Lookup("k1")
	if ok {
		t.Errorf("bad case: exist k1 value: %v", v)
	}
}
