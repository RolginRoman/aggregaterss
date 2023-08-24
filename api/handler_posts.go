package api

import (
	"fmt"
	"net/http"

	"github.com/rolginroman/aggregaterss/api/models"
	"github.com/rolginroman/aggregaterss/internal/database"
)

func (apiConfig *ApiConfig) handlerGetPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error retrieving posts: %v", err))
		return
	}
	respondWithJSON(w, 200, models.ConvertPostModelsToExternal(posts))
}
