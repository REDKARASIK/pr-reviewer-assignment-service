package response

import (
	"encoding/json"
	"net/http"
)

type ErrBodyResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrBodyResponse `json:"error"`
}

func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, ErrorResponse{
		Error: ErrBodyResponse{
			Code:    code,
			Message: message,
		},
	})
}
