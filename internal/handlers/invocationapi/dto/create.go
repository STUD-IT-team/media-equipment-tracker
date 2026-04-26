package dto

import (
	"time"

	"media-equipment-tracker/internal/application/invocationservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateInvocationRequest struct {
	EventName      string   `json:"event_name" binding:"required"`
	StartTime      string   `json:"start_time" binding:"required"`
	EndTime        string   `json:"end_time" binding:"required"`
	OrganizationID *string  `json:"organization_id" binding:"omitempty"`
	DepartmentID   *string  `json:"department_id" binding:"omitempty"`
	EquipmentIDs   []string `json:"equipment_ids" binding:"required,min=1,dive"`
}

func DeserializeCreateInvocationRequest(c *gin.Context) (*invocationservice.CreateInvocationRequest, error) {
	var req CreateInvocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	var orgID *uuid.UUID
	if req.OrganizationID != nil {
		id, err := uuid.Parse(*req.OrganizationID)
		if err != nil {
			return nil, err
		}
		orgID = &id
	}

	var depID *uuid.UUID
	if req.DepartmentID != nil {
		id, err := uuid.Parse(*req.DepartmentID)
		if err != nil {
			return nil, err
		}
		depID = &id
	}

	equipmentIDs := make([]uuid.UUID, 0, len(req.EquipmentIDs))
	for _, id := range req.EquipmentIDs {
		equipmentID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		equipmentIDs = append(equipmentIDs, equipmentID)
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	return &invocationservice.CreateInvocationRequest{
		EventName:      req.EventName,
		StartTime:      startTime,
		EndTime:        endTime,
		OrganizationID: orgID,
		DepartmentID:   depID,
		EquipmentIDs:   equipmentIDs,
	}, nil
}

type CreateInvocationResponse struct {
	ShortInvocationItem
}

func SerializeCreateInvocationResponse(_ *gin.Context, inv *domain.EquipmentInvocation) any {
	return CreateInvocationResponse{
		ShortInvocationItem: ShortFromInvocation(inv),
	}
}
