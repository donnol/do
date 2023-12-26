package do

import (
	"testing"
	"unsafe"
)

func TestCopy(t *testing.T) {
	// value is same, so the pointer
	{
		var a int = 1
		ap := &a
		r := ap
		Assert(t, *r, a)
		Assert(t, unsafe.Pointer(r) == unsafe.Pointer(ap), true)

		// change will effect origin variable
		*r = 2
		Assert(t, 2, a)
		Assert(t, *r, a)
		Assert(t, unsafe.Pointer(r) == unsafe.Pointer(ap), true)
	}

	// value is same, but pointer isn't
	{
		var a int = 1
		ap := &a
		r := Copy(ap)
		Assert(t, *r, a)
		Assert(t, unsafe.Pointer(r) != unsafe.Pointer(ap), true)

		// change will not effect origin variable
		*r = 2
		Assert(t, a != 2, true)
		Assert(t, *r != a, true)
		Assert(t, unsafe.Pointer(r) != unsafe.Pointer(ap), true)
	}

	type m struct{ name string }

	// value is same, so the pointer
	{
		var a = m{name: "jd"}
		ap := &a
		r := ap
		Assert(t, *r, a)
		Assert(t, unsafe.Pointer(r) == unsafe.Pointer(ap), true)

		// change will effect origin variable
		r.name = "jj"
		Assert(t, "jj", a.name)
		Assert(t, *r, a)
		Assert(t, unsafe.Pointer(r) == unsafe.Pointer(ap), true)
	}

	// value is same, but pointer isn't
	{
		var a = m{name: "jd"}
		ap := &a
		r := Copy(ap)
		Assert(t, *r, a)
		Assert(t, unsafe.Pointer(r) != unsafe.Pointer(ap), true)

		// change will not effect origin variable
		r.name = "jj"
		Assert(t, a.name != "jj", true)
		Assert(t, *r != a, true)
		Assert(t, unsafe.Pointer(r) != unsafe.Pointer(ap), true)
	}
}
