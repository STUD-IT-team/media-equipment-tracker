package dto

import (
	"fmt"
	"time"

	"media-equipment-tracker/internal/application/studioservice/studiosearch"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeserializeSearchStudioRequest(c *gin.Context) (*studiosearch.SearchStudioInvocationRequest, error) {
	var searchString *string
	if val := c.Query("search"); val != "" {
		searchString = &val
	}

	var status *domain.StudioInvocationStatus
	if val := c.Query("status"); val != "" {
		s, err := mapStudioStatus(val)
		if err != nil {
			return nil, err
		}
		status = &s
	}

	var startTime *time.Time
	if val := c.Query("start_time"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		startTime = &t
	}

	var endTime *time.Time
	if val := c.Query("end_time"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		endTime = &t
	}

	var adminID *uuid.UUID
	if val := c.Query("admin_id"); val != "" {
		id, err := uuid.Parse(val)
		if err != nil {
			return nil, err
		}
		adminID = &id
	}

	var userID *uuid.UUID
	if val := c.Query("user_id"); val != "" {
		id, err := uuid.Parse(val)
		if err != nil {
			return nil, err
		}
		userID = &id
	}

	var departmentID *uuid.UUID
	if val := c.Query("department_id"); val != "" {
		id, err := uuid.Parse(val)
		if err != nil {
			return nil, err
		}
		departmentID = &id
	}

	var organizationID *uuid.UUID
	if val := c.Query("organization_id"); val != "" {
		id, err := uuid.Parse(val)
		if err != nil {
			return nil, err
		}
		organizationID = &id
	}

	return &studiosearch.SearchStudioInvocationRequest{
		SearchString:   searchString,
		Statuses:       []domain.StudioInvocationStatus{*status},
		StartTime:      startTime,
		EndTime:        endTime,
		AdminID:        adminID,
		UserID:         userID,
		DepartmentID:   departmentID,
		OrganizationID: organizationID,
	}, nil
}

func mapStudioStatus(status string) (domain.StudioInvocationStatus, error) {
	switch status {
	case "created":
		return domain.StudioCreated, nil
	case "under_review":
		return domain.StudioUnderReview, nil
	case "changes_required":
		return domain.StudioChangesRequired, nil
	case "approved":
		return domain.StudioApproved, nil
	case "completed":
		return domain.StudioCompleted, nil
	case "cancelled":
		return domain.StudioCancelled, nil
	default:
		return "", fmt.Errorf("incorrect status")
	}
}

type SearchStudioResponse struct {
	Items []ShortStudioItem `json:"items"`
}

func SerializeSearchStudioResponse(_ *gin.Context, items []*domain.StudioInvocation) any {
	response := SearchStudioResponse{
		Items: make([]ShortStudioItem, 0, len(items)),
	}

	for _, item := range items {
		response.Items = append(response.Items, ShortFromStudio(item))
	}

	return response
}
