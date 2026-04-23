package errs

import (
	"errors"
	"fmt"
	"media-equipment-tracker/internal/domain"
	"strings"
)

type EntityNotFoundError struct {
	Entity string
	ID     interface{}
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %v not found", e.Entity, e.ID)
}

func NewEntityNotFoundError(entity string, id interface{}) EntityNotFoundError {
	return EntityNotFoundError{Entity: entity, ID: id}
}

type EntityAlreadyExistsError struct {
	Entity string
	ID     interface{}
}

func (e EntityAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s with ID %v already exists", e.Entity, e.ID)
}

func NewEntityAlreadyExistsError(entity string, id interface{}) EntityAlreadyExistsError {
	return EntityAlreadyExistsError{Entity: entity, ID: id}
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

func NewValidationError(field, message string) ValidationError {
	return ValidationError{Field: field, Message: message}
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	msgs := make([]string, 0, len(ve))
	for _, v := range ve {
		msgs = append(msgs, v.Error())
	}
	return fmt.Sprintf("validation errors: %s", strings.Join(msgs, "; "))
}

func NewValidationErrors(errs ...ValidationError) ValidationErrors {
	return ValidationErrors(errs)
}

type ConcurrencyError struct {
	Entity string
	ID     interface{}
}

func (e ConcurrencyError) Error() string {
	return fmt.Sprintf("concurrency conflict on %s with ID %v", e.Entity, e.ID)
}

func NewConcurrencyError(entity string, id interface{}) ConcurrencyError {
	return ConcurrencyError{Entity: entity, ID: id}
}

type RepositoryError struct {
	Op  string
	Err error
}

func (e RepositoryError) Error() string {
	return fmt.Sprintf("repository error during %s: %v", e.Op, e.Err)
}

func (e RepositoryError) Unwrap() error {
	return e.Err
}

func NewRepositoryError(op string, err error) RepositoryError {
	return RepositoryError{Op: op, Err: err}
}

type RoleAuthError struct {
	RequiredRoles []domain.RoleAuth
	ActualRoles   []domain.RoleAuth
}

func (e RoleAuthError) Error() string {
	return fmt.Sprintf("insufficient role authorization. required roles: %v, actual roles: %v", e.RequiredRoles, e.ActualRoles)
}

func NewRoleAuthError(requiredRoles []domain.RoleAuth, actualRoles []domain.RoleAuth) RoleAuthError {
	return RoleAuthError{RequiredRoles: requiredRoles, ActualRoles: actualRoles}
}

func IsRoleAuthError(err error) bool {
	var e RoleAuthError
	return errors.As(err, &e)
}

func IsEntityNotFoundError(err error) bool {
	var e EntityNotFoundError
	return errors.As(err, &e)
}

func IsEntityAlreadyExistsError(err error) bool {
	var e EntityAlreadyExistsError
	return errors.As(err, &e)
}

func IsValidationError(err error) bool {
	var ve ValidationErrors
	if errors.As(err, &ve) {
		return true
	}
	var v ValidationError
	return errors.As(err, &v)
}

func IsConcurrencyError(err error) bool {
	var e ConcurrencyError
	return errors.As(err, &e)
}

func IsRepositoryError(err error) bool {
	var e RepositoryError
	return errors.As(err, &e)
}
