package constant

type statusValue int

type statusList struct {
	ACTIVE statusValue
	BLOCK  statusValue
}

// StatusValues is an enum representation of the values handled by the handler
var StatusValues = &statusList{
	ACTIVE: 1,
	BLOCK:  2,
}
