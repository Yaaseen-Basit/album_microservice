package main

import (
	"context"
	"flag"
	"log"
	"net/http" // Import the standard "net/http" package
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"albumproject/internal/album"
	localhttp "albumproject/internal/transport/http"
)

func main() {
	// Parse command line arguments
	port := flag.String("port", "8080", "HTTP port for the server")
	flag.Parse()

	// Create MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Create repository instance
	repo := album.NewMongoRepository(client, "albumdb", "albummeta")
	// Create service instance
	service := album.NewAlbumService(repo)

	// Create HTTP router and server
	router := mux.NewRouter()
	// http.MakeHTTPHandlers(router, service) // Commented out because it's not used
	localhttp.MakeHTTPHandlers(router, service)
	server := &http.Server{
		Addr:    ":" + *port,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", *port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server gracefully stopped")
}
