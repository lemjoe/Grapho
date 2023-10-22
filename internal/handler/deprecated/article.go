package deprecated

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (h *handler) ShowArticle(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang")
	defaultLang := "en"
	localizer := i18n.NewLocalizer(h.bundle, lang, defaultLang)

	artclPath := r.URL.Query().Get("md")
	md, err := os.ReadFile("articles/" + artclPath) // just pass the file name
	if err != nil {
		log.Print("MD file open error: ", err, artclPath)
	}
	// always normalize newlines!
	html := append(MdToHTML(md), toTheTop[:]...)

	homeButton := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "HomeButton",
			Other: "Back to home page",
		},
	})

	doc, err := h.services.ArticleService.GetArticleInfo(artclPath) //RetrieveArticle(artclPath)
	if err != nil {
		log.Println(err)
		return
	}
	//doc.Unmarshal(article)

	HomePageVars := PageVariables{ //store the date and time in a struct
		MDArticle:    template.HTML(html),
		Path:         artclPath,
		Title:        doc.Title,
		HomeButton:   homeButton,
		Author:       doc.Author,
		CreationDate: doc.CreationDate.Format("2006-Jan-02 15:04 MST"),
		UpdateDate:   doc.ModificationDate.Format("2006-Jan-02 15:04 MST"),
	}

	t, err := template.ParseFiles("lib/templates/view.html") //parse the html file homepage.html
	if err != nil {                                          // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

// ArticleList
func (h *handler) ArticleList(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang")
	defaultLang := "en"
	localizer := i18n.NewLocalizer(h.bundle, lang, defaultLang)

	// Translation strings
	listOfArticles := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ListOfArticles",
			Other: "List of available articles:",
		},
	})
	// editButton := localizer.MustLocalize(&i18n.LocalizeConfig{
	// 	DefaultMessage: &i18n.Message{
	// 		ID:    "EditButton",
	// 		Other: "edit",
	// 	},
	// })
	// deleteButton := localizer.MustLocalize(&i18n.LocalizeConfig{
	// 	DefaultMessage: &i18n.Message{
	// 		ID:    "DeleteButton",
	// 		Other: "delete",
	// 	},
	// })
	homeButton := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "HomeButton",
			Other: "Back to home page",
		},
	})
	addButton := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "AddButton",
			Other: "Add an article",
		},
	})
	lastModification := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "LastModification",
			Other: "Last modification",
		},
	})
	pageTitle := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Title",
			Other: "Articles list",
		},
	})

	docs, err := h.services.ArticleService.GetArticlesList()
	if err != nil {
		log.Println(err)
	}

	html := "<h1>" + listOfArticles + "</h1><ul>"
	editImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/edit-pen.png\" alt=\"Edit\" title=\"Edit\">"
	deleteImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/red-trash-can.png\" alt=\"Edit\" title=\"Edit\">"

	if len(docs) == 0 {
		html += "<p>There is no articles here! Why don't you add one?"
	}
	for _, article := range docs {

		html += "<li>" + "<a href='show?md=" + article.FileName + "'>" + article.Title + "</a><i> by <b>" + article.Author + "</b> (" + lastModification + ": " + article.ModificationDate.Format("2006-Jan-02 15:04 MST") + ") </i><a href='edit?md=" + article.FileName + "'><i>" + editImg + "</i></a> | <a href='delete?md=" + article.FileName + "'><i>" + deleteImg + "</i></a></li>"
	}

	html += "</ul>"

	HomePageVars := PageVariables{ //store the date and time in a struct
		MDArticle:  template.HTML(html),
		HomeButton: homeButton,
		AddButton:  addButton,
		Title:      pageTitle,
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

// DeleteArticle
func (h *handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {

	artclPath := r.URL.Query().Get("md")
	err := os.Remove("articles/" + artclPath)
	if err != nil {
		log.Print("MD file delete error: ", err)
	}

	err = h.services.ArticleService.DeleteArticle(artclPath)
	if err != nil {
		log.Print("DB entry delete error: ", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("Successfully Deleted File")
}

// UploadArticle
func (h *handler) UploadArticle(w http.ResponseWriter, r *http.Request) {
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
func (h *handler) DownloadArticle(w http.ResponseWriter, r *http.Request) {

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
func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern

	hash := md5.Sum(fileBytes)
	fileName := hex.EncodeToString(hash[:])
	// tempFile, err := os.Create("articles/" + fileName)
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer tempFile.Close()

	// write this byte array to our temporary file
	//tempFile.Write(fileBytes)
	//!WARN remove admin from author name
	_, err = h.services.ArticleService.CreateNewArticle(fileName, title, "admin", fileBytes) //CreateNewArticle(fileName, title)
	if err != nil {
		log.Println(err)
	}

	// return that we have successfully uploaded our file!
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("Successfully Uploaded File")
}

// SaveFile
func (h *handler) SaveFile(w http.ResponseWriter, r *http.Request) {
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
