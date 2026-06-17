package service

import (
	"errors"
	"fmt"
	"log"
	"test-constructor/internal/client"
	"test-constructor/internal/domain"
	"test-constructor/internal/dto"
	"test-constructor/internal/repository"
)

type EventConfigService interface {
	CreateConfig(creatorID uint, req dto.CreateEventConfigRequest) (*dto.CreateEventConfigResponse, error)
	UpdateConfig(configID, creatorID uint, req dto.UpdateEventConfigRequest) (*dto.UpdateEventConfigResponse, error)
	GetConfig(id uint) (*dto.EventConfigResponse, error)
}

type eventConfigService struct {
	questionRepo       repository.QuestionRepository
	testRepo           repository.TestRepository
	eventConfigRepo    repository.EventConfigRepository
	extraThresholdRepo repository.ExtraThresholdRepository
	txManager          repository.TransactionManager
	crmClient          client.CRMClient
	validationService  ValidationService
}

func NewEventConfigService(
	questionRepo repository.QuestionRepository,
	testRepo repository.TestRepository,
	eventConfigRepo repository.EventConfigRepository,
	extraThresholdRepo repository.ExtraThresholdRepository,
	txManager repository.TransactionManager,
	crmClient client.CRMClient,
	validationService ValidationService,
) EventConfigService {
	return &eventConfigService{
		questionRepo:       questionRepo,
		testRepo:           testRepo,
		eventConfigRepo:    eventConfigRepo,
		extraThresholdRepo: extraThresholdRepo,
		txManager:          txManager,
		crmClient:          crmClient,
		validationService:  validationService,
	}
}

func (s *eventConfigService) validateBaseRequest(eventID, testID uint, threshold int) error {
	if eventID < 1 {
		return errors.New("event ID должен быть положительным")
	}
	if testID < 1 {
		return errors.New("test ID должен быть положительным")
	}
	if threshold < 1 {
		return errors.New("пороговое значение должно быть положительным")
	}
	return nil
}

func (s *eventConfigService) CreateConfig(creatorID uint, req dto.CreateEventConfigRequest) (*dto.CreateEventConfigResponse, error) {
	if err := s.validateBaseRequest(req.EventID, req.TestID, req.Threshold); err != nil {
		return nil, err
	}

	_, err := s.testRepo.FindByID(req.TestID)
	if err != nil {
		return nil, fmt.Errorf("тест с ID %d не найден", req.TestID)
	}

	if err := s.validationService.ValidateThreshold(req.TestID, req.Threshold); err != nil {
		return nil, err
	}

	for i, extra := range req.ExtraThreshold {
		if _, err := s.testRepo.FindByID(extra.TestID); err != nil {
			return nil, fmt.Errorf("дополнительный тест #%d с ID %d не найден", i+1, extra.TestID)
		}

		if err := s.validationService.ValidateThreshold(extra.TestID, extra.TestThreshold); err != nil {
			return nil, fmt.Errorf("дополнительный тест #%d: %w", i+1, err)
		}
	}

	tx, err := s.txManager.Begin()
	if err != nil {
		return nil, errors.New("ошибка базы данных")
	}

	eventConfig := domain.EventConfig{
		EventID:          req.EventID,
		SpecializationID: req.SpecializationID,
		TestID:           req.TestID,
		CreatorID:        creatorID,
		SuccessText:      req.SuccessText,
		FailText:         req.FailText,
		TimeLimit:        req.TimeLimit,
		Threshold:        req.Threshold,
	}

	if err := s.eventConfigRepo.CreateWithTx(tx, &eventConfig); err != nil {
		s.txManager.Rollback(tx)
		return nil, fmt.Errorf("ошибка создания конфигурации: %w", err)
	}

	for i, eThreshold := range req.ExtraThreshold {
		mainMaxScore, _ := s.questionRepo.GetMaxScoreByTestID(req.TestID)
		if eThreshold.Threshold > mainMaxScore {
			s.txManager.Rollback(tx)
			return nil, fmt.Errorf(
				"порог перехода (%d) дополнительного теста #%d не может быть выше "+
					"максимального балла основного теста (%d)",
				eThreshold.Threshold, i+1, mainMaxScore,
			)
		}

		extraConfig := domain.EventConfig{
			EventID:          req.EventID,
			SpecializationID: req.SpecializationID,
			TestID:           eThreshold.TestID,
			CreatorID:        creatorID,
			SuccessText:      req.SuccessText,
			FailText:         req.FailText,
			TimeLimit:        req.TimeLimit,
			Threshold:        eThreshold.TestThreshold,
			IsExtra:          true,
		}

		if err := s.eventConfigRepo.CreateWithTx(tx, &extraConfig); err != nil {
			s.txManager.Rollback(tx)
			return nil, fmt.Errorf("ошибка создания дополнительного теста: %w", err)
		}

		extraThreshold := domain.ExtraThreshold{
			ConfigID:      eventConfig.ConfigID,
			Threshold:     eThreshold.Threshold,
			Message:       eThreshold.Message,
			ExtraConfigID: extraConfig.ConfigID,
		}

		if err := s.extraThresholdRepo.CreateWithTx(tx, &extraThreshold); err != nil {
			s.txManager.Rollback(tx)
			return nil, fmt.Errorf("ошибка создания дополнительного порога: %w", err)
		}
	}

	if err := s.txManager.Commit(tx); err != nil {
		return nil, errors.New("ошибка сохранения изменений")
	}

	return &dto.CreateEventConfigResponse{
		ConfigID: eventConfig.ConfigID,
		Message:  "Конфигурация создана",
	}, nil
}

func (s *eventConfigService) UpdateConfig(configID, creatorID uint, req dto.UpdateEventConfigRequest) (*dto.UpdateEventConfigResponse, error) {
	if req.EventID < 1 || req.SpecializationID < 1 || req.TestID < 1 {
		return nil, errors.New("ID должен быть положительным")
	}

	if req.Threshold < 1 {
		return nil, errors.New("пороговое значение должно быть положительным")
	}

	if err := s.validationService.ValidateThreshold(req.TestID, req.Threshold); err != nil {
		return nil, err
	}

	for i, extra := range req.ExtraThreshold {
		if err := s.validationService.ValidateThreshold(extra.TestID, extra.TestThreshold); err != nil {
			return nil, fmt.Errorf("дополнительный тест #%d: %w", i+1, err)
		}
	}

	tx, err := s.txManager.Begin()
	if err != nil {
		return nil, errors.New("ошибка базы данных")
	}

	existingConfig, err := s.eventConfigRepo.FindByID(configID)
	if err != nil {
		s.txManager.Rollback(tx)
		return nil, errors.New("конфигурация не найдена")
	}

	existingConfig.EventID = req.EventID
	existingConfig.SpecializationID = req.SpecializationID
	existingConfig.TestID = req.TestID
	existingConfig.SuccessText = req.SuccessText
	existingConfig.FailText = req.FailText
	existingConfig.TimeLimit = req.TimeLimit
	existingConfig.Threshold = req.Threshold

	if err := s.eventConfigRepo.UpdateWithTx(tx, existingConfig); err != nil {
		s.txManager.Rollback(tx)
		return nil, fmt.Errorf("ошибка обновления конфигурации: %w", err)
	}

	existingThresholds, err := s.extraThresholdRepo.FindByConfigID(configID)
	if err != nil {
		s.txManager.Rollback(tx)
		return nil, fmt.Errorf("ошибка получения дополнительных порогов: %w", err)
	}

	existingExtraConfigMap := make(map[uint]*domain.EventConfig)
	for _, et := range existingThresholds {
		if et.ExtraConfigID > 0 {
			config := et.ExtraConfig
			existingExtraConfigMap[et.ExtraConfigID] = &config
		}
	}

	if err := s.extraThresholdRepo.DeleteByConfigIDWithTx(tx, configID); err != nil {
		s.txManager.Rollback(tx)
		return nil, fmt.Errorf("ошибка удаления старых порогов: %w", err)
	}

	for i, eThreshold := range req.ExtraThreshold {
		mainMaxScore, _ := s.questionRepo.GetMaxScoreByTestID(req.TestID)
		if eThreshold.Threshold > mainMaxScore {
			s.txManager.Rollback(tx)
			return nil, fmt.Errorf(
				"порог перехода (%d) дополнительного теста #%d не может быть выше "+
					"максимального балла основного теста (%d)",
				eThreshold.Threshold, i+1, mainMaxScore,
			)
		}

		existingExtraConfig, exists := s.findExistingExtraConfig(existingExtraConfigMap, eThreshold.TestID)

		if exists {
			existingExtraConfig.EventID = req.EventID
			existingExtraConfig.SpecializationID = req.SpecializationID
			existingExtraConfig.TestID = eThreshold.TestID
			existingExtraConfig.SuccessText = req.SuccessText
			existingExtraConfig.FailText = req.FailText
			existingExtraConfig.TimeLimit = req.TimeLimit
			existingExtraConfig.Threshold = eThreshold.TestThreshold

			if err := s.eventConfigRepo.UpdateWithTx(tx, existingExtraConfig); err != nil {
				s.txManager.Rollback(tx)
				return nil, fmt.Errorf("ошибка обновления дополнительного теста: %w", err)
			}

			delete(existingExtraConfigMap, existingExtraConfig.ConfigID)
		} else {
			newExtraConfig := domain.EventConfig{
				EventID:          req.EventID,
				SpecializationID: req.SpecializationID,
				TestID:           eThreshold.TestID,
				CreatorID:        creatorID,
				SuccessText:      req.SuccessText,
				FailText:         req.FailText,
				TimeLimit:        req.TimeLimit,
				Threshold:        eThreshold.TestThreshold,
				IsExtra:          true,
			}

			if err := s.eventConfigRepo.CreateWithTx(tx, &newExtraConfig); err != nil {
				s.txManager.Rollback(tx)
				return nil, fmt.Errorf("ошибка создания дополнительного теста: %w", err)
			}

			existingExtraConfig = &newExtraConfig
		}

		extraThreshold := domain.ExtraThreshold{
			ConfigID:      configID,
			Threshold:     eThreshold.Threshold,
			Message:       eThreshold.Message,
			ExtraConfigID: existingExtraConfig.ConfigID,
		}

		if err := s.extraThresholdRepo.CreateWithTx(tx, &extraThreshold); err != nil {
			s.txManager.Rollback(tx)
			return nil, fmt.Errorf("ошибка создания дополнительного порога: %w", err)
		}
	}

	for _, unusedConfig := range existingExtraConfigMap {
		if err := s.eventConfigRepo.Delete(unusedConfig.ConfigID); err != nil {
			log.Printf("ошибка удаления неиспользуемых дополнительных конфигураций: %s", err.Error())
		}
	}

	if err := s.txManager.Commit(tx); err != nil {
		return nil, errors.New("ошибка сохранения изменений")
	}

	return &dto.UpdateEventConfigResponse{
		ConfigID: configID,
		Message:  "Конфигурация обновлена",
	}, nil
}

func (s *eventConfigService) GetConfig(id uint) (*dto.EventConfigResponse, error) {
	config, err := s.eventConfigRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("конфигурация не найдена")
	}

	response := &dto.EventConfigResponse{
		ConfigID:         config.ConfigID,
		EventID:          config.EventID,
		SpecializationID: config.SpecializationID,
		TestID:           config.TestID,
		SuccessText:      config.SuccessText,
		FailText:         config.FailText,
		TimeLimit:        config.TimeLimit,
		Threshold:        config.Threshold,
		TestLink:         config.TestLink.String(),
		IsExtra:          config.IsExtra,
	}

	for _, et := range config.ExtraThreshold {
		response.ExtraThreshold = append(response.ExtraThreshold, dto.ExtraThresholdResponse{
			Threshold: et.Threshold,
			Message:   et.Message,
			TestID:    et.ExtraConfig.TestID,
		})
	}

	return response, nil
}

func (s *eventConfigService) findExistingExtraConfig(
	existingConfigs map[uint]*domain.EventConfig,
	testID uint,
) (*domain.EventConfig, bool) {
	for _, config := range existingConfigs {
		if config.TestID == testID {
			return config, true
		}
	}
	return nil, false
}
