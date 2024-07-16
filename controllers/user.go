package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/utils"
)

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Create user logic
	var createUserPayload dto.CreateUserPayloadDTO
	err := json.NewDecoder(r.Body).Decode(&createUserPayload)
	if err != nil {
		log.Printf("Error decoding create user payload: %v", err)
		utils.Dispatch400Error(w, "Invalid Payload", err)
		return
	}

	user, err := c.dbRepository.CreateUser(&createUserPayload)
	if err != nil {
		utils.Dispatch500Error(w, err)
		return
	}
	go c.task.AddUserToGetAllRepoQueue(user)
	utils.Dispatch200(w, "user created successfully", user)
}
