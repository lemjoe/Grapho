package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/lemjoe/md-blog/internal/models"
)

// SingUp
func (h *Handler) SingUp(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	translation := Localizer([]string{"homeButton"}, lang, h.bundle)

	t, err := template.ParseFiles("lib/templates/sing-up.html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	SingUpPageVars := models.PageVariables{ //store the date and time in a struct
		HomeButton: translation["homeButton"],
	}
	err = t.Execute(w, SingUpPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                    // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
func (h *Handler) SingIn(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	translation := Localizer([]string{"homeButton"}, lang, h.bundle)

	t, err := template.ParseFiles("lib/templates/sing-in.html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	SingInPageVars := models.PageVariables{ //store the date and time in a struct
		HomeButton: translation["homeButton"],
	}
	err = t.Execute(w, SingInPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                    // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
