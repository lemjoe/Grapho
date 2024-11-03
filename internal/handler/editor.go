package handler

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/lemjoe/Grapho/internal/models"
	"github.com/lemjoe/Grapho/internal/service"
)

func (h *Handler) Editor(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
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
	} else if !curUser.IsWriter {
		logger.Error("Wrong permissions. You are unable to write articles: status code 403")
		h.SendCode(w, r, statusCodes[http.StatusForbidden])
		return
	}

	lang := curUser.Settings["language"]
	translation := Localizer(localization, lang, h.bundle)

	md, err := os.ReadFile("articles/" + artclPath) // just pass the file name
	if err != nil {
		logger.Error("MD file open error: ", err)
	}
	html := MdToHTML(md, true)

	HomePageVars := models.PageVariables{ //store the date and time in a struct
		Md:          string(md),
		MDArticle:   template.HTML(html),
		Path:        artclPath,
		UserName:    curUser.FullName,
		Theme:       curUser.Settings["theme"],
		Translation: translation,
		Title:       translation["titleEdit"],
	}
	t, err := template.ParseFiles("lib/templates/editor.html") //parse the html file homepage.html
	if err != nil {                                            // if there is an error
		logger.Error("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		logger.Error("template executing error: ", err) //log it
	}
}

func (h *Handler) MDConvert(w http.ResponseWriter, r *http.Request) {

	md, _ := io.ReadAll(r.Body)
	rg := regexp.MustCompile(`(?:[\t ]*(?:\r?\n|\r))`)
	html := MdToHTML(md, true)
	str := string(html)
	result := rg.ReplaceAllString(str, "")
	html = []byte(result)
	w.Header().Set("Content-Type", "application/json")
	responseJSON := map[string]string{"msg": string(html)}
	json.NewEncoder(w).Encode(responseJSON)
}
