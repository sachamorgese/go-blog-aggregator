package config

import (
	"BlogAggregator/internal/database"
	"BlogAggregator/server-actions"
	"database/sql"
	"fmt"
	"github.com/go-chi/cors"
	"net/http"
)

var CorsOptions = cors.Options{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
	AllowedHeaders: []string{"*"},
}

type ApiConfig struct {
	Db *sql.DB
}

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) MiddlewareAuth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")

		if apiKey == "" || apiKey[0:6] != "ApiKey" {
			err := fmt.Errorf("Missing or invalid API key")

			fmt.Println(err.Error())

			server_actions.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		dbQueries := database.New(cfg.Db)

		user, err := dbQueries.GetUserByApiKey(r.Context(), apiKey[7:])

		if err != nil {
			server_actions.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}

		handler(w, r, user)
	}
}
