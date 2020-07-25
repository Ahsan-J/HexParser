package constant

type avlIOElement struct {
	Bytes       byte
	ID          int
	Type        int
	ValueMin    int
	ValueMax    int
	Multiplier  float32
	Units       string
	Description string
}

var a = map[string]int{"foo": 1, "bar": 2}
