package ros

import (
	"io"
)

// Reader implements io.Reader interface and provides a no-copy way to read N
// bytes via Next() method. Reader is used by generated message to de-serialize
// byte arrays ([]uint8) without/ copying underlying data.
type Reader struct {
	s []byte
	i int
}

// NewReader creates new Reader and adopts the byte slice.
// The caller must not modify the slice after this call.
func NewReader(s []byte) *Reader {
	return &Reader{s, 0}
}

// Read implements the io.Reader interface. Like the other reader
// implementations, this implementation copies data from the original slice
// into "b".
func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}
	n = copy(b, r.s[r.i:])
	r.i += n
	return
}

// Next returns a slice containing the next n bytes from the buffer, advancing
// the buffer as if the bytes had been returned by Read. The resulting slice is
// a sub-slice of the original slice.
//
// Asking for more bytes than available would returns only the remaining bytes.
// Calling Next on an empty buffer, or after the buffer has been exhausted,
// returns an empty slice.
func (r *Reader) Next(n int) []byte {
	m := len(r.s) - r.i
	if n > m {
		n = m
	}
	data := r.s[r.i : r.i+n]
	r.i += n
	return data
}
