package service

import (
	"errors"
	"test-constructor/internal/domain"
	"test-constructor/internal/dto"
	"test-constructor/internal/repository"

	"gorm.io/gorm"
)

type TestSelectionService interface {
	GetSelection(userID uint, eventID uint) (*dto.TestSelectionResponse, error)
}

type testSelectionService struct {
	eventConfigRepo    repository.EventConfigRepository
	userEventRepo      repository.UserEventRepository
	attemptRepo        repository.AttemptRepository
	extraThresholdRepo repository.ExtraThresholdRepository
}

func NewTestSelectionService(
	eventConfigRepo repository.EventConfigRepository,
	userEventRepo repository.UserEventRepository,
	attemptRepo repository.AttemptRepository,
	extraThresholdRepo repository.ExtraThresholdRepository,
) TestSelectionService {
	return &testSelectionService{
		eventConfigRepo:    eventConfigRepo,
		userEventRepo:      userEventRepo,
		attemptRepo:        attemptRepo,
		extraThresholdRepo: extraThresholdRepo,
	}
}

func (s *testSelectionService) GetSelection(userID uint, eventID uint) (*dto.TestSelectionResponse, error) {
	userEvent, err := s.userEventRepo.FindByUserAndEvent(userID, eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("вы не записаны на это мероприятие")
		}
		return nil, err
	}

	specID := userEvent.SpecializationID

	mainConfigs, err := s.eventConfigRepo.FindMainConfigsByEventAndSpec(eventID, specID)
	if err != nil {
		return nil, err
	}

	specificConfigs, _ := s.eventConfigRepo.FindByEventAndSpecialization(eventID, specID, false)
	commonConfigs, _ := s.eventConfigRepo.FindCommonConfigsByEventID(eventID)
	extraConfigs, _ := s.eventConfigRepo.FindByEventAndSpecialization(eventID, specID, true)

	allConfigs := append(specificConfigs, commonConfigs...)
	allConfigs = append(allConfigs, extraConfigs...)

	configIDs := make([]uint, len(allConfigs))
	for i, cfg := range allConfigs {
		configIDs[i] = cfg.ConfigID
	}
	attempts, _ := s.attemptRepo.FindByConfigIDsAndUser(configIDs, userID)
	attemptMap := make(map[uint]*domain.Attempt)
	for i := range attempts {
		attemptMap[attempts[i].ConfigID] = &attempts[i]
	}

	var tests []dto.TestInfo

	for _, mainCfg := range mainConfigs {
		mainAttempt, hasMainAttempt := attemptMap[mainCfg.ConfigID]
		mainPassed := hasMainAttempt && mainAttempt.EndTime != nil && mainAttempt.Passed

		isCommon := mainCfg.SpecializationID == 0

		if !hasMainAttempt {
			tests = append(tests, dto.TestInfo{
				ConfigID:    mainCfg.ConfigID,
				TestID:      mainCfg.TestID,
				TestLink:    mainCfg.TestLink.String(),
				Title:       mainCfg.Test.Title,
				Description: mainCfg.Test.Description,
				TimeLimit:   mainCfg.TimeLimit,
				IsExtra:     false,
				IsCommon:    isCommon,
				Status:      "available",
			})
		} else if mainAttempt.EndTime == nil {
			tests = append(tests, dto.TestInfo{
				ConfigID:    mainCfg.ConfigID,
				TestID:      mainCfg.TestID,
				TestLink:    mainCfg.TestLink.String(),
				Title:       mainCfg.Test.Title,
				Description: mainCfg.Test.Description,
				TimeLimit:   mainCfg.TimeLimit,
				IsExtra:     false,
				IsCommon:    isCommon,
				Status:      "in_progress",
				AttemptID:   mainAttempt.AttemptID,
			})
		} else if mainPassed {
			tests = append(tests, dto.TestInfo{
				ConfigID:    mainCfg.ConfigID,
				TestID:      mainCfg.TestID,
				TestLink:    mainCfg.TestLink.String(),
				Title:       mainCfg.Test.Title,
				Description: mainCfg.Test.Description,
				TimeLimit:   mainCfg.TimeLimit,
				IsExtra:     false,
				IsCommon:    isCommon,
				Status:      "completed",
				Score:       mainAttempt.Score,
				MaxScore:    mainAttempt.MaxScore,
				Passed:      true,
				AttemptID:   mainAttempt.AttemptID,
			})
		} else {
			replacements, _ := s.extraThresholdRepo.FindReplacementsForConfigID(mainCfg.ConfigID)

			tests = append(tests, dto.TestInfo{
				ConfigID:    mainCfg.ConfigID,
				TestID:      mainCfg.TestID,
				TestLink:    mainCfg.TestLink.String(),
				Title:       mainCfg.Test.Title,
				Description: mainCfg.Test.Description,
				TimeLimit:   mainCfg.TimeLimit,
				IsExtra:     false,
				IsCommon:    isCommon,
				Status:      "completed",
				Score:       mainAttempt.Score,
				MaxScore:    mainAttempt.MaxScore,
				Passed:      false,
				AttemptID:   mainAttempt.AttemptID,
			})

			for _, replacement := range replacements {
				extraCfg := replacement.ExtraConfig

				canAccess := mainAttempt.Score >= replacement.Threshold
				extraAttempt, hasExtraAttempt := attemptMap[extraCfg.ConfigID]

				status := "locked"
				if canAccess {
					if !hasExtraAttempt {
						status = "available"
					} else if extraAttempt.EndTime == nil {
						status = "in_progress"
					} else {
						status = "completed"
					}
				}

				testInfo := dto.TestInfo{
					ConfigID:       extraCfg.ConfigID,
					TestID:         extraCfg.TestID,
					TestLink:       extraCfg.TestLink.String(),
					Title:          extraCfg.Test.Title,
					Description:    extraCfg.Test.Description,
					TimeLimit:      extraCfg.TimeLimit,
					IsExtra:        true,
					IsCommon:       isCommon,
					ReplacedTestID: mainCfg.TestID,
					ReplacedTitle:  mainCfg.Test.Title,
					Status:         status,
				}

				if hasExtraAttempt && extraAttempt.EndTime != nil {
					testInfo.Score = extraAttempt.Score
					testInfo.MaxScore = extraAttempt.MaxScore
					testInfo.Passed = extraAttempt.Passed
					testInfo.AttemptID = extraAttempt.AttemptID
				} else if hasExtraAttempt && extraAttempt.EndTime == nil {
					testInfo.AttemptID = extraAttempt.AttemptID
				}

				tests = append(tests, testInfo)
			}
		}
	}

	allCompleted := s.checkAllTestsCompleted(userID, eventID, specID)
	eventPassed := s.checkEventPassed(userID, eventID, specID)

	return &dto.TestSelectionResponse{
		EventID:          eventID,
		SpecializationID: specID,
		Tests:            tests,
		AllCompleted:     allCompleted,
		EventPassed:      eventPassed,
	}, nil
}

func (s *testSelectionService) checkAllTestsCompleted(userID, eventID, specializationID uint) bool {
	mainConfigs, err := s.eventConfigRepo.FindMainConfigsByEventAndSpec(eventID, specializationID)
	if err != nil {
		return false
	}

	for _, mainCfg := range mainConfigs {
		mainAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, mainCfg.ConfigID)
		if err != nil || mainAttempt.EndTime == nil {
			return false
		}

		if mainAttempt.Passed {
			continue
		}

		if !s.hasCompletedReplacement(userID, mainCfg.ConfigID) {
			return false
		}
	}

	return true
}

func (s *testSelectionService) hasPassedAnyReplacement(userID, mainConfigID uint) bool {
	replacements, err := s.extraThresholdRepo.FindReplacementsForConfigID(mainConfigID)
	if err != nil {
		return false
	}

	mainAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, mainConfigID)
	if err != nil || mainAttempt.EndTime == nil {
		return false
	}

	for _, replacement := range replacements {
		if mainAttempt.Score < replacement.Threshold {
			continue
		}

		extraAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, replacement.ExtraConfigID)
		if err == nil && extraAttempt.EndTime != nil && extraAttempt.Passed {
			return true
		}
	}
	return false
}

func (s *testSelectionService) hasCompletedReplacement(userID, mainConfigID uint) bool {
	replacements, err := s.extraThresholdRepo.FindReplacementsForConfigID(mainConfigID)
	if err != nil {
		return false
	}

	mainAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, mainConfigID)
	if err != nil || mainAttempt.EndTime == nil {
		return false
	}

	for _, replacement := range replacements {
		if mainAttempt.Score < replacement.Threshold {
			continue
		}

		extraAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, replacement.ExtraConfigID)
		if err == nil && extraAttempt.EndTime != nil {
			return true
		}
	}
	return false
}

func (s *testSelectionService) checkEventPassed(userID, eventID, specializationID uint) bool {
	mainConfigs, err := s.eventConfigRepo.FindMainConfigsByEventAndSpec(eventID, specializationID)
	if err != nil || len(mainConfigs) == 0 {
		return false
	}

	for _, mainCfg := range mainConfigs {
		mainAttempt, err := s.attemptRepo.FindByUserAndConfig(userID, mainCfg.ConfigID)
		if err != nil || mainAttempt.EndTime == nil {
			return false
		}

		if mainAttempt.Passed {
			continue
		}

		if !s.hasPassedAnyReplacement(userID, mainCfg.ConfigID) {
			return false
		}
	}
	return true
}
