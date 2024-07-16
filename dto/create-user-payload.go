package dto

type CreateUserPayloadDTO struct {
	Username string `json:"username"`
	FullName string `json:"fullName"`
}
