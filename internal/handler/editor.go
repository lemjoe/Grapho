package handler

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/lemjoe/md-blog/internal/models"
)

func (h *Handler) Editor(w http.ResponseWriter, r *http.Request) {

	artclPath := r.URL.Query().Get("md")
	md, err := os.ReadFile("articles/" + artclPath) // just pass the file name
	if err != nil {
		log.Print("MD file open error: ", err)
	}
	html := MdToHTML(md)

	HomePageVars := models.PageVariables{ //store the date and time in a struct
		Md:        string(md),
		MDArticle: template.HTML(html),
		Path:      artclPath,
	}

	t, err := template.ParseFiles("lib/templates/editor.html") //parse the html file homepage.html
	if err != nil {                                            // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func (h *Handler) MDConvert(w http.ResponseWriter, r *http.Request) {

	md, _ := io.ReadAll(r.Body)
	rg := regexp.MustCompile(`(?:[\t ]*(?:\r?\n|\r))`)
	html := MdToHTML(md)
	str := string(html)
	result := rg.ReplaceAllString(str, "")
	html = []byte(result)
	w.Header().Set("Content-Type", "application/json")
	responseJSON := map[string]string{"msg": string(html)}
	json.NewEncoder(w).Encode(responseJSON)
}
