package blockchain

import (
	"fmt"
	// "github.com/OhYee/goutils/bytes"
	"strings"
)

// HashCode hash value
type HashCode []byte

// NewHashCodeFromBytes return a HashCode
func NewHashCodeFromBytes(b []byte) HashCode {
	return b

}

// NewHashCodeFromBytes return a HashCode
func NewHashCodeFromString(s string) HashCode {
	b := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		b[i*2] = hex2ord(s[i*2])*16 + hex2ord(s[i*2+1])
	}
	return b
}

func hex2ord(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

// ToBytes return the []byte of HashCode
func (h HashCode) ToBytes() []byte {
	return h
	// return bytes.FromString(string(h))
}

func (h HashCode) String() string {
	ss := make([]string, len(h))
	for idx, bb := range h {
		ss[idx] = fmt.Sprintf("%02x", bb)
	}
	return strings.Join(ss, "")
	// return string(h)
}
