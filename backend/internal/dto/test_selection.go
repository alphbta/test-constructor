package dto

type TestSelectionResponse struct {
	EventID          uint       `json:"event_id"`
	SpecializationID uint       `json:"specialization_id"`
	Tests            []TestInfo `json:"tests"`
	AllCompleted     bool       `json:"all_completed"`
	EventPassed      bool       `json:"event_passed"`
}

type TestInfo struct {
	ConfigID       uint   `json:"config_id"`
	TestID         uint   `json:"test_id"`
	TestLink       string `json:"test_link"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	TimeLimit      int    `json:"time_limit"`
	IsExtra        bool   `json:"is_extra"`
	IsCommon       bool   `json:"is_common"`
	Status         string `json:"status"` // available, locked, in_progress, completed
	Score          int    `json:"score,omitempty"`
	MaxScore       int    `json:"max_score,omitempty"`
	Passed         bool   `json:"passed,omitempty"`
	AttemptID      uint   `json:"attempt_id,omitempty"`
	ReplacedTestID uint   `json:"replaced_test_id,omitempty"`
	ReplacedTitle  string `json:"replaced_title,omitempty"`
}
