package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/lemjoe/md-blog/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) ChangeTheme(w http.ResponseWriter, r *http.Request) {
	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	theme := "light"
	curSettings := curUser.Settings
	if curSettings["theme"] == "light" {
		theme = "dark"
	}
	curSettings["theme"] = theme
	err := h.services.UserService.ChangeUserSettings(curUser.Id, curSettings)
	if err != nil {
		log.Println("Unable to change current theme: ", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// User settings
func (h *Handler) Settings(w http.ResponseWriter, r *http.Request) {
	curUser := h.GetCurrentUser(w.Header().Get("userID"))

	if curUser.FullName == "Guest" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	lang := curUser.Settings["language"]
	translation := Localizer([]string{"homeButton"}, lang, h.bundle)

	t, err := template.ParseFiles("lib/templates/user-settings.html") //parse the html file homepage.html
	if err != nil {                                                   // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	UserSettingsPageVars := models.PageVariables{ //store the date and time in a struct
		HomeButton: translation["homeButton"],
		Theme:      curUser.Settings["theme"],
	}
	err = t.Execute(w, UserSettingsPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                          // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func (h *Handler) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to change user password")

	curUser := h.GetCurrentUser(w.Header().Get("userID"))

	if curUser.FullName == "Guest" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	old_passwd := r.FormValue("old_password")
	new_passwd := r.FormValue("new_password")
	r_new_passwd := r.FormValue("r_new_password")

	err := bcrypt.CompareHashAndPassword([]byte(curUser.Password), []byte(old_passwd))
	if err != nil {
		log.Print("invalid password: ", err)
	} else {
		if new_passwd != r_new_passwd {
			log.Println("passwords must match!")
		} else {
			err := h.services.UserService.ChangeUserPassword(curUser.Id, new_passwd)
			if err != nil {
				log.Print("unable to change password: ", err)
			} else {
				log.Println("password was successfully changed")
				alert := AlertMessage{
					Title:       "Congratulations!",
					Description: "Your password was successfully changed",
				}
				h.SendAlert(w, r, alert)
				return
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
