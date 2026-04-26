package accessservice

import "errors"

type EquipmentAccessError struct {
	Message string
}

func (e EquipmentAccessError) Error() string {
	return e.Message
}

func NewEquipmentAccessError(message string) EquipmentAccessError {
	return EquipmentAccessError{
		Message: message,
	}
}

func IsEquipmentAccessError(err error) bool {
	return errors.Is(err, EquipmentAccessError{})
}
