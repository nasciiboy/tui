package tui

// Simple byte buffer for marshaling data.

import (
	"errors"
)

// A Buffer is a variable-sized buffer of uint64 with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer struct {
	buf       []uint64     // contents are the bytes buf[off : len(buf)]
	off       int          // read at &buf[off], write at &buf[len(buf)]
	bootstrap [64]uint64   // memory to hold first slice; helps small buffers avoid allocation.
	lastRead  readOp       // last read operation, so that Unread* can work correctly.
}

// The readOp constants describe the last action performed on
// the buffer, so that UnreadRune and UnreadByte can check for
// invalid usage. opReadRuneX constants are choosen such that
// converted to int they correspond to the rune size that was read.
type readOp int

const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid          = 0  // Non-read operation.
	opReadRune1        = 1  // Read rune of size 1.
	opReadRune2        = 2  // Read rune of size 2.
	opReadRune3        = 3  // Read rune of size 3.
	opReadRune4        = 4  // Read rune of size 4.
)

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("bytes.Buffer: too large")

// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like Read, Write, Reset, or Truncate).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (b *Buffer) Data() []uint64 { return b.buf[b.off:] }

func (b *Buffer) CellData() []Cell {
  c := make( []Cell, b.Len() )

  for i, d := range b.Data() {
    attrs, color, _, r := extractData( d )
    c[i] = Cell{ Attrs: attrs, Color: color, Ch: r }
  }

  return c
}

// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
func (b *Buffer) String() (str string) {
	if b == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}

  for _, cell := range b.buf[b.off:] {
    if (cell & hasKey) == 0 {
      r   := rune(cell & runeMask)
      str += string(r)
      continue
    }

    str += " "
  }

	return
}

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (b *Buffer) Len() int { return len(b.buf) - b.off }

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
func (b *Buffer) Cap() int { return cap(b.buf) }

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (b *Buffer) Truncate(n int) {
	b.lastRead = opInvalid
	switch {
	case n < 0 || n > b.Len():
		panic("bytes.Buffer: truncation out of range")
	case n == 0:
		// Reuse buffer space.
		b.off = 0
	}
	b.buf = b.buf[0 : b.off+n]
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (b *Buffer) Reset() { b.Truncate(0) }

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) grow(n int) int {
	m := b.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.off != 0 {
		b.Truncate(0)
	}
	if len(b.buf)+n > cap(b.buf) {
		var buf []uint64
		if b.buf == nil && n <= len(b.bootstrap) {
			buf = b.bootstrap[0:]
		} else if m+n <= cap(b.buf)/2 {
			// We can slide things down instead of allocating a new
			// slice. We only need m+n <= cap(b.buf) to slide, but
			// we instead let capacity get twice as large so we
			// don't spend all our time copying.
			copy(b.buf[:], b.buf[b.off:])
			buf = b.buf[:m]
		} else {
			// not enough space anywhere
			buf = makeSlice(2*cap(b.buf) + n)
			copy(buf, b.buf[b.off:])
		}
		b.buf = buf
		b.off = 0
	}
	b.buf = b.buf[0 : b.off+m+n]
	return b.off + m
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) Grow(n int) {
	if n < 0 {
		panic("bytes.Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[0:m]
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (b *Buffer) Write(p []uint64) (n int, err error) {
	b.lastRead = opInvalid
	m := b.grow(len(p))
	return copy(b.buf[m:], p), nil
}

func (b *Buffer) ReadFrom(r Buffer) (n int, err error) {
	return b.Write( r.Data() )
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
func (b *Buffer) WriteString(s string) (n int, err error) {
	b.lastRead = opInvalid
	m := b.grow(len(s))
  p := make( []uint64, 0, 64 )
  for _, r := range( s ) {
    p = append( p, uint64(r))
  }

	return copy(b.buf[m:], p), nil
}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []uint64 {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]uint64, n)
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
func (b *Buffer) WriteU64( u uint64 ) error {
	b.lastRead = opInvalid
	m := b.grow(1)
	b.buf[m] = u
	return nil
}


func (b *Buffer) WriteCell( cell Cell ) error {
  var c uint64
  c = uint64(cell.Ch) | uint64(cell.Color) << 48 | uint64(cell.Attrs) << 56

	b.lastRead = opInvalid
	m := b.grow(1)
	b.buf[m] = c
	return nil
}

func (b *Buffer) WriteCells( cells []Cell ) error {
  for _, c := range cells {
    err := b.WriteCell( c )

    if err != nil { return err }
  }

	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match bufio.Writer's WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with ErrTooLarge.
func (b *Buffer) WriteRune(r rune) error {
  return b.WriteU64( uint64(r) )
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (b *Buffer) Next(n int) []uint64 {
	b.lastRead = opInvalid
	m := b.Len()
	if n > m {
		n = m
	}
	data := b.buf[b.off : b.off+n]
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return data
}
