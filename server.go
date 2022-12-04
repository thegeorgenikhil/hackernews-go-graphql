package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/thegeorgenikhil/hackernews-go-graphql/graph"
	"github.com/thegeorgenikhil/hackernews-go-graphql/graph/generated"
	"github.com/thegeorgenikhil/hackernews-go-graphql/internal/auth"
	database "github.com/thegeorgenikhil/hackernews-go-graphql/pkg/db/mysql"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	url := os.Getenv("URL") // frontend url

	if port == "" {
		port = defaultPort
	}
	if url == "" {
		url = "http://localhost:3000"
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	// Set up CORS middleware
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{url},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	router.Use(cors.Handler)

	database.InitDB()
	database.Migrate()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
