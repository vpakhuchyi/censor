// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"unicode/utf8"
	"unsafe"
)

/*
	NOTE: This is a copy of the strings.Builder implementation from the Go standard library.
	It has copyCheck method removed to reduce the allocations.
*/

// A Builder is used to efficiently build a string using [Builder.Write] methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Builder.
type Builder struct {
	buf []byte
}

func New() *Builder {
	return &Builder{buf: make([]byte, 0, 64)}
}

// String returns the accumulated string.
func (b *Builder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *Builder) Len() int { return len(b.buf) }

// Cap returns the capacity of the builder's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *Builder) Cap() int { return cap(b.buf) }

// Reset resets the [Builder] to be empty.
func (b *Builder) Reset() {
	b.buf = b.buf[:0]
}

// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *Builder) grow(n int) {
	buf := make([]byte, 2*cap(b.buf)+n)
	copy(buf, b.buf)
	b.buf = buf
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *Builder) Grow(n int) {
	if n < 0 {
		panic("strings.Builder.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		b.grow(n)
	}
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p).
func (b *Builder) Write(p []byte) int {
	b.buf = append(b.buf, p...)
	return len(p)
}

// WriteByte appends the byte c to b's buffer.
func (b *Builder) WriteByte(c byte) {
	b.buf = append(b.buf, c)
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r.
func (b *Builder) WriteRune(r rune) int {
	n := len(b.buf)
	b.buf = utf8.AppendRune(b.buf, r)
	return len(b.buf) - n
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s.
func (b *Builder) WriteString(s string) int {
	b.buf = append(b.buf, s...)
	return len(s)
}
