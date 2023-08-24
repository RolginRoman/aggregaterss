package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rolginroman/aggregaterss/api/models"
	"github.com/rolginroman/aggregaterss/internal/database"
)

func (apiConfig *ApiConfig) handlerFeedCreate(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	createdFeed, err := apiConfig.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed: %v", err))
		return
	}
	respondWithJSON(w, 201, models.ConvertFeedModelToExternal(createdFeed))
}

func (apiConfig *ApiConfig) handlerFeedsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiConfig.DB.GetFeedsByUserId(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error retrieving feeds: %v", err))
		return
	}
	respondWithJSON(w, 200, models.ConvertFeedModelsToExternal(feeds))

}

func (apiConfig *ApiConfig) handlerFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiConfig.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error retrieving feeds: %v", err))
		return
	}
	respondWithJSON(w, 200, models.ConvertFeedModelsToExternal(feeds))

}
