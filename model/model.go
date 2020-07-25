package model

// MetaType is the meta field in json specifying the data about parsed string
type MetaType struct {
	DataFieldLength int    `json:"dataFieldLength"`
	IsValid         bool   `json:"isValid"`
	HexStringLength int    `json:"hexStringLength"`
	Message         string `json:"message"`
	HexString       string `json:"hexString"`
	ServerTime      string `json:"serverTime"`
}

// GPSType is the data mapping of GPS data with hex string
type GPSType struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Altitude  int     `json:"altitude"`
	Angle     int     `json:"angle"`
	Satellite int     `json:"satellite"`
	Speed     int     `json:"speed"`
}

// IOElementType is the IO section of the received IO
type IOElementType struct {
	EventIOID      int                 `json:"eventID"`
	TotalIO        int                 `json:"totalIO"`
	GenerationType int                 `json:"generationType,omitempty"`
	IOEvents       map[string][]IOType `json:"ioEvents"`
}

// IOType is the single instance of N byte of IO event
type IOType struct {
	ID     int `json:"id"`
	Value  int `json:"value"`
	Length int `json:"length,omitempty"`
}

// AVLType is the AVL structure for Codec8, Codec8E and Codec16
type AVLType struct {
	Time      string        `json:"time"`
	Priority  int           `json:"priority"`
	GPSData   GPSType       `json:"gps"`
	IOElement IOElementType `json:"ioElement"`
}

// DecodedCodec is the compatible struct to handle Codec8, Codec8E and Codec16
type DecodedCodec struct {
	ID             string    `json:"id"`
	Meta           MetaType  `json:"_meta"`
	CodecID        string    `json:"codecId,omitempty"`
	AVLData        []AVLType `json:"avl,omitempty"`
	ResponseString string    `json:"response,omitempty"`
	IMEI           string    `json:"deviceIMEI,omitempty"`
}

// EncodedCodec is the representation of hex string from a particular command
type EncodedCodec struct {
	Meta     MetaType `json:"_meta"`
	CodecID  string   `json:"codecId"`
	CodecHex string   `json:"hex"`
	IMEI     string   `json:"deviceIMEI,omitempty"`
}

// CodecConfig handles the configuration to tune up between Codec8, Codec8E and Codec16
type CodecConfig struct {
	NumberOfData2 int
	NumberOfData1 int
	// GPS Byte mappings
	Timestamp int
	Priority  int
	Longitude int
	Latitude  int
	Altitude  int
	Angle     int
	Satellite int
	Speed     int
	// IO Byte mappings
	EventID          int
	TotalIO          int
	TotalOneByteIO   int
	OneByteIOID      int
	OneByteIOValue   int
	TotalTwoByteIO   int
	TwoByteIOID      int
	TwoByteIOValue   int
	TotalFourByteIO  int
	FourByteIOID     int
	FourByteIOValue  int
	TotalEightByteIO int
	EightByteIOID    int
	EightByteIOValue int
}
