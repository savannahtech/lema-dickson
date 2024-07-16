package controllers

import (
	"github.com/midedickson/github-service/database"
	"github.com/midedickson/github-service/requester"
	"github.com/midedickson/github-service/tasks"
)

type Controller struct {
	requester    requester.Requester
	dbRepository database.DBRepository
	task         tasks.Task
}

func NewController(
	requester requester.Requester,
	dbRepository database.DBRepository,
	task tasks.Task,
) *Controller {
	return &Controller{
		requester:    requester,
		dbRepository: dbRepository,
		task:         task,
	}
}
