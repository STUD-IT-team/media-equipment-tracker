package dto

import (
	"fmt"
	"time"

	"media-equipment-tracker/internal/application/invocationservice/invocationsearch"
	"media-equipment-tracker/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func mapInvocationStatus(status string) (domain.EquipmentInvocationStatus, error) {
	switch status {
	case "created":
		return domain.InvocationCreated, nil
	case "under_review":
		return domain.InvocationUnderReview, nil
	case "changes_required":
		return domain.InvocationChangesRequired, nil
	case "approved":
		return domain.InvocationApproved, nil
	case "equipment_issued":
		return domain.InvocationEquipmentIssued, nil
	case "equipment_returned":
		return domain.InvocationEquipmentReturned, nil
	case "completed":
		return domain.InvocationCompleted, nil
	case "cancelled":
		return domain.InvocationCancelled, nil
	default:
		return "", fmt.Errorf("invalid status: %s", status)
	}
}

func DeserializeSearchInvocationRequest(c *gin.Context) (invocationsearch.SearchInvocationRequest, error) {
	var req invocationsearch.SearchInvocationRequest

	if val := c.Query("status"); val != "" {
		status, err := mapInvocationStatus(val)
		if err != nil {
			return req, err
		}
		req.Statuses = []domain.EquipmentInvocationStatus{status}
	}

	if val := c.Query("department_id"); val != "" {
		id, err := uuid.Parse(val)
		if err != nil {
			return req, err
		}
		req.DepartmentID = &id
	}

	if val := c.Query("organization_id"); val != "" {
		id, err := uuid.Parse(val)
		if err != nil {
			return req, err
		}
		req.OrganizationID = &id
	}

	if val := c.Query("user_id"); val != "" {
		id, err := uuid.Parse(val)
		if err != nil {
			return req, err
		}
		req.UserID = &id
	}

	if val := c.Query("start_time"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return req, err
		}
		req.StartTime = &t
	}

	if val := c.Query("end_time"); val != "" {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return req, err
		}
		req.EndTime = &t
	}

	if val := c.Query("search"); val != "" {
		req.SearchString = &val
	}

	return req, nil
}

type SearchInvocationResponse struct {
	Items []ShortInvocationItem `json:"items"`
}

func SerializeSearchInvocationResponse(_ *gin.Context, items []*domain.EquipmentInvocation) any {
	response := SearchInvocationResponse{
		Items: make([]ShortInvocationItem, 0, len(items)),
	}

	for _, item := range items {
		response.Items = append(response.Items, ShortFromInvocation(item))
	}

	return response
}
