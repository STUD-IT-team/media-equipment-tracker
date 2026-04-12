//go:build integration

package postgresuser_test

import (
	"media-equipment-tracker/internal/domain"

	"github.com/google/uuid"
)

func newUser() *domain.User {
	return &domain.User{
		ID:           uuid.New(),
		FullName:     "test",
		Email:        uuid.NewString() + "@mail.ru",
		HashPassword: "hash",
		Nice:         domain.DefaultNice,
	}
}

func newOrg() *domain.Organization {
	return &domain.Organization{
		ID:   uuid.New(),
		Name: "org",
	}
}

func newDept() *domain.Department {
	return &domain.Department{
		ID:   uuid.New(),
		Name: "dept",
	}
}
