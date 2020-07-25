package codec8e

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
	var pEventID = 2
	var pTotalIO = 2
	var pTotalOneByteIO = 2
	var pOneByteIOID = 2
	var pOneByteIOValue = 1
	var pTotalTwoByteIO = 2
	var pTwoByteIOID = 2
	var pTwoByteIOValue = 2
	var pTotalFourByteIO = 2
	var pFourByteIOID = 2
	var pFourByteIOValue = 4
	var pTotalEightByteIO = 2
	var pEightByteIOID = 2
	var pEightByteIOValue = 8
	var pTotalXBytesIO = 2
	var pXBytesIOID = 2
	var pXBytesIOLength = 2

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
		eventID := int(binary.BigEndian.Uint16(hexBytes[p : p+pEventID]))
		AVLData.IOElement.EventIOID = eventID
		p += pEventID

		totalIO := int(binary.BigEndian.Uint16(hexBytes[p : p+pTotalIO]))
		AVLData.IOElement.TotalIO = totalIO
		p += pTotalIO // N = N1 + N2 + N4 + N8

		totalOneByteIO := int(binary.BigEndian.Uint16(hexBytes[p : p+pTotalOneByteIO]))
		p += pTotalOneByteIO

		for i := 0; i < totalOneByteIO; i++ { // Iterating through One Byte event
			var io model.IOType

			io.ID = int(binary.BigEndian.Uint16(hexBytes[p : p+pOneByteIOID]))
			p += pOneByteIOID

			io.Value = int(hexBytes[p])
			p += pOneByteIOValue

			IOMap["1"] = append(IOMap["1"], io)
			// AVLData.IOElement.OneBytes = append(AVLData.IOElement.OneBytes, io)
		}

		totalTwoByteIO := int(binary.BigEndian.Uint16(hexBytes[p : p+pTotalTwoByteIO]))
		p += pTotalTwoByteIO

		for i := 0; i < totalTwoByteIO; i++ { // Iterating through Two Bytes event
			var io model.IOType

			io.ID = int(binary.BigEndian.Uint16(hexBytes[p : p+pTwoByteIOID]))
			p += pTwoByteIOID

			io.Value = int(binary.BigEndian.Uint16(hexBytes[p : p+pTwoByteIOValue]))
			p += pTwoByteIOValue

			IOMap["2"] = append(IOMap["2"], io)
			// AVLData.IOElement.TwoBytes = append(AVLData.IOElement.TwoBytes, io)
		}

		totalFourByteIO := int(binary.BigEndian.Uint16(hexBytes[p : p+pTotalFourByteIO]))
		p += pTotalFourByteIO

		for i := 0; i < totalFourByteIO; i++ { // Iterating through Four Bytes event
			var io model.IOType
			io.ID = int(binary.BigEndian.Uint16(hexBytes[p : p+pFourByteIOID]))
			p += pFourByteIOID

			io.Value = int(binary.BigEndian.Uint32(hexBytes[p : p+pFourByteIOValue]))
			p += pFourByteIOValue

			IOMap["4"] = append(IOMap["4"], io)
			// AVLData.IOElement.FourBytes = append(AVLData.IOElement.FourBytes, io)
		}

		totalEightByteIO := int(binary.BigEndian.Uint16(hexBytes[p : p+pTotalEightByteIO]))
		p += pTotalEightByteIO

		for i := 0; i < totalEightByteIO; i++ { // Iterating through Eight Bytes
			var io model.IOType

			io.ID = int(binary.BigEndian.Uint16(hexBytes[p : p+pEightByteIOID]))
			p += pEightByteIOID

			io.Value = int(binary.BigEndian.Uint32(hexBytes[p : p+pEightByteIOValue]))
			p += pEightByteIOValue

			IOMap["8"] = append(IOMap["8"], io)
			// AVLData.IOElement.EightBytes = append(AVLData.IOElement.EightBytes, io)
		}

		totalXByteIO := int(binary.BigEndian.Uint16(hexBytes[p : p+pTotalXBytesIO]))
		p += pTotalXBytesIO

		for i := 0; i < totalXByteIO; i++ { // Iterating through Eight Bytes
			var io model.IOType

			io.ID = int(binary.BigEndian.Uint16(hexBytes[p : p+pXBytesIOID]))
			p += pXBytesIOID

			io.Length = int(binary.BigEndian.Uint16(hexBytes[p : p+pXBytesIOLength]))

			io.Value = int(binary.BigEndian.Uint32(hexBytes[p : p+io.Length]))
			p += io.Length
			IOMap["x"] = append(IOMap["x"], io)
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
		return nil, errors.New("Something went wrong")
	}

	return AVLDataArray, nil
}
