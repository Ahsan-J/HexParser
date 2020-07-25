package codec14

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
)

// Decode requires a hexBytes (Array of Hex bytes) and Parses to return AVLType
func Decode(hexBytes []byte) (string, string, error) {

	var pResponseQuantity1 = 1
	var pResponseQuantity2 = 1
	var pType = 1
	var pIMEI = 8
	var pResponseSize = 4

	p := 0
	responseQuantity1 := int(hexBytes[p])
	p += pResponseQuantity1

	responseType := hexBytes[p]
	p += pType

	if responseType != 0x06 {
		return "", "", errors.New("Cannot decode a non-response string")
	}

	responseSize := int(binary.BigEndian.Uint32(hexBytes[p : p+pResponseSize]))
	p += pResponseSize

	imei := hex.EncodeToString(hexBytes[p : p+pIMEI])
	p += pIMEI

	response := hexBytes[p : p+responseSize-pIMEI]
	p += responseSize - pIMEI

	responseQuantity2 := int(hexBytes[p])
	p += pResponseQuantity2

	if responseQuantity1 != responseQuantity2 {
		return "", "", errors.New("Something went wrong")
	}

	return string(response), imei, nil
}
