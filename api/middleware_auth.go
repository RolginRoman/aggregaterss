package api

import (
	"fmt"
	"net/http"

	"github.com/rolginroman/aggregaterss/internal/auth"
	"github.com/rolginroman/aggregaterss/internal/database"
)

type authenticatedRequestHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *ApiConfig) middlewareAuth(handler authenticatedRequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth, error := auth.GetApiKey(r.Header)
		if error != nil {
			respondWithError(w, 403, fmt.Sprintf("Invalid authentication info: %v", error))
			return
		}
		user, err := apiConfig.DB.GetUserByApiKey(r.Context(), auth)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("User not found: %v", err))
			return
		}

		handler(w, r, user)
	}
}
