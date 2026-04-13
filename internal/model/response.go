package model

type MessageResponse struct {
	Message string `json:"message"`
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}
