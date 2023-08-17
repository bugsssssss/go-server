package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bugsssssss/rssag/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v ", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldnt create feed: %v", err))
	}

	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseFeedsToFeeds(feeds))
}

func (apiCfg *apiConfig) handlerDeleteFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIDStr := chi.URLParam(r, "feedID")

	feedID, err := uuid.Parse(feedIDStr)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeed(r.Context(), database.DeleteFeedParams{
		ID:     feedID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed: %v", err))
	}

	respondWithJSON(w, 200, struct{}{})
}
