package models

import "html/template"

type PageVariables struct {
	Md           string
	MDArticle    template.HTML
	HomeButton   string
	AddButton    string
	Title        string
	Path         string
	Author       string
	CreationDate string
	UpdateDate   string
}
