package deprecated

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

// SingUp
func (h *handler) SingUp(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("lib/templates/sing-up.html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, time.Now()) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
func (h *handler) SingIn(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("lib/templates/sing-in.html") //parse the html file homepage.html
	if err != nil {                                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, time.Now()) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
