package models

import "html/template"

type PageVariables struct {
	Id                  string
	Md                  string
	MDArticle           template.HTML
	HomeButton          string
	AddButton           string
	Title               string
	Path                string
	Author              string
	CreationDate        string
	UpdateDate          string
	UserName            string
	BodyLoudText        string
	BodyText            string
	Theme               string
	ManagedUserName     string
	ManagedUserFullName string
	ManagedUserEmail    string
	ManagedUserIsAdmin  string
	AdminPanel          template.HTML
}
