package models

import (
	"gorm.io/gorm"
)

type Commit struct {
	gorm.Model
	RepositoryName string      `gorm:"repository_name"`
	Repository     *Repository `gorm:"foreignKey:RepositoryName"`
	Message        string      `gorm:"message" json:"message"`
	Author         string      `gorm:"author" json:"author"`
	Date           string      `gorm:"string" json:"date"`
	URL            string      `gorm:"html_url" json:"html_url"`
	SHA            string      `gorm:"sha" json:"sha"`
}
