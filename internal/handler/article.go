package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lemjoe/md-blog/internal/models"
)

// Home page (Articles list)
func (h *Handler) GetArticlesList(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang")
	translation := Localizer([]string{"listOfArticles", "homeButton", "addButton", "lastModification", "pageTitle"}, lang, h.bundle)
	// log.Println(translation)
	docs, err := h.services.ArticleService.GetArticlesList()
	if err != nil {
		log.Println(err)
	}

	html := "<h1>" + translation["listOfArticles"] + "</h1><ul>"
	editImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/edit-pen.png\" alt=\"Edit\" title=\"Edit\">"
	deleteImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/red-trash-can.png\" alt=\"Edit\" title=\"Edit\">"

	if len(docs) == 0 {
		html += "<p>There is no articles here! Why don't you add one?"
	}
	for _, article := range docs {
		log.Println(article)
		html += "<li>" + "<a href='show?md=" + article.Id + "'>" + article.Title + "</a><i> by <b>" + article.Author + "</b> (" + translation["lastModification"] + ": " + article.ModificationDate.Format("2006-Jan-02 15:04 MST") + ") </i><a href='edit?md=" + article.Id + "'><i>" + editImg + "</i></a> | <a href='delete?md=" + article.Id + "'><i>" + deleteImg + "</i></a></li>"
	}

	html += "</ul>"

	HomePageVars := models.PageVariables{ //store the date and time in a struct
		MDArticle:  template.HTML(html),
		HomeButton: translation["homeButton"],
		AddButton:  translation["addButton"],
		Title:      translation["pageTitle"],
	}

	t, err := template.ParseFiles("lib/templates/home.html") //parse the html file homepage.html
	if err != nil {                                          // if there is an error
		log.Print("Template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("Template executing error: ", err) //log it
	}

}

// Show article
func (h *Handler) ShowArticle(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang")
	translation := Localizer([]string{"homeButton"}, lang, h.bundle)

	artclPath := r.URL.Query().Get("md")
	md, err := h.services.FileService.ReadFile("articles/" + artclPath)
	if err != nil {
		log.Print("MD file open error: ", err, artclPath)
	}
	// always normalize newlines!
	html := append(MdToHTML(md), toTheTop[:]...)

	doc, err := h.services.ArticleService.GetArticleInfo(artclPath) //RetrieveArticle(artclPath)
	if err != nil {
		log.Println(err)
		return
	}
	//doc.Unmarshal(article)

	ArticlePageVars := models.PageVariables{ //store the date and time in a struct
		MDArticle:    template.HTML(html),
		Path:         artclPath,
		Title:        doc.Title,
		HomeButton:   translation["homeButton"],
		Author:       doc.Author,
		CreationDate: doc.CreationDate.Format("2006-Jan-02 15:04 MST"),
		UpdateDate:   doc.ModificationDate.Format("2006-Jan-02 15:04 MST"),
	}

	t, err := template.ParseFiles("lib/templates/view.html") //parse the html file homepage.html
	if err != nil {                                          // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, ArticlePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

// DeleteArticle
func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {

	artclPath := r.URL.Query().Get("md")

	err := h.services.ArticleService.DeleteArticle(artclPath)
	if err != nil {
		log.Print("DB entry delete error: ", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("Successfully Deleted File")
}

// UploadArticle
func (h *Handler) UploadArticle(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("lib/templates/upload.html") //parse the html file homepage.html
	if err != nil {                                            // if there is an error
		log.Print("template parsing error: ", err) // log it
		fmt.Fprintln(w, "template parsing error: ", err)
		return
	}
	err = t.Execute(w, time.Now()) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                // if there is an error
		log.Print("template executing error: ", err) //log it
		fmt.Fprintln(w, "template executing error: ", err)
		return
	}
}

// DownloadArticle
func (h *Handler) DownloadArticle(w http.ResponseWriter, r *http.Request) {

	artclPath := r.URL.Query().Get("md")
	md, err := h.services.ArticleService.GetArticleBody(artclPath)
	if err != nil {
		log.Print("MD file open error: ", err)
		fmt.Fprintln(w, "MD file open error: ", err)
		return
	}

	doc, err := h.services.ArticleService.GetArticleInfo(artclPath)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		return
	}

	fileName := strings.ReplaceAll(doc.Title, " ", "_") + ".md"

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "text/markdown")
	w.Header().Set("Content-Length", strconv.Itoa(len(md)))

	//stream the body to the client without fully loading it into memory
	reader := bytes.NewReader(md)
	_, err = io.Copy(w, reader)
	if err != nil {
		log.Print("Unable to download a file: ", err)
		fmt.Fprintln(w, "Unable to download a file: ", err)
		return
	}

}

// Upload
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	log.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	title := r.FormValue("title")
	if err != nil {
		log.Print("Error Retrieving the File", err)
		return
	}

	ftype := handler.Filename[len(handler.Filename)-3:]
	if ftype != ".md" {
		log.Println("File type must be text/markdown" + ftype)
		return
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error Reading File", err)
		return
	}

	_, err = h.services.ArticleService.CreateNewArticle(title, "admin", fileBytes) //CreateNewArticle(fileName, title)
	if err != nil {
		log.Println("Error Creating Article", err)
		return
	}

	// return that we have successfully uploaded our file!
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("Successfully Uploaded File")
}

// SaveFile
func (h *Handler) SaveFile(w http.ResponseWriter, r *http.Request) {
	md := []byte(r.FormValue("textEditArea"))
	artclPath := r.FormValue("articlePath")
	err := os.WriteFile("articles/"+artclPath, md, 0644)
	if err != nil {
		log.Print("MD file write error: ", err, artclPath)
		return
	} else {
		log.Println("Successfully Edited File")
	}

	err = h.services.ArticleService.UpdateArticle(artclPath)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
