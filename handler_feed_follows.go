package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/rolginroman/aggregaterss/internal/database"
)

func (apiConfig *apiConfig) handlerFeedFollowCreate(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	createdFeedFollow, err := apiConfig.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed follow: %v", err))
		return
	}
	respondWithJSON(w, 201, convertFeedFollowModelToExternal(createdFeedFollow))
}

func (apiConfig *apiConfig) handlerFeedFollowsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiConfig.DB.GetFeedsForUser(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error retrieving feed follows: %v", err))
		return
	}
	respondWithJSON(w, 200, convertFeedModelsToExternal(feeds))

}

func (apiConfig *apiConfig) handlerFeedFollowDelete(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowIdString := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(feedFollowIdString)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Incorrect request: %v", err))
	}

	err = apiConfig.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed follows: %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})

}
