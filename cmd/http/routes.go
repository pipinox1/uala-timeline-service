package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ddchi "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"
	"net/http"
	"uala-timeline-service/config"
)

func SetupRouterAndRoutes(config *config.Config, deps *config.Dependencies) chi.Router {
	router := chi.NewRouter()

	router.Route("/health", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	router.Route("/api/v1/user_timeline", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Use(ddchi.Middleware(ddchi.WithServiceName(config.ServiceName)))
		r.Get("/{user_id}", getUserTimeline(deps))
		r.Post("/add", addPostToUserTimeline(deps))
	})

	return router
}
