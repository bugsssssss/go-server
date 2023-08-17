package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bugsssssss/rssag/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// ? struct for db queries
type apiConfig struct {
	DB *database.Queries
}

func main() {

	fmt.Println("Hello welcome to my first API server on GO:")

	// ? saying that our env is in .env file
	godotenv.Load(".env")

	// ? Pulling PORT env from .env file
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the env")
	}
	// ? Pulling DB env from .env file
	dbURL := os.Getenv("DB_URL")

	// ? checking if there is any error while getting dbURL from .env
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the env")
	}

	// ? connecting to our postgres db by url
	conn, err := sql.Open("postgres", dbURL)

	// ? checking
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	// ? if everything is okay then we create new cpnnection and new config
	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	// ? assign NewRouter imported from chi
	router := chi.NewRouter()

	// ? standart conf for handling
	// ? telling methods we are gonna need and link types
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// ? creating var for the v1 api endpoint
	v1Router := chi.NewRouter()

	// ? API endpoints, takes a handler as second parameter
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Get("/test", handlerReadiness)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeed)
	v1Router.Delete("/feeds/{feedID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeed))

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	// ? it add before every link phrase '/v1'
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
