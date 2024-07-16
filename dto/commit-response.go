package dto

import "encoding/json"

type CommitResponseDTO struct {
	SHA     string `json:"sha"`
	Message string
	Author  string
	Date    string
	URL     string `json:"html_url"`
}

type nestedCommit struct {
	Message string `json:"message"`
	Author  struct {
		Name string `json:"name"`
		Date string `json:"date"`
	} `json:"author"`
}

type tempCommitResponseDTO struct {
	SHA    string       `json:"sha"`
	Commit nestedCommit `json:"commit"`
	URL    string       `json:"html_url"`
}

func (c *CommitResponseDTO) UnmarshalJSON(data []byte) error {
	var temp tempCommitResponseDTO
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	c.SHA = temp.SHA
	c.Message = temp.Commit.Message
	c.Author = temp.Commit.Author.Name
	c.Date = temp.Commit.Author.Date
	c.URL = temp.URL
	return nil
}
