package dto

import (
	"fmt"
	"media-equipment-tracker/internal/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Nice      int       `json:"nice"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Nice:     user.Nice,
		IsAdmin:  user.IsAdmin,
	}
}

type UserInvocationsResponse struct {
	EquipmentInvocations []EquipmentInvocationResponse `json:"equipment_invocations"`
	StudioInvocations    []StudioInvocationResponse    `json:"studio_invocations"`
}

func UserToUserInvocationsResponse(user *domain.User) UserInvocationsResponse {
	equipmentInvocations := make([]EquipmentInvocationResponse, len(user.EquipmentInvocations))
	studioInvocations := make([]StudioInvocationResponse, len(user.StudioInvocations))
	for i := range user.EquipmentInvocations {
		equipmentInvocations[i] = EquipmentInvocationToResponse(user.EquipmentInvocations[i])
	}
	for i := range user.StudioInvocations {
		studioInvocations[i] = StudioInvocationToResponse(user.StudioInvocations[i])
	}
	return UserInvocationsResponse{
		EquipmentInvocations: equipmentInvocations,
		StudioInvocations:    studioInvocations,
	}
}

type EquipmentInvocationResponse struct {
	ID        uuid.UUID `json:"id"`
	EventName string    `json:"eventName"`

	Status              string     `json:"status"`
	StartTime           time.Time  `json:"startTime"`
	EndTime             time.Time  `json:"endTime"`
	EquipmentReturnTime *time.Time `json:"equipmentReturnTime,omitempty"`
	SdCardReturnTime    *time.Time `json:"sdCardReturnTime,omitempty"`
	CuratorComment      string     `json:"curatorComment,omitempty"`

	UserID           uuid.UUID  `json:"userID"`
	UserName         string     `json:"userName,omitempty"`
	AdminID          *uuid.UUID `json:"adminID,omitempty"`
	AdminName        string     `json:"adminName,omitempty"`
	OrganizationID   *uuid.UUID `json:"organizationID,omitempty"`
	OrganizationName string     `json:"organizationName,omitempty"`
	DepartmentID     *uuid.UUID `json:"departmentID,omitempty"`
	DepartmentName   string     `json:"departmentName,omitempty"`

	EquipmentCount int `json:"equipmentCount"`
	// created_at updated_at
}

type StudioInvocationResponse struct {
	ID                  uuid.UUID `json:"id"`
	EventName           string    `json:"eventName"`
	ShootingDescription string    `json:"shootingDescription"`
	UserID              uuid.UUID `json:"userID"`
	Status              string    `json:"status"`
	StartTime           time.Time `json:"startTime"`
	EndTime             time.Time `json:"endTime"`
	NeedsChromakey      bool      `json:"needsChromakey"`
	NeedsCyclorama      bool      `json:"needsCyclorama"`
	NeedsBlackFabric    bool      `json:"needsBlackFabric"`
	CuratorComment      string    `json:"curatorComment,omitempty"`

	UserName         string     `json:"userName,omitempty"`
	AdminID          *uuid.UUID `json:"adminID,omitempty"`
	AdminName        string     `json:"adminName,omitempty"`
	OrganizationID   *uuid.UUID `json:"organizationID,omitempty"`
	OrganizationName string     `json:"organizationName,omitempty"`
	DepartmentID     *uuid.UUID `json:"departmentID,omitempty"`
	DepartmentName   string     `json:"departmentName,omitempty"`

	// created_at updated_at
}

func EquipmentInvocationToResponse(inv *domain.EquipmentInvocation) EquipmentInvocationResponse {
	resp := EquipmentInvocationResponse{
		ID:                  inv.ID,
		EventName:           inv.EventName,
		StartTime:           inv.StartTime,
		EndTime:             inv.EndTime,
		EquipmentReturnTime: inv.EquipmentReturnTime,
		SdCardReturnTime:    inv.SdCardReturnTime,
		Status:              string(inv.Status),
		CuratorComment:      inv.CuratorComment,
		UserID:              inv.UserID,
		EquipmentCount:      len(inv.Equipment),
	}

	if inv.Organization != nil {
		resp.OrganizationID = &inv.Organization.ID
		resp.OrganizationName = inv.Organization.Name
	} else if inv.OrganizationID != nil {
		resp.OrganizationID = inv.OrganizationID
	}

	if inv.Department != nil {
		resp.DepartmentID = &inv.Department.ID
		resp.DepartmentName = inv.Department.Name
	} else if inv.DepartmentID != nil {
		resp.DepartmentID = inv.DepartmentID
	}

	if inv.User != nil {
		resp.UserName = inv.User.FullName
	}

	if inv.Admin != nil {
		resp.AdminID = &inv.Admin.ID
		resp.AdminName = inv.Admin.FullName
	} else if inv.AdminID != nil {
		resp.AdminID = inv.AdminID
	}

	return resp
}

func StudioInvocationToResponse(inv *domain.StudioInvocation) StudioInvocationResponse {
	resp := StudioInvocationResponse{
		ID:                  inv.ID,
		EventName:           inv.EventName,
		ShootingDescription: inv.ShootingDescription,
		UserID:              inv.UserID,
		Status:              string(inv.Status),
		StartTime:           inv.StartTime,
		EndTime:             inv.EndTime,
		NeedsChromakey:      inv.NeedsChromakey,
		NeedsCyclorama:      inv.NeedsCyclorama,
		NeedsBlackFabric:    inv.NeedsBlackFabric,
		CuratorComment:      inv.CuratorComment,
	}

	if inv.Organization != nil {
		resp.OrganizationID = &inv.Organization.ID
		resp.OrganizationName = inv.Organization.Name
	} else if inv.OrganizationID != nil {
		resp.OrganizationID = inv.OrganizationID
	}

	if inv.Department != nil {
		resp.DepartmentID = &inv.Department.ID
		resp.DepartmentName = inv.Department.Name
	} else if inv.DepartmentID != nil {
		resp.DepartmentID = inv.DepartmentID
	}

	if inv.User != nil {
		resp.UserName = inv.User.FullName
	}

	if inv.Admin != nil {
		resp.AdminID = &inv.Admin.ID
		resp.AdminName = inv.Admin.FullName
	} else if inv.AdminID != nil {
		resp.AdminID = inv.AdminID
	}

	return resp
}

type InvocationFilterDto struct {
	Status     string `form:"status" json:"status" enums:"created,under_review,changes_required,approved,equipment_issued,equipment_returned,completed,cancelled"`
	Type       string `form:"type,default=all" json:"type" enums:"equipment,studio,all"`
	ActiveOnly bool   `form:"active_only,default=false" json:"active_only"`
}

type UpdateUserDto struct {
	FullName string `json:"full_name" example:"Иванов Иван Иванович"`
	Email    string `json:"email" example:"ivan@example.com"`
}

func UpdateUserDtoToUser(dto UpdateUserDto) (*domain.User, error) {
	user := &domain.User{
		FullName: dto.FullName,
		Email:    dto.Email,
	}
	return user, nil
}

type UpdateUserByAdminDto struct {
	FullName string `json:"full_name" example:"Иванов Иван Иванович"`
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `form:"password" json:"password"`
	Nice     int    `form:"nice" json:"nice"`
	IsAdmin  bool   `form:"is_admin" json:"is_admin"`
	// TODO organizations, departments
}

func UpdateUserByAdminDtoToUser(dto UpdateUserByAdminDto) (*domain.User, error) {
	user := &domain.User{
		FullName: dto.FullName,
		Email:    dto.Email,
		Nice:     dto.Nice,
		IsAdmin:  dto.IsAdmin,
	}
	if dto.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.HashPassword = string(hashedPassword)
	}
	return user, nil
}

// ----------------------------------

type UserFilterDto struct {
	Search  string `form:"search" json:"search" description:"Поиск по ФИО/email"`
	IsAdmin *bool  `form:"isAdmin" json:"isAdmin"`
	NiceMin int    `form:"niceMin" json:"niceMin"`
	NiceMax int    `form:"niceMax" json:"niceMax"`
	Page    int    `form:"page,default=1" json:"page"`
	Limit   int    `form:"limit,default=20" json:"limit"`
}
type UserDetailResponse struct {
	UserResponse
	Organizations []OrganizationBriefResponse `json:"organizations,omitempty"`
	Departments   []UserDepartmentResponse    `json:"departments,omitempty"`
}

type OrganizationBriefResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UserDepartmentResponse struct {
	DepartmentID   uuid.UUID `json:"departmentID"`
	DepartmentName string    `json:"departmentName"`
	Role           string    `json:"role"` // trainee, activist
}

// UpdateUserAdminDto - обновление пользователя администратором
type UpdateUserAdminDto struct {
	FullName string `json:"fullName" example:"Иванов Иван Иванович"`
	Email    string `json:"email" example:"ivan@example.com"`
	IsAdmin  *bool  `json:"isAdmin,omitempty" description:"Права администратора"`
}

// UpdateNiceDto - изменение баллов порядочности
type UpdateNiceDto struct {
	Nice   int    `json:"nice" binding:"required,min=1" example:"150"`
	Reason string `json:"reason,omitempty" example:"За активное участие в мероприятии"`
}

// UserListResponse - список пользователей с пагинацией
type UserListResponse struct {
	Items      []UserResponse `json:"items"`
	Total      int64          `json:"total,omitempty"`
	Page       int            `json:"page,omitempty"`
	Limit      int            `json:"limit,omitempty"`
	TotalPages int            `json:"totalPages,omitempty"`
}

// CreateUserRequest - создание пользователя (если потребуется)
type CreateUserRequest struct {
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	IsAdmin  bool   `json:"isAdmin"`
}

// EquipmentInInvocationResponse - ответ с информацией о оборудовании в заявке
type EquipmentInInvocationResponse struct {
	InvocationID uuid.UUID `json:"invocationID"`
	EquipmentID  uuid.UUID `json:"equipmentID"`
	Status       string    `json:"status"` // pending, issued, returned
}

// UserToUserDetailResponse преобразует User в UserDetailResponse с вложенными связями
func UserToUserDetailResponse(user *domain.User) UserDetailResponse {
	resp := UserDetailResponse{
		UserResponse: UserToUserResponse(user),
	}

	// Маппинг организаций
	if len(user.Organizations) > 0 {
		resp.Organizations = make([]OrganizationBriefResponse, len(user.Organizations))
		for i, org := range user.Organizations {
			resp.Organizations[i] = OrganizationToBriefResponse(org)
		}
	}

	// Маппинг отделов
	if len(user.Departments) > 0 {
		resp.Departments = make([]UserDepartmentResponse, len(user.Departments))
		for i, ud := range user.Departments {
			resp.Departments[i] = UserDepartmentToResponse(ud)
		}
	}

	return resp
}

// OrganizationToBriefResponse преобразует Organization в OrganizationBriefResponse
func OrganizationToBriefResponse(org *domain.Organization) OrganizationBriefResponse {
	return OrganizationBriefResponse{
		ID:   org.ID,
		Name: org.Name,
	}
}

// UserDepartmentToResponse преобразует UserDepartment в UserDepartmentResponse
func UserDepartmentToResponse(ud *domain.UserDepartment) UserDepartmentResponse {
	departmentName := ""
	if ud.Department != nil {
		departmentName = ud.Department.Name
	}

	return UserDepartmentResponse{
		DepartmentID:   ud.DepartmentID,
		DepartmentName: departmentName,
		Role:           string(ud.Role),
	}
}

// UsersToUserListResponse преобразует слайс пользователей в UserListResponse
func UsersToUserListResponse(users []*domain.User, total int64, page, limit int) UserListResponse {
	items := make([]UserResponse, len(users))
	for i, user := range users {
		items[i] = UserToUserResponse(user)
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	return UserListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// EquipmentInInvocationToResponse преобразует EquipmentInInvocation в EquipmentInInvocationResponse
func EquipmentInInvocationToResponse(eqi *domain.EquipmentInInvocation) EquipmentInInvocationResponse {
	resp := EquipmentInInvocationResponse{
		InvocationID: eqi.InvocationID,
		EquipmentID:  eqi.EquipmentID,
		Status:       string(eqi.Status),
	}
	return resp
}

// EquipmentInvocationsToResponses преобразует слайс EquipmentInvocation в слайс EquipmentInvocationResponse
func EquipmentInvocationsToResponses(invs []*domain.EquipmentInvocation) []EquipmentInvocationResponse {
	if invs == nil {
		return []EquipmentInvocationResponse{}
	}

	responses := make([]EquipmentInvocationResponse, len(invs))
	for i, inv := range invs {
		responses[i] = EquipmentInvocationToResponse(inv)
	}
	return responses
}

// StudioInvocationsToResponses преобразует слайс StudioInvocation в слайс StudioInvocationResponse
func StudioInvocationsToResponses(invs []*domain.StudioInvocation) []StudioInvocationResponse {
	if invs == nil {
		return []StudioInvocationResponse{}
	}

	responses := make([]StudioInvocationResponse, len(invs))
	for i, inv := range invs {
		responses[i] = StudioInvocationToResponse(inv)
	}
	return responses
}

// UserInvocationsToResponse преобразует бронирования пользователя в UserInvocationsResponse
func UserInvocationsToResponse(equipmentInvs []*domain.EquipmentInvocation, studioInvs []*domain.StudioInvocation) UserInvocationsResponse {
	return UserInvocationsResponse{
		EquipmentInvocations: EquipmentInvocationsToResponses(equipmentInvs),
		StudioInvocations:    StudioInvocationsToResponses(studioInvs),
	}
}

// UpdateUserAdminDtoToUser обновляет существующего пользователя из UpdateUserAdminDto
func UpdateUserAdminDtoToUser(dto *UpdateUserAdminDto, user *domain.User) {
	if dto.FullName != "" {
		user.FullName = dto.FullName
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.IsAdmin != nil {
		user.IsAdmin = *dto.IsAdmin
	}
}

// CreateUserRequestToUser преобразует CreateUserRequest в User
func CreateUserRequestToUser(req *CreateUserRequest) *domain.User {
	return &domain.User{
		FullName:     req.FullName,
		Email:        req.Email,
		HashPassword: "", // будет установлен после хэширования
		Nice:         domain.DefaultNice,
		IsAdmin:      req.IsAdmin,
	}
}

// UserFilterDtoToUserFilter преобразует UserFilterDto в domain.UserFilter
func UserFilterDtoToUserFilter(dto *UserFilterDto) domain.UserFilter {
	filter := domain.UserFilter{}

	if dto.Search != "" {
		filter.FullName = dto.Search
		filter.Email = dto.Search
	}

	return filter
}
