package client

import (
	"encoding/binary"
	"log"
	"tunProxy/crypt"
)

func pdt(b []byte) {
	for i := 0; i < len(b); i++ {
		log.Printf("%02x ", b[i])
	}
	log.Panicln(" ")
}
func magic() []byte {
	return []byte{0x89, 0x64}
}
func EncodeMessage(in, key, iv []byte) []byte {
	encrypt := crypt.AesEncrypt(in, key, iv)
	bytes := magic()
	lenb := make([]byte, 4)
	binary.BigEndian.PutUint32(lenb, uint32(len(encrypt)))
	bytes = append(bytes, lenb...)
	return append(bytes, encrypt...)
}
func DecodeMessage(in, key, iv []byte) []byte {
	lenb := int(binary.BigEndian.Uint32(in[2:6]))
	b := in[6:]
	if lenb != len(b) {
		log.Fatalf("warning : actual length : %d of data received does not equal to the length :%d   in message head\n", len(b), lenb)
		pdt(in)
	}
	return crypt.AesDecrypt(b, key, iv)
}
