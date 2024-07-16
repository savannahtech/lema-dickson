package tasks

import (
	"github.com/midedickson/github-service/database"
	"github.com/midedickson/github-service/models"
	"github.com/midedickson/github-service/requester"
)

type AsyncTask struct {
	GetAllRepoForUserQueue       chan *models.User
	FetchNewlyRequestedRepoQueue chan *RepoRequest
	CheckForUpdateOnAllRepoQueue chan string
	requester                    requester.Requester
	dbRepository                 database.DBRepository
}

func NewAsyncTask(requester requester.Requester, dbRepository database.DBRepository) *AsyncTask {
	return &AsyncTask{
		GetAllRepoForUserQueue:       make(chan *models.User),
		FetchNewlyRequestedRepoQueue: make(chan *RepoRequest),
		CheckForUpdateOnAllRepoQueue: make(chan string),
		requester:                    requester,
		dbRepository:                 dbRepository,
	}
}
