package errmap

import "errors"

var (
	ErrmapConflict = errors.New("conflict")
	ErrmapNotFound = errors.New("not found")
	ErrmapInvalidPassword = errors.New("invalid password")
	ErrmapInvalidStock = errors.New("invalid stock")
)