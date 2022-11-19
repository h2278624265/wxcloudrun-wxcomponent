package utils

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(n int) []byte {
	data := int32(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func BytesToInt(bys []byte) int {
	bytebuf := bytes.NewBuffer(bys)
	var data int32
	binary.Read(bytebuf, binary.BigEndian, &data)
	return int(data)
}
