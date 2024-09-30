package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/lemjoe/Grapho/internal/models"
	"github.com/lemjoe/Grapho/internal/service"
)

// Home page (Articles list)
func (h *Handler) GetArticlesList(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	theme := curUser.Settings["theme"]
	lang := curUser.Settings["language"]
	translation := Localizer([]string{"listOfArticles", "homeButton", "addButton", "lastModification", "pageTitle"}, lang, h.bundle)

	docs, err := h.services.ArticleService.GetArticlesList()
	if err != nil {
		logger.Error(err)
	}

	html := "<h1>" + translation["listOfArticles"] + "</h1><ul>"
	editImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/" + theme + "/edit-pen.png\" alt=\"Edit\" title=\"Edit\">"
	deleteImg := "<img style=\"padding: 0px; display: inline-block\" width=\"16\" height=\"16\" src=\"../images/" + theme + "/red-trash-can.png\" alt=\"Edit\" title=\"Edit\">"

	if len(docs) == 0 {
		html += "<p>There is no articles here! Why don't you add one?"
	}
	for _, article := range docs {
		logger.Info(article)
		html += "<li>" + "<a href='show?md=" + article.Id + "'>" + article.Title + "</a><i> by <b>" + article.Author + "</b> (" + translation["lastModification"] + ": " + article.ModificationDate.Format("2006-Jan-02 15:04 MST") + ") </i><a href='edit?md=" + article.Id + "'><i>" + editImg + "</i></a> | <a href='delete?md=" + article.Id + "'><i>" + deleteImg + "</i></a></li>"
	}

	html += "</ul>"

	HomePageVars := models.PageVariables{ //store the date and time in a struct
		MDArticle:  template.HTML(html),
		HomeButton: translation["homeButton"],
		AddButton:  translation["addButton"],
		Title:      translation["pageTitle"],
		UserName:   curUser.FullName,
		Theme:      theme,
	}

	t, err := template.ParseFiles("lib/templates/home.html") //parse the html file homepage.html
	if err != nil {                                          // if there is an error
		logger.Error("Template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		logger.Error("Template executing error: ", err) //log it
	}

}

// Show article
func (h *Handler) ShowArticle(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	lang := curUser.Settings["language"]
	translation := Localizer([]string{"homeButton"}, lang, h.bundle)

	artId := r.URL.Query().Get("md")
	md, err := h.services.FileService.ReadFile("articles/" + artId)
	if err != nil {
		logger.Error("MD file open error: ", err, artId)
	}
	// always normalize newlines!
	html := append(MdToHTML(md), toTheTop[:]...)

	doc, err := h.services.ArticleService.GetArticleInfo(artId)
	if err != nil {
		logger.Error(err)
		return
	}
	//doc.Unmarshal(article)

	ArticlePageVars := models.PageVariables{ //store the date and time in a struct
		MDArticle:    template.HTML(html),
		Path:         artId,
		Id:           artId,
		Title:        doc.Title,
		HomeButton:   translation["homeButton"],
		Author:       doc.Author,
		CreationDate: doc.CreationDate.Format("2006-Jan-02 15:04 MST"),
		UpdateDate:   doc.ModificationDate.Format("2006-Jan-02 15:04 MST"),
		UserName:     curUser.FullName,
		Theme:        curUser.Settings["theme"],
	}

	t, err := template.ParseFiles("lib/templates/view.html") //parse the html file homepage.html
	if err != nil {                                          // if there is an error
		logger.Error("template parsing error: ", err) // log it
	}
	err = t.Execute(w, ArticlePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                     // if there is an error
		logger.Error("template executing error: ", err) //log it
	}
}

// DeleteArticle
func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		logger.Error("Unauthorized status code 401")
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	artclPath := r.URL.Query().Get("md")

	doc, err := h.services.ArticleService.GetArticleInfo(artclPath)
	if err != nil {
		logger.Error(err)
		return
	}

	// Send 403 wrong user
	if !curUser.IsAdmin && curUser.Id != doc.AuthorId {
		logger.Error("Wrong user. Action forbidden: status code 403")
		h.SendCode(w, r, statusCodes[http.StatusForbidden])
		return
	} else {
		logger.Info("OK user!")
		err := h.services.ArticleService.DeleteArticle(artclPath)
		if err != nil {
			logger.Error("DB entry delete error: ", err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		logger.Info("Successfully Deleted File")
	}
}

// UploadArticle
func (h *Handler) UploadArticle(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	lang := curUser.Settings["language"]
	translation := Localizer([]string{"homeButton"}, lang, h.bundle)

	t, err := template.ParseFiles("lib/templates/upload.html") //parse the html file homepage.html
	if err != nil {                                            // if there is an error
		logger.Error("template parsing error: ", err) // log it
		return
	}
	UploadPageVars := models.PageVariables{ //store the date and time in a struct
		HomeButton: translation["homeButton"],
		UserName:   curUser.FullName,
		Theme:      curUser.Settings["theme"],
	}
	err = t.Execute(w, UploadPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                    // if there is an error
		logger.Error("template executing error: ", err) //log it
		fmt.Fprintln(w, "template executing error: ", err)
		return
	}
}

// DownloadArticle
func (h *Handler) DownloadArticle(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	artclPath := r.URL.Query().Get("md")
	md, err := h.services.ArticleService.GetArticleBody(artclPath)
	if err != nil {
		logger.Error("MD file open error: ", err)
		fmt.Fprintln(w, "MD file open error: ", err)
		return
	}

	doc, err := h.services.ArticleService.GetArticleInfo(artclPath)
	if err != nil {
		logger.Error(err)
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
		logger.Error("Unable to download a file: ", err)
		fmt.Fprintln(w, "Unable to download a file: ", err)
		return
	}

}

// Upload
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	curUserString := curUser.UserName
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	logger.Info("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 1 MB files.
	r.ParseMultipartForm(1 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	title := r.FormValue("title")
	if err != nil {
		logger.Error("Error Retrieving the File", err)
		return
	}

	ftype := handler.Filename[len(handler.Filename)-3:]
	if ftype != ".md" {
		logger.Error("File type must be text/markdown" + ftype)
		return
	}
	defer file.Close()
	logger.Infof("Uploaded File: %+v\n", handler.Filename)
	logger.Infof("File Size: %+v\n", handler.Size)
	logger.Infof("MIME Header: %+v\n", handler.Header)

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Error Reading File", err)
		return
	}

	_, err = h.services.ArticleService.CreateNewArticle(title, curUserString, fileBytes)
	if err != nil {
		logger.Error("Error Creating Article", err)
		return
	}

	// return that we have successfully uploaded our file!
	http.Redirect(w, r, "/", http.StatusSeeOther)
	logger.Info("Successfully Uploaded File")
}

// SaveFile
func (h *Handler) SaveFile(w http.ResponseWriter, r *http.Request) {
	md := []byte(r.FormValue("textEditArea"))
	artclPath := r.FormValue("articlePath")
	err := os.WriteFile("articles/"+artclPath, md, 0644)
	logger := service.GetLogger()
	if err != nil {
		logger.Error("MD file write error: ", err, artclPath)
		return
	} else {
		logger.Info("Successfully Edited File")
	}

	err = h.services.ArticleService.UpdateArticle(artclPath)
	if err != nil {
		logger.Error(err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
