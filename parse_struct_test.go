package do

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func TestResolve(t *testing.T) {
	sm, fm, err := resolve("Error")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = sm, fm
	// jsonPrint(os.Stdout, sm)
	// jsonPrint(os.Stdout, fm)
}

func TestCollectStructComment(t *testing.T) {
	for _, cas := range []any{
		&Error[error]{},
	} {
		s, err := ResolveStruct(cas)
		if err != nil {
			t.Fatal(err)
		}
		fields := s.GetFields()
		_ = fields
		// jsonPrint(os.Stdout, fields)
	}
}

func TestResolveStruct(t *testing.T) {
	s, err := ResolveStruct(&Error[error]{})
	if err != nil {
		t.Fatal(err)
	}
	_ = s
	// jsonPrint(os.Stdout, s)
}

var _ = jsonPrint

func jsonPrint(w io.Writer, in any) {
	var data []byte
	if v, ok := in.([]byte); ok {
		data = v
	} else {
		var err error
		data, err = json.Marshal(in)
		if err != nil {
			panic(err)
		}
	}
	var buf = new(bytes.Buffer)
	if err := json.Indent(buf, data, "", "\t"); err != nil {
		panic(err)
	}
	if _, err := buf.WriteTo(w); err != nil {
		panic(err)
	}
}
