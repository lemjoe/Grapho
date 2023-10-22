package deprecated

import (
	"html/template"
	"log"
	"net/http"
	"os"

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
