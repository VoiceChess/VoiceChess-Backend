package models

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

type ChatResponse struct {
	Response string `json:"response,omitempty"`
	Screen   string `json:"screen,omitempty"`
	Error    string `json:"error,omitempty"`
}
