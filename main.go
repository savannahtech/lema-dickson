package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/midedickson/github-service/controllers"
	"github.com/midedickson/github-service/database"
	"github.com/midedickson/github-service/requester"
	"github.com/midedickson/github-service/routes"
	"github.com/midedickson/github-service/tasks"
)

func main() {
	log.Println("Starting server...")

	database.ConnectToDB()
	database.AutoMigrate()
	// Use a WaitGroup to manage goroutines
	var wg sync.WaitGroup

	repoRequester := requester.NewRepositoryRequester()
	dbRepository := database.NewSqliteDBRepository(database.DB)
	tasks := tasks.NewAsyncTask(repoRequester, dbRepository)
	controller := controllers.NewController(repoRequester, dbRepository, tasks)

	// Start goroutines to fetch repositories and check for updates
	wg.Add(1)
	go tasks.GetAllRepoForUser(&wg)
	wg.Add(1)
	go tasks.FetchNewlyRequestedRepo(&wg)
	wg.Add(1)
	go tasks.CheckForUpdateOnAllRepo(&wg)
	go tasks.AddSignalToCheckForUpdateOnAllRepoQueue()

	// create mux router
	r := mux.NewRouter()
	routes.ConnectRoutes(r, controller)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	server := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	log.Println("Server started on :8080")
	<-stop
	log.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close channels to signal workers to stop
	close(tasks.GetAllRepoForUserQueue)
	close(tasks.FetchNewlyRequestedRepoQueue)
	close(tasks.CheckForUpdateOnAllRepoQueue)

	// Wait for all goroutines to complete
	wg.Wait()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
