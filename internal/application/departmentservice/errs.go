package departmentservice

import (
	"errors"
)

var (
	ErrDepartmentHasUsers       = errors.New("department has users")
	ErrDepartmentHasInvocations = errors.New("department has invocations")
)
