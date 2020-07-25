package codec8

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/Ahsan-J/HexParser/model"
)

// Decode requires a hexBytes (Array of Hex bytes) and Parses to return AVLType
func Decode(hexBytes []byte) ([]model.AVLType, error) {
	var AVLDataArray []model.AVLType
	var pNumberOfData1 = 1
	var pNumberOfData2 = 1

	// GPS Byte mappings
	var pTimestamp = 8
	var pPriority = 1
	var pLongitude = 4
	var pLatitude = 4
	var pAltitude = 2
	var pAngle = 2
	var pSatellite = 1
	var pSpeed = 2

	// IO Byte mappings
	var pEventID = 1
	var pTotalIO = 1
	var pTotalOneByteIO = 1
	var pOneByteIOID = 1
	var pOneByteIOValue = 1
	var pTotalTwoByteIO = 1
	var pTwoByteIOID = 1
	var pTwoByteIOValue = 2
	var pTotalFourByteIO = 1
	var pFourByteIOID = 1
	var pFourByteIOValue = 4
	var pTotalEightByteIO = 1
	var pEightByteIOID = 1
	var pEightByteIOValue = 8

	p := 0
	numberOfAVLData := int(hexBytes[p])
	p += pNumberOfData1

	for i := 0; i < numberOfAVLData; i++ {
		IOMap := make(map[string][]model.IOType)
		var AVLData model.AVLType
		packetTimeRaw := int64(binary.BigEndian.Uint64(hexBytes[p : p+pTimestamp]))
		AVLData.Time = time.Unix(packetTimeRaw/1000, 0).Format(time.RFC3339)
		p += pTimestamp

		AVLData.Priority = int(hexBytes[p])
		p += pPriority

		// GPS Elements

		AVLData.GPSData.Longitude = float32(int32(binary.BigEndian.Uint32(hexBytes[p:p+pLongitude]))) / 3600000
		p += pLongitude

		AVLData.GPSData.Latitude = float32(int32(binary.BigEndian.Uint32(hexBytes[p:p+pLatitude]))) / 3600000
		p += pLatitude

		AVLData.GPSData.Altitude = int(int16(binary.BigEndian.Uint16(hexBytes[p : p+pAltitude])))
		p += pAltitude

		AVLData.GPSData.Angle = int(int16(binary.BigEndian.Uint16(hexBytes[p : p+pAngle])))
		p += pAngle

		AVLData.GPSData.Satellite = int(hexBytes[p])
		p += pSatellite

		AVLData.GPSData.Speed = int(binary.BigEndian.Uint16(hexBytes[p : p+pSpeed]))
		p += pSpeed

		// if int(AVLData.GPSData.Longitude) == 0 || int(AVLData.GPSData.Latitude) == 0 {
		// 	fmt.Println("Case where the co ordinates are invalid")
		// }

		// if AVLData.GPSData.Speed == 0 {
		// 	fmt.Println("The GPS data is not valid or the car is in rest state")
		// }

		// IO Elements

		// @TODO: check the event id's
		eventID := int(hexBytes[p])
		AVLData.IOElement.EventIOID = eventID
		p += pEventID

		totalIO := int(hexBytes[p])
		AVLData.IOElement.TotalIO = totalIO
		p += pTotalIO // N = N1 + N2 + N4 + N8

		totalOneByteIO := int(hexBytes[p])
		p += pTotalOneByteIO

		for i := 0; i < totalOneByteIO; i++ { // Iterating through One Byte event
			var io model.IOType

			io.ID = int(hexBytes[p])
			p += pOneByteIOID

			io.Value = int(hexBytes[p])
			p += pOneByteIOValue

			IOMap["1"] = append(IOMap["1"], io)
			// AVLData.IOElement.OneBytes = append(AVLData.IOElement.OneBytes, io)
		}

		totalTwoByteIO := int(hexBytes[p])
		p += pTotalTwoByteIO

		for i := 0; i < totalTwoByteIO; i++ { // Iterating through Two Bytes event
			var io model.IOType

			io.ID = int(hexBytes[p])
			p += pTwoByteIOID

			io.Value = int(binary.BigEndian.Uint16(hexBytes[p : p+pTwoByteIOValue]))
			p += pTwoByteIOValue

			IOMap["2"] = append(IOMap["2"], io)
			// AVLData.IOElement.TwoBytes = append(AVLData.IOElement.TwoBytes, io)
		}

		totalFourByteIO := int(hexBytes[p])
		p += pTotalFourByteIO

		for i := 0; i < totalFourByteIO; i++ { // Iterating through Four Bytes event
			var io model.IOType
			io.ID = int(hexBytes[p])
			p += pFourByteIOID

			io.Value = int(binary.BigEndian.Uint32(hexBytes[p : p+pFourByteIOValue]))
			p += pFourByteIOValue

			IOMap["4"] = append(IOMap["4"], io)
			// AVLData.IOElement.FourBytes = append(AVLData.IOElement.FourBytes, io)
		}

		totalEightByteIO := int(hexBytes[p])
		p += pTotalEightByteIO

		for i := 0; i < totalEightByteIO; i++ { // Iterating through Eight Bytes
			var io model.IOType

			io.ID = int(hexBytes[p])
			p += pEightByteIOID

			io.Value = int(binary.BigEndian.Uint32(hexBytes[p : p+pEightByteIOValue]))
			p += pEightByteIOValue

			IOMap["8"] = append(IOMap["8"], io)
			// AVLData.IOElement.EightBytes = append(AVLData.IOElement.EightBytes, io)
		}

		if totalIO != (totalOneByteIO + totalTwoByteIO + totalFourByteIO + totalEightByteIO) {
			return nil, errors.New("Events section was corrupted or invalid")
		}
		AVLData.IOElement.IOEvents = IOMap
		AVLDataArray = append(AVLDataArray, AVLData)
	}

	numberOfData2 := int(hexBytes[p])
	p += pNumberOfData2

	if numberOfAVLData != numberOfData2 {
		return nil, errors.New("AVL Data count is unequal")
	}

	return AVLDataArray, nil
}
