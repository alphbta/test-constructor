package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CRMClient interface {
	GetEvents() ([]Event, error)
	GetEventSpecializations(eventID int) (*CRMEventResponse, error)
	CreateTestSession(applicationID, testID uint, sessionID string, expiresAt time.Time) error
	SendTestResult(applicationID uint, result CRMResultData) error
}

type crmClient struct {
	baseURL string
	token   string
	client  *http.Client
}

type Event struct {
	ID              int              `json:"id"`
	Name            string           `json:"name"`
	StartDate       string           `json:"start_date"`
	EndDate         string           `json:"end_date"`
	Specializations []Specialization `json:"specializations"`
}

type EventDetailResponse struct {
	Specializations []Specialization `json:"specializations"`
}

type CRMEventResponse struct {
	Specializations []Specialization `json:"specializations"`
}

type Specialization struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CRMResultData struct {
	SessionID   string `json:"session_id"`
	TestID      string `json:"test_id"`
	Score       int    `json:"score"`
	MaxScore    int    `json:"max_score"`
	IsPassed    bool   `json:"is_passed"`
	CompletedAt string `json:"completed_at"`
	StartedAt   string `json:"started_at"`
}

type CRMCreateSessionData struct {
	TestID    uint   `json:"test_id"`
	SessionID string `json:"session_id"`
	ExpiresAt string `json:"expires_at"`
}

func NewCRMClient(baseURL, token string) CRMClient {
	return &crmClient{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{},
	}
}

func (c *crmClient) GetEvents() ([]Event, error) {
	url := c.baseURL + "/api/users/events/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("X-Service-Token", c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к CRM: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CRM вернул ошибку %d: %s", resp.StatusCode, string(body))
	}

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return events, nil
}

func (c *crmClient) GetEventSpecializations(eventID int) (*CRMEventResponse, error) {
	url := c.baseURL + fmt.Sprintf("/api/users/events/%d/", eventID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("X-Service-Token", c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к CRM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("мероприятие не найдено")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CRM вернул ошибку %d: %s", resp.StatusCode, string(body))
	}

	var eventData CRMEventResponse
	if err := json.NewDecoder(resp.Body).Decode(&eventData); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	return &eventData, nil
}

func (c *crmClient) CreateTestSession(applicationID, testID uint, sessionID string, expiresAt time.Time) error {
	url := c.baseURL + fmt.Sprintf("/api/users/integration/applications/%d/test-sessions/", applicationID)
	data := CRMCreateSessionData{
		TestID:    testID,
		SessionID: sessionID,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z"),
	}
	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Token", c.token)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("CRM вернул ошибку %d: %v", resp.StatusCode, errResp)
	}
	return nil
}

func (c *crmClient) SendTestResult(applicationID uint, result CRMResultData) error {
	url := c.baseURL + fmt.Sprintf("/api/users/integration/applications/%d/test-results/", applicationID)
	body, _ := json.Marshal(result)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Token", c.token)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки результатов: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("CRM вернул ошибку %d: %v", resp.StatusCode, errResp)
	}
	return nil
}
