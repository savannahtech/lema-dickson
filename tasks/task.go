package tasks

import (
	"github.com/midedickson/github-service/models"
)

type Task interface {
	AddUserToGetAllRepoQueue(user *models.User)
	AddRequestToFetchNewlyRequestedRepoQueue(username, repoName string)
}
