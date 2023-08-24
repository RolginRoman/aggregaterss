package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rolginroman/aggregaterss/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
}

func CreateRouter(api ApiConfig) http.Handler {
	v1Router := chi.NewRouter()

	v1Router.Get("/ready", handlerReadiness)
	v1Router.Post("/error", handlerError)

	v1Router.Post("/users", api.handlerUserCreate)
	v1Router.Get("/users", api.middlewareAuth(api.handlerUserGetByApiKey))

	v1Router.Post("/feeds", api.middlewareAuth(api.handlerFeedCreate))
	v1Router.Get("/feeds", api.handlerFeeds)
	v1Router.Get("/feeds", api.middlewareAuth(api.handlerFeedsByUser))

	v1Router.Post("/feed_follows", api.middlewareAuth(api.handlerFeedFollowCreate))
	v1Router.Get("/feed_follows", api.middlewareAuth(api.handlerFeedFollowsByUser))
	v1Router.Delete("/feed_follows/{feedFollowId}", api.middlewareAuth(api.handlerFeedFollowDelete))

	v1Router.Get("/posts", api.middlewareAuth(api.handlerGetPostsByUser))

	return v1Router
}
