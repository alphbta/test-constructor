package dto

type EventResponse struct {
	Name            string                        `json:"name" example:"Хакатон 2024"`
	StartDate       string                        `json:"start_date" example:"2024-12-01T00:00:00Z"`
	EndDate         string                        `json:"end_date" example:"2024-12-02T00:00:00Z"`
	Specializations []EventSpecializationResponse `json:"specializations"`
	TotalTests      int                           `json:"total_tests,omitempty" example:"5"`
}

type EventsListResponse struct {
	Events []EventResponse `json:"events"`
	Total  int             `json:"total" example:"10"`
}

type EventSpecializationResponse struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Frontend Developer"`
}

type EventSpecializationsListResponse struct {
	EventID         int                           `json:"event_id" example:"1"`
	Specializations []EventSpecializationResponse `json:"specializations"`
}
