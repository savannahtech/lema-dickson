package controllers

import (
	"log"
	"net/http"

	"github.com/midedickson/github-service/utils"
)

func (c *Controller) GetRepositoryInfo(w http.ResponseWriter, r *http.Request) {
	owner, err := utils.GetPathParam(r, "owner")
	if err != nil {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	if owner == "" {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	repoName, err := utils.GetPathParam(r, "repo")
	if err != nil {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	if repoName == "" {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	user, err := c.dbRepository.GetUser(owner)
	if err != nil {
		utils.Dispatch500Error(w, err)
		return
	}
	if user == nil {
		utils.Dispatch404Error(w, "User with this github username not found, please register this github username", err)
		return
	}
	repo, err := c.dbRepository.GetRepository(user.ID, repoName)
	if err != nil {
		utils.Dispatch500Error(w, err)
		return
	}
	if repo == nil {
		go c.task.AddRequestToFetchNewlyRequestedRepoQueue(user.Username, repoName)
		utils.Dispatch404Error(w, "Repository not found on Github; kindly check back again.", err)
		return
	}

	utils.Dispatch200(w, "Repository Information Fetched Successfully", repo)
}

func (c *Controller) GetRepositoryCommits(w http.ResponseWriter, r *http.Request) {
	repoName, err := utils.GetPathParam(r, "repo")
	if err != nil {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	if repoName == "" {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	commits, err := c.dbRepository.GetRepositoryCommits(repoName)
	if err != nil {
		log.Printf("%v", err)
		utils.Dispatch500Error(w, err)
		return
	}
	utils.Dispatch200(w, "Repository Commits Fetched Successfully", commits)
}

func (c *Controller) GetRepositories(w http.ResponseWriter, r *http.Request) {
	repoSearchParams := &utils.RepositorySearchParams{}
	owner, err := utils.GetPathParam(r, "owner")
	if err != nil || owner == "" {
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}
	user, err := c.dbRepository.GetUser(owner)
	if err != nil {
		utils.Dispatch500Error(w, err)
		return
	}
	if user == nil {
		utils.Dispatch404Error(w, "User with this github username not found, please register this github username", err)
		return
	}
	utils.ParseQueryParams(r, repoSearchParams)
	repositories, err := c.dbRepository.SearchRepository(user.ID, repoSearchParams)
	if err != nil {
		log.Printf("%v", err)
		utils.Dispatch500Error(w, err)
		return
	}
	utils.Dispatch200(w, "Repositories Fetched Successfully", repositories)
}
