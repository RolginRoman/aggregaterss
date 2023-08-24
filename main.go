package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/rolginroman/aggregaterss/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("There is no PORT env variable")
	}

	dbUrlString := os.Getenv("DB_URL")

	if dbUrlString == "" {
		log.Fatal("There is no DB_URL env variable")
	}

	connection, err := sql.Open("postgres", dbUrlString)
	if err != nil {
		log.Fatal("Can't connect to DB", err)
	}

	apiCfg := apiConfig{
		DB: database.New(connection),
	}

	go startScraping(apiCfg.DB, 10, time.Minute*1)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/ready", handlerReadiness)
	v1Router.Post("/error", handlerError)
	v1Router.Post("/users", apiCfg.handlerUserCreate)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUserGetByApiKey))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedCreate))
	v1Router.Get("/feeds", apiCfg.handlerFeeds)
	v1Router.Get("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedsByUser))
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowCreate))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsByUser))
	v1Router.Delete("/feed_follows/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowDelete))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server's starting on port %v", portString)

	error := server.ListenAndServe()
	if error != nil {
		log.Fatal(error)
	}
}
