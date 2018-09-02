package cmd

import (
	"net/http"

	"github.com/go-chi/chi"
)

type server struct {
	Router *chi.Router
	DB     *db.Sql
	Server *http.Server
}
