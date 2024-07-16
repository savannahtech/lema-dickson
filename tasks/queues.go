package tasks

import (
	"log"

	"github.com/midedickson/github-service/models"
)

func (t *AsyncTask) AddUserToGetAllRepoQueue(user *models.User) {
	t.GetAllRepoForUserQueue <- user
}

func (t *AsyncTask) AddRequestToFetchNewlyRequestedRepoQueue(username, repoName string) {
	log.Println("Adding request to fetch newly requested")

	t.FetchNewlyRequestedRepoQueue <- &RepoRequest{
		Username: username,
		RepoName: repoName,
	}
	log.Println("Added request to fetch newly requested")

}

func (t *AsyncTask) AddSignalToCheckForUpdateOnAllRepoQueue() {
	t.CheckForUpdateOnAllRepoQueue <- "signal"
}
