// Package putbackreader
// This implements io.Reader, giving you the option to "put back"
// bytes that were already read.  This is mainly used with a
// json.Decoder and you're working with json-like serial input where
// a separator character is used, which pisses the json decoder off
// right good.
//
// Note: This package is not thread safe.
package putbackreader

import (
	"io"
)

// PutBackReader - This implements io.Reader
type PutBackReader struct {
	putBack []byte
	r       io.Reader
}

// NewPutBackReader - allocate a PutBackReader
func NewPutBackReader(r io.Reader) *PutBackReader {
	return &PutBackReader{r: r}
}

// Read - implements io.Reader.
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
func (pbr *PutBackReader) SetBackBytes(b []byte) {
	pbr.putBack = b
}

// BackBytes - get the internal byte slice of back bytes.
func (pbr *PutBackReader) BackBytes() []byte {
	return pbr.putBack
}
