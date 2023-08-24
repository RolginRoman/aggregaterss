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

func (apiConfig *ApiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	createdUser, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJSON(w, 201, models.ConvertUserModelToExternal(createdUser))
}

func (apiConfig *ApiConfig) handlerUserGetByApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, models.ConvertUserModelToExternal(user))
}
