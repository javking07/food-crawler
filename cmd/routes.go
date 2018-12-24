package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type HealthResponse struct {
	Status string
}

func (a *App) Health(w http.ResponseWriter, r *http.Request) {
	var Health struct {
		Status         string `json:"status"`
		DatabaseStatus string `json:"database_status"`
	}
	Health.Status = "ok"

	if a.AppDatabaseCQL != nil {
		err := a.AppDatabaseCQL.Query("SELECT now() FROM system.local").Exec()

		if err != nil {
			Health.DatabaseStatus = "not ok"
		} else {
			Health.DatabaseStatus = "ok"
		}

	}
	respondWithJSON(w, http.StatusOK, Health)

}

// todo create middleware for capturing trace data on request
func (a *App) WithTracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

//ValidateParams ensures the proper headers are in the request before passing request to key lookup functions
func (a *App) ValidateParams(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// extract params from url
		ctx := r.Context()
		id := chi.URLParam(r, "id")
		if id == "" {
			a.AppLogger.Warn().Msg("no id present in request")
			respondWithError(w, http.StatusBadRequest, "no id present in request")
			return
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			a.AppLogger.Warn().Msgf("id %v has invalid format", idInt)
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("id %v has invalid format", id))
			return
		}

		// add url params to context
		ctx = context.WithValue(ctx, ContextKeyID, idInt)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
