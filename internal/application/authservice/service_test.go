package authservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"media-equipment-tracker/internal/application/authservice"
	"media-equipment-tracker/internal/domain"
)

type AuthSuite struct {
	suite.Suite

	svc authservice.AuthUserService

	tokenMaker *TokenMakerMock
	hasher     *HasherMock
	userRep    *UserRepoMock
	tokenRep   *TokenRepoMock
}

func (s *AuthSuite) SetupTest() {
	s.tokenMaker = new(TokenMakerMock)
	s.hasher = new(HasherMock)
	s.userRep = new(UserRepoMock)
	s.tokenRep = new(TokenRepoMock)

	svc, err := authservice.NewAuthUser(
		s.tokenMaker, s.hasher, time.Hour, s.userRep, s.tokenRep,
	)
	s.Require().NoError(err)
	s.svc = svc
}

func newUser(admin bool) *domain.User {
	return &domain.User{
		ID:           uuid.New(),
		Email:        "test@mail.ru",
		HashPassword: "hash",
		IsAdmin:      admin,
	}
}

func (s *AuthSuite) TestNewAuthUser_Invalid() {
	svc, err := authservice.NewAuthUser(nil, nil, 0, nil, nil)
	s.Nil(svc)
	s.ErrorIs(err, authservice.ErrCreateAuthUserService)
}

func (s *AuthSuite) TestLoginUser_Success() {
	u := newUser(true)

	s.userRep.On("GetByEmail", u.Email, mock.Anything).Return(u, nil)
	s.hasher.On("CheckPassword", "pass", u.HashPassword).Return(nil)

	s.tokenMaker.On("CreateToken", mock.MatchedBy(func(p *domain.TokenPayload) bool {
		return p.UserID == u.ID && len(p.Roles) == 1
	})).Return("token", nil)

	s.tokenRep.On("Add", "token").Return(true)

	token, err := s.svc.LoginUser(context.Background(), authservice.LoginUserRequest{
		Email: u.Email, Password: "pass",
	})

	s.NoError(err)
	s.Equal("token", token)
}

func (s *AuthSuite) TestLoginUser_AddFail() {
	u := newUser(false)

	s.userRep.On("GetByEmail", u.Email, mock.Anything).Return(u, nil)
	s.hasher.On("CheckPassword", mock.Anything, mock.Anything).Return(nil)
	s.tokenMaker.On("CreateToken", mock.Anything).Return("token", nil)
	s.tokenRep.On("Add", "token").Return(false)

	_, err := s.svc.LoginUser(context.Background(), authservice.LoginUserRequest{Email: u.Email})
	s.Error(err)
}

func (s *AuthSuite) TestRegisterUser() {
	s.hasher.On("HashPassword", "pass").Return("hash", nil)
	s.userRep.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	err := s.svc.RegisterUser(context.Background(), authservice.RegisterUserRequest{
		FullName: "test", Email: "mail", Password: "pass",
	})

	s.NoError(err)
}

func (s *AuthSuite) TestVerifyByToken_Success() {
	payload := &domain.TokenPayload{Roles: []domain.RoleAuth{}}

	s.tokenMaker.On("VerifyToken", "token").Return(payload, nil)

	res, err := s.svc.VerifyByToken("token", []domain.RoleAuth{domain.AdminRole})

	s.NoError(err)
	s.Equal(payload, res)
}

func (s *AuthSuite) TestVerifyByToken_BadRole() {
	payload := &domain.TokenPayload{Roles: []domain.RoleAuth{domain.AdminRole}}

	s.tokenMaker.On("VerifyToken", "token").Return(payload, nil)

	_, err := s.svc.VerifyByToken("token", []domain.RoleAuth{domain.AdminRole})

	s.ErrorIs(err, authservice.ErrIncorrectRole)
}

func (s *AuthSuite) TestLogoutUser() {
	s.tokenRep.On("Delete", "token").Return()

	s.svc.LogoutUser(context.Background(), "token")

	s.tokenRep.AssertCalled(s.T(), "Delete", "token")
}

type TokenMakerMock struct{ mock.Mock }

func (m *TokenMakerMock) CreateToken(p *domain.TokenPayload) (string, error) {
	args := m.Called(p)
	return args.String(0), args.Error(1)
}
func (m *TokenMakerMock) VerifyToken(t string) (*domain.TokenPayload, error) {
	args := m.Called(t)
	if v := args.Get(0); v != nil {
		return v.(*domain.TokenPayload), args.Error(1)
	}
	return nil, args.Error(1)
}

type HasherMock struct{ mock.Mock }

func (m *HasherMock) HashPassword(p string) (string, error) {
	args := m.Called(p)
	return args.String(0), args.Error(1)
}
func (m *HasherMock) CheckPassword(p, h string) error {
	return m.Called(p, h).Error(0)
}

type UserRepoMock struct{ mock.Mock }

func (m *UserRepoMock) GetByEmail(email string, _ ...domain.UserOption) (*domain.User, error) {
	args := m.Called(email, mock.Anything)
	if v := args.Get(0); v != nil {
		return v.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *UserRepoMock) Create(u *domain.User) error {
	return m.Called(u).Error(0)
}

// unused
func (*UserRepoMock) Get(uuid.UUID, ...domain.UserOption) (*domain.User, error) { return nil, nil }
func (*UserRepoMock) List(...domain.UserOption) ([]*domain.User, error)         { return nil, nil }
func (*UserRepoMock) Reload(*domain.User, ...domain.UserOption) error           { return nil }
func (*UserRepoMock) Update(*domain.User) error                                 { return nil }
func (*UserRepoMock) Delete(uuid.UUID) error                                    { return nil }

type TokenRepoMock struct{ mock.Mock }

func (m *TokenRepoMock) Add(t string) bool { return m.Called(t).Bool(0) }
func (m *TokenRepoMock) Delete(t string)   { m.Called(t) }
func (*TokenRepoMock) Check(string) bool   { return false }

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
