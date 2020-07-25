package helper

import (
	"bytes"
	"encoding/binary"
	"log"
)

// ContainsAll check if a []byte has byte in it
func ContainsAll(s []byte, e byte) bool {
	for _, a := range s {
		if a != e {
			return false
		}
	}
	return true
}

// FailOnError is a general helper to deal with error logging
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func insertNth(s string, n int) string {
	var buffer bytes.Buffer
	var n1 = n - 1
	var l1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%n == n1 && i != l1 {
			buffer.WriteRune(' ')
		}
	}
	return buffer.String()
}

// CRC16 calculates the CRC bytes from hexbytes
func CRC16(hexBytes []byte) []byte {

	var CRC uint32 = 0
	var bitNumber uint32 = 0
	var carry uint32 = 0

	for _, hexByte := range hexBytes {
		CRC = CRC ^ uint32(hexByte)
		bitNumber = 0
		for bitNumber != 8 {
			carry = uint32(CRC) & 1
			CRC = CRC >> 1
			if carry == 1 {
				CRC = CRC ^ uint32(40961)
			}
			bitNumber = bitNumber + 1
		}
	}
	res := make([]byte, 4)
	binary.BigEndian.PutUint32(res, CRC)
	return res
}
