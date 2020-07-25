package codec12

import (
	"encoding/binary"
	"errors"
)

// Decode requires a hexBytes (Array of Hex bytes) and Parses to return AVLType
func Decode(hexBytes []byte) (string, error) {

	var pResponseQuantity1 = 1
	var pResponseQuantity2 = 1
	var pType = 1
	var pResponseSize = 4

	p := 0
	responseQuantity1 := int(hexBytes[p])
	p += pResponseQuantity1

	responseType := hexBytes[p]
	p += pType

	if responseType != 0x06 {
		return "", errors.New("Cannot decode a non-response string")
	}

	responseSize := int(binary.BigEndian.Uint32(hexBytes[p : p+pResponseSize]))
	p += pResponseSize

	response := hexBytes[p : p+responseSize]
	p += responseSize

	responseQuantity2 := int(hexBytes[p])
	p += pResponseQuantity2

	if responseQuantity1 != responseQuantity2 {
		return "", errors.New("Something went wrong")
	}

	return string(response), nil
}
