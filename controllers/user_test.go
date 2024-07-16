package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/midedickson/github-service/controllers"
	"github.com/midedickson/github-service/dto"
	"github.com/midedickson/github-service/mocks"
	"github.com/midedickson/github-service/models"
	"github.com/midedickson/github-service/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	// Initialize the mocks
	mockDBRepository := new(mocks.MockDBRepository)
	mockRequester := new(mocks.MockRequester)
	mockTask := new(mocks.MockTask)

	// Create the controller with mocked dependencies
	controller := controllers.NewController(mockRequester, mockDBRepository, mockTask)

	// Define the input payload and the expected user
	createUserPayload := &dto.CreateUserPayloadDTO{
		Username: "testuser",
	}
	user := &models.User{
		Username: "testuser",
	}

	// Set up the expectations
	mockDBRepository.On("CreateUser", createUserPayload).Return(user, nil)
	var wg sync.WaitGroup
	wg.Add(1)
	mockTask.On("AddUserToGetAllRepoQueue", user).Run(func(args mock.Arguments) {
		wg.Done()
	}).Return()

	// Create a new HTTP request with the input payload
	body, _ := json.Marshal(createUserPayload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// Call the CreateUser method
	controller.CreateUser(rr, req)
	wg.Wait()

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, rr.Code)
	var response utils.APIResponse
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, true, response.Success)
	assert.Equal(t, "user created successfully", response.Message)

	// Assert that the expectations were met
	mockDBRepository.AssertExpectations(t)
	mockTask.AssertExpectations(t)
	mockRequester.AssertExpectations(t)
}
