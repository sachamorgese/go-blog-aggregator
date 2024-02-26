package main

import (
	"BlogAggregator/config"
	"BlogAggregator/routes"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
		fmt.Println(err.Error())
		return
	}

	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		fmt.Println("Error connecting to the database")
		panic(err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("PORT is not set")
		return
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(config.CorsOptions))

	r.Mount("/v1", routes.V1Routes(db))

	server := http.Server{
		Handler: r,
		Addr:    ":" + port,
	}

	fmt.Println("Server is running on port " + port)
	errServer := server.ListenAndServe()

	if errServer != nil {
		panic(errServer)
	}

}
