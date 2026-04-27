package organizationservice

import (
	"errors"
)

var (
	ErrOrganizationHasUsers       = errors.New("organization has users")
	ErrOrganizationHasInvocations = errors.New("organization has invocations")
)
