package http

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"uala-timeline-service/config"
	"uala-timeline-service/internal/application"
)

var (
	ErrInvalidUser = errors.New("invalid user")
)

func getUserTimelineByDay(deps *config.Dependencies) http.HandlerFunc {
	createPost := application.NewGetUserTimeline(deps.TimelineService)
	return func(w http.ResponseWriter, r *http.Request) {
		var cmd application.GetUserTimelineCommand
		err := json.NewDecoder(r.Body).Decode(&cmd)
		if err != nil {
			handleError(w, err)
			return
		}

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			handleError(w, ErrInvalidUser)
			return
		}
		cmd.UserID = userID
		response, err := createPost.Exec(r.Context(), &cmd)
		if err != nil {
			handleError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			handleError(w, err)
			return
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var errorResp ErrorResponse

	switch {
	case errors.Is(err, ErrInvalidUser):
		errorResp = ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid user",
			Code:       "BAD_REQUEST",
		}
	default:
		errorResp = ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
			Code:       "INTERNAL_ERROR",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorResp.StatusCode)

	jsonResp, jsonErr := json.Marshal(errorResp)
	if jsonErr != nil {
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
	return
}

type ErrorResponse struct {
	StatusCode int    `json:"status,omitempty"`
	Message    string `json:"message,omitempty"`
	Code       string `json:"code,omitempty"`
}
