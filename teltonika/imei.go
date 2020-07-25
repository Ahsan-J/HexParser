package teltonika

import (
	"encoding/binary"
	"errors"
	"regexp"
)

// parseIMEI parses the hex bytes to output the IMEI number of the devices
func parseIMEI(hexBytes []byte) (string, error) {
	var pIMEILength = 2

	lengthIMEI := int(binary.BigEndian.Uint16(hexBytes[:len(hexBytes)-pIMEILength]))

	if len(hexBytes) != pIMEILength+lengthIMEI {
		return "", errors.New("Invalid")
	}

	deviceIMEI := string(hexBytes[pIMEILength:])

	if res, _ := regexp.MatchString("^[0-9]+$", deviceIMEI); !res { // case where IMEI is not numeric
		return "", errors.New("Not a Valid IMEI")
	}

	return deviceIMEI, nil
}
