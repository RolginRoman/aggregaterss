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
	"github.com/rolginroman/aggregaterss/api"
	"github.com/rolginroman/aggregaterss/internal/database"
	"github.com/rolginroman/aggregaterss/internal/scraper"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("There is no PORT env variable")
	}

	connection, err := initDbConnection()
	if err != nil {
		log.Fatal("Cannot establish DB connection", err)
	}
	apiCfg := api.ApiConfig{
		DB: database.New(connection),
	}

	go scraper.StartScraping(apiCfg.DB, 10, time.Minute*1)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Mount("/v1", api.CreateRouter(apiCfg))

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

func initDbConnection() (*sql.DB, error) {
	dbUrlString := os.Getenv("DB_URL")
	if dbUrlString == "" {
		log.Fatal("There is no DB_URL env variable")
	}

	connection, err := sql.Open("postgres", dbUrlString)
	if err != nil {
		log.Fatal("Can't connect to DB", err)
	}
	return connection, nil
}
