package teltonika

import (
	"fmt"
	"strings"

	"github.com/Ahsan-J/HexParser/codec12"
	"github.com/Ahsan-J/HexParser/codec14"
	"github.com/Ahsan-J/HexParser/model"
)

// GenerateHex creates a hex representation of command
func GenerateHex(command string, codecID string, imei int) model.EncodedCodec {
	var encodedCodec model.EncodedCodec
	var err error

	switch strings.ToUpper(codecID) {
	case "0C":
		encodedCodec, err = codec12.Encode(command)
		if err != nil {
			fmt.Println(err.Error())
		}

		break
	case "0E":
		encodedCodec, err = codec14.Encode(command, imei)
		if err != nil {
			fmt.Println(err.Error())
		}

		break
	default:
		encodedCodec.Meta.IsValid = false
		encodedCodec.Meta.Message = "Hex encoding codec not supported"
		break
	}

	return encodedCodec
}
