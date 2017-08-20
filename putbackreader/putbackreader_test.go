package putbackreader

import (
	"bytes"
	"io"
	"testing"
)

var payload = []byte(`Lorem ipsum dolor sit amet,
consectetur adipiscing elit, sed do eiusmod
tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation
ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit
esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
occaecat cupidatat non proident, sunt in culpa qui officia
deserunt mollit anim id est laborum.`)

func Test01StraightRead(t *testing.T) {
	pbr := NewPutBackReader(bytes.NewReader(payload))
	buf := new(bytes.Buffer)
	buf.ReadFrom(pbr)
	if bytes.Compare(payload, buf.Bytes()) != 0 {
		t.Error("Failed on straight read test")
		t.Fail()
	}
}

func Test02PutBackRead(t *testing.T) {
	b := make([]byte, 20, 20)
	pbr := NewPutBackReader(bytes.NewReader(payload))
	pbr.Read(b)
	// put it back
	pbr.SetBackBytes(b)
	buf := new(bytes.Buffer)
	// now read all again
	buf.ReadFrom(pbr)
	if bytes.Compare(payload, buf.Bytes()) != 0 {
		t.Error("Failed on put back read test")
		t.Fail()
	}
}

func Test02PutBackBiggerRead(t *testing.T) {
	pbr := NewPutBackReader(bytes.NewReader(payload))
	// drain it.
	buf := new(bytes.Buffer)
	buf.ReadFrom(pbr)
	// put it back.
	pbr.SetBackBytes(buf.Bytes())
	// small buffer reads
	b := make([]byte, 10, 10)
	pc := make([]byte, 0, len(payload))
	for {
		n, err := pbr.Read(b)
		if err != nil && err != io.EOF {
			t.Errorf("read failed unexpectedly: %s", err)
			t.Fail()
		}
		if n == 0 {
			break
		}
		pc = append(pc, b[:n]...)
	}
	if bytes.Compare(payload, pc) != 0 {
		t.Errorf("Failed on put back read test, %s", string(pc))
		t.Fail()
	}

}
