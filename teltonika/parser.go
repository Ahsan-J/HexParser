package teltonika

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"regexp"
	"strings"
	"time"

	"github.com/Ahsan-J/HexParser/codec12"
	"github.com/Ahsan-J/HexParser/codec13"
	"github.com/Ahsan-J/HexParser/codec14"
	"github.com/Ahsan-J/HexParser/codec16"
	"github.com/Ahsan-J/HexParser/codec8"
	"github.com/Ahsan-J/HexParser/codec8e"
	"github.com/Ahsan-J/HexParser/helper"
	"github.com/Ahsan-J/HexParser/model"
	"github.com/teris-io/shortid"
)

var pPreamble = 4
var pDataFieldLength = 4
var pCodecID = 1
var pCRC16 = 4

// Parser parses the teltonika hex string. refer to URL: https://wiki.teltonika-gps.com/view/Codec for more info
func Parser(hexstring string) model.DecodedCodec {
	hexstring = strings.Replace(hexstring, " ", "", -1) // Trim all whitespaces

	data := model.DecodedCodec{}
	data.ID, _ = shortid.Generate()
	data.Meta.IsValid = true
	data.Meta.HexStringLength = len(hexstring)
	data.Meta.HexString = hexstring
	data.Meta.ServerTime = time.Now().Format(time.RFC3339)

	if hexstring == "" { // case where the hex string is an empty string
		data.Meta.IsValid = false
		data.Meta.Message = "Hex received is empty"
		return data
	}

	if res, _ := regexp.MatchString("^[0-9A-Fa-f]+$", hexstring); !res { // case where hex is not in valid hexadecimal format
		data.Meta.IsValid = false
		data.Meta.Message = "Not a valid Hex string"
		return data
	}

	p := 0 // a position variable to handle the position of string

	hexBytes, _ := hex.DecodeString(hexstring)

	if !helper.ContainsAll(hexBytes[p:p+pPreamble], 0x00) { // case where the packet starts with non-zero bytes.
		var err error
		data.IMEI, err = parseIMEI(hexBytes) // try parsing for IMEI string
		if err != nil {
			data.Meta.IsValid = false
			data.Meta.Message = "Non-Zero bytes occured. Seems like an error in hexa code"
		}
		return data
	}
	p += pPreamble

	data.Meta.DataFieldLength = int(binary.BigEndian.Uint32(hexBytes[p : p+pDataFieldLength])) // size is calculated starting from Codec ID to Number of Data 2
	p += pDataFieldLength

	// case where check is performed to see if hex length is equally correct and not currupted
	totalSize := data.Meta.DataFieldLength + pPreamble + pDataFieldLength + pCRC16
	if totalSize != len(hexBytes) {
		var remarks string
		if totalSize > len(hexBytes) {
			remarks = "Hex string received is shorter"
		} else {
			remarks = "Hex string received is larger"
		}
		data.Meta.IsValid = false
		data.Meta.Message = "DataFieldLength failed to verify the string. " + remarks
		return data
	}

	data.CodecID = strings.ToUpper(hex.EncodeToString([]byte{hexBytes[p]}))

	// case where checking is needed to perform against CRC16 module
	if !bytes.Equal(hexBytes[len(hexBytes)-pCRC16:], helper.CRC16(hexBytes[p:len(hexBytes)-pCRC16])) {
		data.Meta.IsValid = false
		data.Meta.Message = "CRC checksum failed"
		return data
	}

	p += pCodecID

	switch data.CodecID {
	case "08": // Codec8
		var err error
		data.AVLData, err = codec8.Decode(hexBytes[p : len(hexBytes)-pCRC16])
		if err != nil {
			data.Meta.IsValid = false
			data.Meta.Message = err.Error()
		}
		break
	case "8E": // Codec8E
		var err error
		data.AVLData, err = codec8e.Decode(hexBytes[p : len(hexBytes)-pCRC16])
		if err != nil {
			data.Meta.IsValid = false
			data.Meta.Message = err.Error()
		}
		break
	case "10": // Codec16
		var err error
		data.AVLData, err = codec16.Decode(hexBytes[p : len(hexBytes)-pCRC16])
		if err != nil {
			data.Meta.IsValid = false
			data.Meta.Message = err.Error()
		}
		break
	case "0C": // Codec12
		var err error
		data.ResponseString, err = codec12.Decode(hexBytes[p : len(hexBytes)-pCRC16])
		if err != nil {
			data.Meta.IsValid = false
			data.Meta.Message = err.Error()
		}
		break
	case "0D": // Codec13
		var err error
		var timeInst time.Time
		data.ResponseString, timeInst, err = codec13.Decode(hexBytes[p : len(hexBytes)-pCRC16])
		if err != nil {
			data.Meta.IsValid = false
			data.Meta.Message = err.Error()
			data.Meta.ServerTime = timeInst.Format(time.RFC3339)
		}
		break
	case "0E": // Codec14
		var err error
		data.ResponseString, data.IMEI, err = codec14.Decode(hexBytes[p : len(hexBytes)-pCRC16])
		if err != nil {
			data.Meta.IsValid = false
			data.Meta.Message = err.Error()
		}
		break
	default:
		data.Meta.IsValid = false
		data.Meta.Message = "Hex decoding codec not supported"
		break
	}

	return data
}
