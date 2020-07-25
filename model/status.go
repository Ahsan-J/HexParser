package model

// Status is the status field used to define a tuple status flag
type Status struct {
	value statusValue
}

// Set modifie the status value
func (status *Status) Set(state statusValue) {

}

// Check returns true if state is set in the value
func (status *Status) Check(state statusValue) bool {
	return false
}

// Unset the status value
func (status *Status) Unset(state statusValue) {

}

// GetValue returns the value representation of the status
func (status *Status) GetValue() int {
	return 0
}
