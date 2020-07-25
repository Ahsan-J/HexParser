package codec13

import (
	"encoding/binary"
	"errors"
	"time"
)

// Decode requires a hexBytes (Array of Hex bytes) and Parses to return AVLType
func Decode(hexBytes []byte) (string, time.Time, error) {

	var pResponseQuantity1 = 1
	var pResponseQuantity2 = 1
	var pType = 1
	var pResponseSize = 4
	var pTimestamp = 4

	p := 0
	responseQuantity1 := int(hexBytes[p])
	p += pResponseQuantity1

	responseType := hexBytes[p]
	p += pType

	if responseType != 0x05 {
		return "", time.Now(), errors.New("Cannot decode a non-response string")
	}

	responseSize := int(binary.BigEndian.Uint32(hexBytes[p : p+pResponseSize]))
	p += pResponseSize

	packetTimeRaw := int32(binary.BigEndian.Uint32(hexBytes[p : p+pTimestamp]))
	packetTime := time.Unix(int64(packetTimeRaw)/1000, 0)
	p += pTimestamp

	response := hexBytes[p : p+responseSize-pTimestamp]
	p += responseSize - pTimestamp

	responseQuantity2 := int(hexBytes[p])
	p += pResponseQuantity2

	if responseQuantity1 != responseQuantity2 {
		return "", time.Now(), errors.New("Something went wrong")
	}

	return string(response), packetTime, nil
}
