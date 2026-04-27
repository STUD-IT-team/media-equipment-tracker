package dto

import (
	"time"

	"media-equipment-tracker/internal/application/studioservice"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateStudioRequest struct {
	EventName           string  `json:"event_name" binding:"required"`
	ShootingDescription string  `json:"shooting_description" binding:"required"`
	StartTime           string  `json:"start_time" binding:"required"`
	EndTime             string  `json:"end_time" binding:"required"`
	OrganizationID      *string `json:"organization_id" binding:"omitempty"`
	DepartmentID        *string `json:"department_id" binding:"omitempty"`
	NeedsChromakey      bool    `json:"needs_chromakey"`
	NeedsCyclorama      bool    `json:"needs_cyclorama"`
	NeedsBlackFabric    bool    `json:"needs_black_fabric"`
}

func DeserializeCreateStudioRequest(c *gin.Context) (*studioservice.CreateStudioInvocationRequest, error) {
	var req CreateStudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
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

	return &studioservice.CreateStudioInvocationRequest{
		EventName:           req.EventName,
		ShootingDescription: req.ShootingDescription,
		StartTime:           startTime,
		EndTime:             endTime,
		OrganizationID:      orgID,
		DepartmentID:        depID,
		NeedsChromakey:      req.NeedsChromakey,
		NeedsCyclorama:      req.NeedsCyclorama,
		NeedsBlackFabric:    req.NeedsBlackFabric,
	}, nil
}

type CreateStudioResponse struct {
	ShortStudioItem
}

func SerializeCreateStudioResponse(_ *gin.Context, inv *domain.StudioInvocation) any {
	return CreateStudioResponse{
		ShortStudioItem: ShortFromStudio(inv),
	}
}
