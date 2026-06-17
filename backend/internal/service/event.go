package service

import (
	"test-constructor/internal/client"
	"test-constructor/internal/dto"
)

type EventService interface {
	GetEvents() (*dto.EventsListResponse, error)
	GetEventSpecializations(eventID int) (*dto.EventSpecializationsListResponse, error)
}

type eventService struct {
	crmClient client.CRMClient
}

func NewEventService(crmClient client.CRMClient) EventService {
	return &eventService{
		crmClient: crmClient,
	}
}

func (s *eventService) GetEvents() (*dto.EventsListResponse, error) {
	events, err := s.crmClient.GetEvents()
	if err != nil {
		return nil, err
	}

	response := &dto.EventsListResponse{
		Events: make([]dto.EventResponse, len(events)),
		Total:  len(events),
	}

	for i, event := range events {
		specializations := make([]dto.EventSpecializationResponse, len(event.Specializations))
		for j, spec := range event.Specializations {
			specializations[j] = dto.EventSpecializationResponse{
				ID:   spec.ID,
				Name: spec.Name,
			}
		}

		response.Events[i] = dto.EventResponse{
			Name:            event.Name,
			StartDate:       event.StartDate,
			EndDate:         event.EndDate,
			Specializations: specializations,
		}
	}

	return response, nil
}

func (s *eventService) GetEventSpecializations(eventID int) (*dto.EventSpecializationsListResponse, error) {
	eventDetail, err := s.crmClient.GetEventSpecializations(eventID)
	if err != nil {
		return nil, err
	}

	specializations := make([]dto.EventSpecializationResponse, len(eventDetail.Specializations))
	for i, spec := range eventDetail.Specializations {
		specializations[i] = dto.EventSpecializationResponse{
			ID:   spec.ID,
			Name: spec.Name,
		}
	}

	return &dto.EventSpecializationsListResponse{
		EventID:         eventID,
		Specializations: specializations,
	}, nil
}
