package codec14

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Ahsan-J/HexParser/helper"
	"github.com/Ahsan-J/HexParser/model"
)

// Encode takes the string command and returns the Hex representation of the command in Codec12
func Encode(command string, imei int) (model.EncodedCodec, error) {
	var codec model.EncodedCodec
	var pIMEIByte = 8
	imeiString := strings.Replace(strconv.Itoa(imei), " ", "", -1)
	var imeiHex string

	if len(imeiString) == 15 {
		imeiHex = "0" + imeiString
	} else {
		imeiHex = imeiString
	}

	if imeiString == "" || len(imeiHex)/2 != pIMEIByte {
		codec.Meta.IsValid = false
		codec.Meta.Message = "Invalid Device IMEI number"
		return codec, errors.New("Invalid Device IMEI number")
	}

	var dataString = ""
	commandHex := hex.EncodeToString([]byte(command))

	commandSize := make([]byte, 4)
	binary.BigEndian.PutUint32(commandSize, uint32((len(commandHex)/2)+pIMEIByte))

	premable := hex.EncodeToString([]byte{0, 0, 0, 0}) // preamble

	dataString += "0E"                            // CodecId = 12
	dataString += "01"                            // Command quantity static to sending only one
	dataString += "05"                            // Type: Command = 05
	dataString += hex.EncodeToString(commandSize) // command size
	dataString += imeiHex                         // device imei
	dataString += commandHex                      // command of X byte
	dataString += "01"                            // Command quantity 2 static to sending only one

	dataHex, _ := hex.DecodeString(dataString)
	CRCString := hex.EncodeToString(helper.CRC16(dataHex))

	DataFieldLength := fmt.Sprintf("%X", len(dataHex))
	finalHex := strings.ToUpper(premable + DataFieldLength + dataString + CRCString)

	codec.Meta.DataFieldLength = len(dataHex)
	codec.Meta.HexStringLength = len(finalHex)
	codec.Meta.IsValid = true
	codec.Meta.ServerTime = time.Now().Format(time.RFC3339)
	codec.CodecHex = finalHex
	codec.CodecID = "0E"
	codec.IMEI = imeiString

	return codec, nil
}
