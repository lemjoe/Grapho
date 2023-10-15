package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"time"

	c "github.com/ostafen/clover/v2"
	d "github.com/ostafen/clover/v2/document"
	q "github.com/ostafen/clover/v2/query"
)

var username = "admin"

func CreateDefaultDB() {

	err := os.Mkdir("db", 0644)
	if err != nil {
		log.Print("Unable to create database: ", err)
	}
	err = os.Mkdir("articles", 0644)
	if err != nil {
		log.Print("Unable to create articles folder: ", err)
	}
	db, err := c.Open("db")
	if err != nil {
		log.Print("Unable to connect database: ", err)
	}
	err = db.CreateCollection("users")
	if err != nil {
		log.Print("Unable to create collection: ", err)
	}
	err = db.CreateCollection("articles")
	if err != nil {
		log.Print("Unable to create collection: ", err)
	}

	doc := d.NewDocument()
	doc.Set("user_name", "admin")
	doc.Set("full_name", "Administrator")
	doc.Set("creation_date", time.Now())
	doc.Set("last_login", time.Date(2004, time.March, 19, 17, 36, 0, 0, time.UTC))
	doc.Set("is_admin", true)
	docId, err := db.InsertOne("users", doc)
	if err != nil {
		log.Print("Unable to insert document: ", err)
	} else {
		log.Println("New object was created and has an id - ", docId)
	}

	md := []byte("# Welcome to **md-blog**, your Personal Blog/Wiki Page\n\nHere you can store, organize and collaborate on information in a way that suits you best. Create, explore and share your knowledge with ease!\n\n* Quickly access and edit your notes using the intuitive web-based interface\n* Customize your pages using Markdown markup language\n* Collaborate with friends and colleagues by inviting them to view or edit specific pages\n* Write your thoughts and share them with everyone ")

	hash := md5.Sum(md)
	fileName := hex.EncodeToString(hash[:])
	err = os.WriteFile("articles/"+fileName, md, 0644)
	if err != nil {
		log.Print("MD file write error: ", err)
	} else {
		log.Println("Successfully Written File")
	}
	doc = d.NewDocument()
	doc.Set("file_name", fileName)
	doc.Set("article_title", "Welcoming words")
	doc.Set("article_author", "Alexander Ignatov")
	doc.Set("creation_date", time.Now())
	doc.Set("modification_date", time.Now())
	doc.Set("is_locked", false)
	docId, err = db.InsertOne("articles", doc)
	if err != nil {
		log.Print("Unable to insert document: ", err)
	} else {
		log.Println("New object was created and has an id - ", docId)
	}

	db.CreateIndex("users", "user_name")
	if err != nil {
		log.Print("Unable to create index: ", err)
	}
	db.CreateIndex("articles", "file_name")
	if err != nil {
		log.Print("Unable to create index: ", err)
	}

	defer db.Close()

	if err != nil {
		log.Println("Something went wrong. Read the log above")
	} else {
		log.Println("Default database was successfully created")
	}
}

func ReadArticlesList() ([]*d.Document, error) {

	db, err := c.Open("db")
	if err != nil {
		return nil, err
	}

	docs, err := db.FindAll(q.NewQuery("articles"))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return docs, nil
}

func RetrieveArticle(fileName string) (*d.Document, error) {

	db, err := c.Open("db")
	if err != nil {
		return nil, err
	}

	doc, err := db.FindFirst(q.NewQuery("articles").Where(q.Field("file_name").Eq(fileName)))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return doc, nil
}

func CreateNewArticle(fileName string, title string) error {

	db, err := c.Open("db")
	if err != nil {
		return err
	}
	fn, err := db.FindFirst(q.NewQuery("users").Where(q.Field("user_name").Eq(username)))
	if err != nil {
		return err
	}
	doc := d.NewDocument()
	doc.Set("file_name", fileName)
	doc.Set("article_title", title)
	doc.Set("article_author", fn.Get("full_name"))
	doc.Set("creation_date", time.Now())
	doc.Set("modification_date", time.Now())
	doc.Set("is_locked", false)
	docId, err := db.InsertOne("articles", doc)
	if err != nil {
		return err
	} else {
		log.Println("New object was created and has an id - ", docId)
	}
	defer db.Close()

	return nil
}

func DeleteArticleFromDB(fileName string) error {

	db, err := c.Open("db")
	if err != nil {
		return err
	}
	err = db.Delete(q.NewQuery("articles").Where(q.Field("file_name").Eq(fileName)))
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}

func UpdateArticle(fileName string) error {

	db, err := c.Open("db")
	if err != nil {
		return err
	}
	changes := make(map[string]interface{})
	changes["modification_date"] = time.Now()
	changes["is_locked"] = false
	log.Println(changes)
	err = db.Update(q.NewQuery("articles").Where(q.Field("file_name").Eq(fileName)), changes)

	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}
