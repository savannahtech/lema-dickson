package models

import (
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	RemoteID        int    `gorm:"remote_id"`
	OwnerID         uint   `gorm:"owner_id"`
	Owner           *User  `gorm:"foreignKey:OwnerID"`
	Name            string `gorm:"name"`
	Description     string `gorm:"description"`
	URL             string `gorm:"html_url"`
	Language        string `gorm:"language"`
	ForksCount      int    `gorm:"forks_count"`
	StarsCount      int    `gorm:"stargazers_count"`
	OpenIssues      int    `gorm:"open_issues_count"`
	Watchers        int    `gorm:"watchers_count"`
	RemoteCreatedAt string `gorm:"remote_created_at"`
	RemoteUpdatedAt string `gorm:"remote_updated_at"`
}
