package service

import (
	"errors"
	"test-constructor/internal/domain"
	"test-constructor/internal/dto"
	"test-constructor/internal/repository"

	"gorm.io/gorm"
)

type UserEventService interface {
	CreateUserEvent(userID uint, req dto.CreateUserEventRequest) error
	GetUserEvents(userID uint) (*dto.UserEventsListResponse, error)
}

type userEventService struct {
	userEventRepo repository.UserEventRepository
}

func NewUserEventService(userEventRepo repository.UserEventRepository) UserEventService {
	return &userEventService{
		userEventRepo: userEventRepo,
	}
}

func (s *userEventService) CreateUserEvent(userID uint, req dto.CreateUserEventRequest) error {
	if req.EventID < 1 {
		return errors.New("event ID должен быть положительным")
	}

	if req.ApplicationID < 1 {
		return errors.New("application ID должен быть положительным")
	}

	existing, err := s.userEventRepo.FindByUserAndEvent(userID, req.EventID)
	if err == nil && existing != nil {
		return errors.New("вы уже записаны на это мероприятие")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("ошибка базы данных")
	}

	userEvent := domain.UserEvent{
		UserID:        userID,
		EventID:       req.EventID,
		ApplicationID: req.ApplicationID,
	}

	if err := s.userEventRepo.Create(&userEvent); err != nil {
		return errors.New("ошибка создания связи")
	}

	return nil
}

func (s *userEventService) GetUserEvents(userID uint) (*dto.UserEventsListResponse, error) {
	userEvents, err := s.userEventRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("ошибка получения данных")
	}

	response := &dto.UserEventsListResponse{
		Events: make([]dto.UserEventResponse, len(userEvents)),
	}

	for i, ue := range userEvents {
		response.Events[i] = dto.UserEventResponse{
			ID:            ue.ID,
			EventID:       ue.EventID,
			UserID:        ue.UserID,
			ApplicationID: ue.ApplicationID,
		}
	}

	return response, nil
}
