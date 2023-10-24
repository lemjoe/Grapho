package models

import "time"

type Article struct {
	// FileName         string    `json:"file_name"`
	Title            string    `json:"article_title"`
	Author           string    `json:"article_author"`
	AuthorId         string    `json:"author_id"`
	CreationDate     time.Time `json:"creation_date"`
	ModificationDate time.Time `json:"modification_date"`
	IsLocked         bool      `json:"is_locked"`
	Id               string    `json:"id"`
}
