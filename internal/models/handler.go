package models

import (
	"html/template"
)

type PageVariables struct {
	Id                  string
	Md                  string
	MDArticle           template.HTML
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
	ManagedUserIsWriter string
	AdminPanel          bool
	Translation         map[string]string
	Settings            map[string]string
	Licenses            []string
	LicensesPath        []string
	IsWriter            bool
	ToTheTop            bool
	UsersInfo           [][7]string
	ArticlesInfo        [][5]string
}
