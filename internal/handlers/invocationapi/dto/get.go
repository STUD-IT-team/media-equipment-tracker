package dto

import (
	"time"

	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeserializeGetInvocationRequest(c *gin.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

type GetInvocationResponse struct {
	ID              string                     `json:"id"`
	EventName       string                     `json:"event_name"`
	StartTime       string                     `json:"start_time"`
	EndTime         string                     `json:"end_time"`
	Status          string                     `json:"status"`
	CuratorComment  string                     `json:"curator_comment"`
	User            *getInvocationUser         `json:"user"`
	Department      *getInvocationDepartment   `json:"department,omitempty"`
	Organization    *getInvocationOrganization `json:"organization,omitempty"`
	Equipment       []getInvocationEquipment   `json:"equipment"`
	IssuedEquipment []getInvocationEquipment   `json:"issued_equipment"`
}

type getInvocationUser struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
	Nice    int    `json:"nice"`
}

type getInvocationDepartment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getInvocationOrganization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getInvocationEquipment struct {
	ID                 string `json:"id"`
	Inventory          string `json:"inventory_number"`
	Name               string `json:"name"`
	ShortName          string `json:"short_name"`
	Category           string `json:"category"`
	Status             string `json:"status"`
	AvailableToTrainee bool   `json:"available_to_trainee"`
}

func SerializeGetInvocationResponse(_ *gin.Context, inv *domain.EquipmentInvocation) any {
	user := getInvocationUser{
		ID:      inv.UserID.String(),
		Name:    inv.User.FullName,
		IsAdmin: inv.User.IsAdmin,
		Nice:    inv.User.Nice,
	}
	var department *getInvocationDepartment
	if inv.Department != nil {
		department = &getInvocationDepartment{
			ID:   inv.DepartmentID.String(),
			Name: inv.Department.Name,
		}
	}
	var organization *getInvocationOrganization
	if inv.Organization != nil {
		organization = &getInvocationOrganization{
			ID:   inv.OrganizationID.String(),
			Name: inv.Organization.Name,
		}
	}

	equipment := make([]getInvocationEquipment, 0, len(inv.Equipment))
	for _, eq := range inv.Equipment {
		equipment = append(equipment, getInvocationEquipment{
			ID:                 eq.EquipmentID.String(),
			Inventory:          eq.Equipment.InventoryNumber,
			Name:               eq.Equipment.Name,
			ShortName:          eq.Equipment.ShortName,
			Category:           eq.Equipment.Category,
			Status:             string(eq.Equipment.Status),
			AvailableToTrainee: eq.Equipment.AvailableToTrainee,
		})
	}

	issuedEquipment := make([]getInvocationEquipment, 0, len(inv.Equipment))
	for _, eq := range inv.Equipment {
		// slog.Info("eq.Status", eq.Status)
		if eq.Status == domain.EquipmentIssued {
			issuedEquipment = append(issuedEquipment, getInvocationEquipment{
				ID:                 eq.EquipmentID.String(),
				Inventory:          eq.Equipment.InventoryNumber,
				Name:               eq.Equipment.Name,
				ShortName:          eq.Equipment.ShortName,
				Category:           eq.Equipment.Category,
				Status:             string(eq.Equipment.Status),
				AvailableToTrainee: eq.Equipment.AvailableToTrainee,
			})
		}
	}

	return GetInvocationResponse{
		ID:              inv.ID.String(),
		EventName:       inv.EventName,
		StartTime:       inv.StartTime.Format(time.RFC3339),
		EndTime:         inv.EndTime.Format(time.RFC3339),
		Status:          string(inv.Status),
		CuratorComment:  inv.CuratorComment,
		User:            &user,
		Department:      department,
		Organization:    organization,
		Equipment:       equipment,
		IssuedEquipment: issuedEquipment,
	}
}
