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

func addPostToUserTimeline(deps *config.Dependencies) http.HandlerFunc {
	createPost := application.NewAddPostToUserTimeline(deps.TimelineService)
	return func(w http.ResponseWriter, r *http.Request) {
		var cmd application.AddPostToUserTimelineCommand
		err := json.NewDecoder(r.Body).Decode(&cmd)
		if err != nil {
			handleError(w, err)
			return
		}
		err = createPost.Exec(r.Context(), &cmd)
		if err != nil {
			handleError(w, err)
			return
		}
	}
}

func getUserTimeline(deps *config.Dependencies) http.HandlerFunc {
	createPost := application.NewGetUserTimeline(deps.TimelineService)
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			handleError(w, ErrInvalidUser)
			return
		}
		cmd := application.GetUserTimelineCommand{
			UserID: userID,
		}
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
