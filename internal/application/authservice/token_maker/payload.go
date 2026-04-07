package tokenmaker

import (
	"time"

	"github.com/google/uuid"
)

// данные полезной нагрузки, хранящиеся внутри тела токена
type RoleAuth string

const (
	UserRole RoleAuth = "user_role"
)

type Payload struct {
	PersonID  uuid.UUID `json:"person_id"`
	Role      RoleAuth  `json:"role"`
	ExpiredAt time.Time `json:"expired_at"` // время когда срок действия токена истечет
}

func NewPayload(personID uuid.UUID, role RoleAuth, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		PersonID:  personID,
		Role:      role,
		ExpiredAt: time.Now().UTC().Add(duration),
	}
	return payload, nil
}

func CreatePayload(personID uuid.UUID, role RoleAuth, expiredAt time.Time) (*Payload, error) {
	payload := &Payload{
		PersonID:  personID,
		Role:      role,
		ExpiredAt: expiredAt,
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (p *Payload) GetPersonID() uuid.UUID {
	return p.PersonID
}

func (p *Payload) GetRole() RoleAuth {
	return p.Role
}

func (p *Payload) GetExpiredAt() time.Time {
	return p.ExpiredAt
}
