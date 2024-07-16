package utils

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RepositorySearchParams struct {
	Name          string `json:"name"`
	Language      string `json:"language"`
	TopStarsCount int    `json:"stars_count"`
}
