package main

import (
	"bytes"
	"fmt"
	bc "github.com/OhYee/blockchain/blockchain"
	"github.com/OhYee/cryptography_and_network_security/RSA"
	"github.com/OhYee/cryptography_and_network_security/hash/sha"
	"github.com/OhYee/goutils"
	gb "github.com/OhYee/goutils/bytes"
)

type Text struct {
	publicKey bc.HashCode
	text      string
	sign      bc.HashCode
}

func NewText(publicKey bc.HashCode, privateKey bc.HashCode, text string) *Text {
	return &Text{
		publicKey: publicKey,
		text:      text,
		sign:      bc.NewHashCodeFromBytes(rsa.Encrypto(sha.SHA256([]byte(text)), privateKey)),
	}
}

func NewTextFromBytes(b []byte) (t *Text, err error) {
	var publicKey, text, sign []byte
	buf := bytes.NewBuffer(b)

	if publicKey, err = gb.ReadWithLength32(buf); err != nil {
		return
	}
	if text, err = gb.ReadWithLength32(buf); err != nil {
		return
	}
	if sign, err = gb.ReadWithLength32(buf); err != nil {
		return
	}
	t = &Text{
		publicKey: bc.NewHashCodeFromBytes(publicKey),
		text:      string(text),
		sign:      bc.NewHashCodeFromBytes(sign),
	}
	return
}

func (text *Text) Copy() *Text {
	return &Text{
		publicKey: text.publicKey,
		text:      text.text,
		sign:      text.sign,
	}
}

func (text *Text) Varify() bool {
	return goutils.Equal(rsa.Decrypto(text.sign.ToBytes(), text.publicKey.ToBytes()), sha.SHA256([]byte(text.text)))
}

func (text *Text) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	gb.WriteWithLength32(buf, text.publicKey.ToBytes())
	gb.WriteWithLength32(buf, []byte(text.text))
	gb.WriteWithLength32(buf, text.sign.ToBytes())
	return buf.Bytes()
}

func (text *Text) String() string {
	return fmt.Sprintf("%s: %s", text.publicKey.String(), text.text)
}

type TextArray struct {
	data []*Text
}

func NewTextArray() *TextArray {
	return &TextArray{
		data: make([]*Text, 0),
	}
}

func NewTextArrayFromBytes(b []byte) (text *TextArray, err error) {
	text = NewTextArray()
	buf := bytes.NewBuffer(b)

	length32, err := gb.ReadInt32(buf)
	for i := int32(0); i < length32; i++ {
		var b []byte
		var t *Text
		if b, err = gb.ReadWithLength32(buf); err != nil {
			return
		}
		t, err = NewTextFromBytes(b)
		if err != nil {
			return
		}

		text.data = append(text.data, t)
	}
	return
}

func (text *TextArray) Copy() bc.BlockData {
	s := NewTextArray()
	s.data = make([]*Text, len(text.data))
	for i := 0; i < len(text.data); i++ {
		s.data[i] = text.data[i].Copy()
	}
	return s
}

func (text *TextArray) Reset() {
	text.data = text.data[:0]
}

func (text *TextArray) Modify(args ...interface{}) {
	for _, v := range args {
		s, ok := v.(*Text)
		if ok && s.Varify() {
			text.data = append(text.data, s)
			systemLogger.Printf("add %s successfully\n", s.String())
		} else {
			systemLogger.Printf("Ignore %v, unkown data\n", v)
		}
	}
}

func (text *TextArray) Verify() bool {
	for _, t := range text.data {
		if !t.Varify() {
			return false
		}
	}
	return true
}

func (text *TextArray) ToBytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(gb.FromInt32(int32(len(text.data))))
	for _, s := range text.data {
		gb.WriteWithLength32(buf, s.ToBytes())
	}
	return buf.Bytes()
}

func (text *TextArray) String(prefix string) string {
	buf := bytes.NewBufferString("")
	for idx, s := range text.data {
		buf.WriteString(fmt.Sprintf("%s%d %s\n", prefix, idx, s.String()))
	}
	return buf.String()
}

func (text *TextArray) FromBytes(b []byte) (data bc.BlockData, err error) {
	data, err = NewTextArrayFromBytes(b)
	return
}
