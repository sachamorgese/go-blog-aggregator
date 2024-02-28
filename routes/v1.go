package routes

import (
	"BlogAggregator/config"
	"BlogAggregator/internal/database"
	"BlogAggregator/server-actions"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status: "ok",
	}

	server_actions.RespondWithJSON(w, http.StatusOK, response)
}

func SendError(w http.ResponseWriter, r *http.Request) {
	server_actions.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

type UserBody struct {
	Name string `json:"name"`
}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		params := UserBody{}
		err := decoder.Decode(&params)

		if err != nil {
			server_actions.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		dbQueries := database.New(db)

		user, err := dbQueries.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      params.Name,
		})

		if err != nil {
			server_actions.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
			return
		}

		server_actions.RespondWithJSON(w, http.StatusCreated, user)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request, users database.User) {
	server_actions.RespondWithJSON(w, http.StatusOK, users)
}

type FeedBody struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func CreateFeed(cfg *config.ApiConfig) config.AuthedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		decoder := json.NewDecoder(r.Body)
		params := FeedBody{}
		err := decoder.Decode(&params)

		if err != nil {
			server_actions.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		dbQueries := database.New(cfg.Db)
		feed, err := dbQueries.CreateFeed(r.Context(), database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      params.Name,
			Url:       params.Url,
			UserID:    user.ID,
		})

		if err != nil {
			fmt.Println(err.Error())
			server_actions.RespondWithError(w, http.StatusInternalServerError, "Error creating feed")
			return
		}

		feedIDNullUUID := uuid.NullUUID{UUID: feed.ID, Valid: true}
		userIDNullUUID := uuid.NullUUID{UUID: user.ID, Valid: true}

		feedFollow, err := dbQueries.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID:    feedIDNullUUID,
			UserID:    userIDNullUUID,
		})

		if err != nil {
			server_actions.RespondWithError(w, http.StatusInternalServerError, "Error creating feed follow")
			return
		}

		server_actions.RespondWithJSON(w, http.StatusCreated, feedFollow)
	}
}

type FeedResponse struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Url           string    `json:"url"`
	UserID        uuid.UUID `json:"user_id"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
}

func getFeeds(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbQueries := database.New(cfg.Db)
		feeds, err := dbQueries.GetFeeds(r.Context())

		if err != nil {
			server_actions.RespondWithError(w, http.StatusInternalServerError, "Error getting feeds")
			return
		}

		new_feeds := make([]FeedResponse, len(feeds))

		for i, feed := range feeds {
			new_feeds[i] = FeedResponse{
				ID:            feed.ID,
				CreatedAt:     feed.CreatedAt,
				UpdatedAt:     feed.UpdatedAt,
				Name:          feed.Name,
				Url:           feed.Url,
				UserID:        feed.UserID,
				LastFetchedAt: feed.LastFetchedAt.Time,
			}
		}

		server_actions.RespondWithJSON(w, http.StatusOK, new_feeds)
	}
}

type FeedFollowBody struct {
	FeedID uuid.UUID `json:"feed_id"`
}

func CreateFeedFollow(cfg *config.ApiConfig) config.AuthedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		decoder := json.NewDecoder(r.Body)
		params := FeedFollowBody{}
		err := decoder.Decode(&params)

		if err != nil || params.FeedID == uuid.Nil {
			server_actions.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		feedIDNullUUID := uuid.NullUUID{UUID: params.FeedID, Valid: true}
		userIDNullUUID := uuid.NullUUID{UUID: user.ID, Valid: true}

		dbQueries := database.New(cfg.Db)
		feedFollow, err := dbQueries.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID:    feedIDNullUUID,
			UserID:    userIDNullUUID,
		})

		if err != nil {
			server_actions.RespondWithError(w, http.StatusInternalServerError, "Error creating feed follow")
			return
		}

		server_actions.RespondWithJSON(w, http.StatusCreated, feedFollow)
	}
}

func DeleteFeedFollow(cfg *config.ApiConfig) config.AuthedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		decoder := json.NewDecoder(r.Body)
		params := FeedFollowBody{}
		err := decoder.Decode(&params)

		if err != nil || params.FeedID == uuid.Nil {
			server_actions.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		feedIDNullUUID := uuid.NullUUID{UUID: params.FeedID, Valid: true}
		userIDNullUUID := uuid.NullUUID{UUID: user.ID, Valid: true}

		dbQueries := database.New(cfg.Db)

		err = dbQueries.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
			FeedID: feedIDNullUUID,
			UserID: userIDNullUUID,
		})

		if err != nil {
			server_actions.RespondWithError(w, http.StatusInternalServerError, "Error deleting feed follow")
			return
		}

		server_actions.RespondWithJSON(w, http.StatusOK, nil)
	}
}

func GetAllFeedFollows(cfg *config.ApiConfig) config.AuthedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {

		dbQueries := database.New(cfg.Db)

		nullUserId := uuid.NullUUID{UUID: user.ID, Valid: true}

		feedFollows, err := dbQueries.GetAllFeedFollowsForUser(r.Context(), nullUserId)

		if err != nil {
			server_actions.RespondWithError(w, http.StatusInternalServerError, "Error getting feed follows")
			return
		}

		server_actions.RespondWithJSON(w, http.StatusOK, feedFollows)
	}
}

func GetPosts(cfg *config.ApiConfig) config.AuthedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {

		dbQueries := database.New(cfg.Db)

		nullUserId := uuid.NullUUID{UUID: user.ID, Valid: true}

		posts, err := dbQueries.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
			UserID: nullUserId,
			Limit:  5,
		})

		if err != nil {
			server_actions.RespondWithError(w, http.StatusInternalServerError, "Error getting feed follows")
			return
		}

		server_actions.RespondWithJSON(w, http.StatusOK, posts)
	}
}

func V1Routes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	apiCfg := &config.ApiConfig{
		Db: db,
	}

	r.Get("/readiness", CheckHealth)
	r.Get("/err", SendError)
	r.Post("/users", CreateUser(db))
	r.Get("/users", apiCfg.MiddlewareAuth(GetUser))
	r.Post("/feeds", apiCfg.MiddlewareAuth(CreateFeed(apiCfg)))
	r.Get("/feeds", getFeeds(apiCfg))
	r.Get("/feed_follows", apiCfg.MiddlewareAuth(GetAllFeedFollows(apiCfg)))
	r.Post("/feed_follows", apiCfg.MiddlewareAuth(CreateFeedFollow(apiCfg)))
	r.Delete("/feed_follows", apiCfg.MiddlewareAuth(DeleteFeedFollow(apiCfg)))
	r.Get("/posts", apiCfg.MiddlewareAuth(GetPosts(apiCfg)))
	return r
}
