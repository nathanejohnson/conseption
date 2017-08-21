// Package putbackreader
//
// This wraps an io.Reader, giving you the option to "put back"
// bytes that were already read once, so that they may be read again.
// I use this in conseption in conjuction with json.Decoder Buffered()
// call, which gives me the option of manipulating the stream after an
// error.  This is mainly used with with json-like serial input where a
// separator character is used, which pisses the json decoder right off.
package putbackreader

import (
	"io"
)

// PutBackReader - This implements io.Reader
type PutBackReader struct {
	putBack []byte
	r       io.Reader
}

// NewPutBackReader - allocate a PutBackReader, wrapping r.
func NewPutBackReader(r io.Reader) *PutBackReader {
	return &PutBackReader{r: r}
}

// Read - implements io.Reader.  Not thread safe.
func (pbr *PutBackReader) Read(b []byte) (int, error) {
	var copied int
	if len(pbr.putBack) > 0 {
		copied = copy(b, pbr.putBack)
		if len(pbr.putBack) > copied {
			pbr.putBack = pbr.putBack[copied:]
			return copied, nil
		} else {
			pbr.putBack = nil
		}
	}

	n, err := pbr.r.Read(b[copied:])
	return n + copied, err
}

// SetBackBytes - "put back" some bytes, presumably that were already
// read.  This gives you the option to edit the stream, so to speak.
//
// Not thread safe.
func (pbr *PutBackReader) SetBackBytes(b []byte) {
	pbr.putBack = b
}

// BackBytes - get the internal byte slice of back bytes.
func (pbr *PutBackReader) BackBytes() []byte {
	return pbr.putBack
}
