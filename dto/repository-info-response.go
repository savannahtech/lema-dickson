package dto

type RepositoryInfoResponseDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	HtmlUrl     string `json:"html_url"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Fork        bool   `json:"fork"`
	Language    string `json:"language"`
	ForksCount  int    `json:"forks_count"`
	StarsCount  int    `json:"stargazers_count"`
	OpenIssues  int    `json:"open_issues_count"`
	Watchers    int    `json:"watchers_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
