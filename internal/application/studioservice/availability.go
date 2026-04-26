package studioservice

import (
	"context"
	"time"

	"media-equipment-tracker/internal/application/studioservice/studiosearch"
	"media-equipment-tracker/internal/domain"
	"media-equipment-tracker/internal/domain/errs"
	"media-equipment-tracker/internal/utils/validate"

	"github.com/google/uuid"
)

type AvailabilityStudioInvocationRequest struct {
	StartTime time.Time `validate:"required"`
	EndTime   time.Time `validate:"gtfield=StartTime"`

	ForInvocation *uuid.UUID
}

type AvailabilityStudioService interface {
	Availability(ctx context.Context, req AvailabilityStudioInvocationRequest) (*StudioInvocationAvailabilityResponse, error)
}

type StudioInvocationAvailabilityResponse struct {
	Available              bool
	ConflictingInvocations []*domain.StudioInvocation
}

type studioInvocationAvailabilityService struct {
	search studiosearch.SearchStudioService
}

func NewStudioInvocationAvailabilityService(
	search studiosearch.SearchStudioService,
) AvailabilityStudioService {
	return &studioInvocationAvailabilityService{
		search: search,
	}
}

var _ AvailabilityStudioService = (*studioInvocationAvailabilityService)(nil)

func (s *studioInvocationAvailabilityService) Availability(ctx context.Context, req AvailabilityStudioInvocationRequest) (*StudioInvocationAvailabilityResponse, error) {
	if err := validate.ValidateStruct(req); err != nil {
		return nil, errs.NewValidationError("AvailabilityStudioInvocationRequest", err.Error())
	}

	invs, err := s.search.Search(ctx, &studiosearch.SearchStudioInvocationRequest{
		StartTime: &req.StartTime,
		EndTime:   &req.EndTime,
	})
	if err != nil {
		return nil, err
	}

	if req.ForInvocation != nil {
		for i, inv := range invs {
			if inv.ID == *req.ForInvocation {
				invs = append(invs[:i], invs[i+1:]...)
				break
			}
		}
	}

	if len(invs) == 0 {
		return &StudioInvocationAvailabilityResponse{
			Available: true,
		}, nil
	}
	return &StudioInvocationAvailabilityResponse{
		Available:              false,
		ConflictingInvocations: invs,
	}, nil
}
