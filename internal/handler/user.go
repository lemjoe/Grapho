package handler

import (
	"html/template"

	"net/http"

	"github.com/lemjoe/Grapho/internal/models"
	"github.com/lemjoe/Grapho/internal/service"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) ChangeTheme(w http.ResponseWriter, r *http.Request) {
	logger := service.GetLogger()
	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	theme := "light"
	curSettings := curUser.Settings
	if curSettings["theme"] == "light" {
		theme = "dark"
	}
	curSettings["theme"] = theme
	err := h.services.UserService.ChangeUserSettings(curUser.Id, curSettings)
	if err != nil {
		logger.Error("Unable to change current theme: ", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// User settings
func (h *Handler) Settings(w http.ResponseWriter, r *http.Request) {
	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	if curUser.FullName == "Guest" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	lang := curUser.Settings["language"]
	translation := Localizer(localization, lang, h.bundle)

	t, err := template.ParseFiles("lib/templates/user-settings.html") //parse the html file homepage.html
	if err != nil {                                                   // if there is an error
		logger.Error("template parsing error: ", err) // log it
	}

	adminInterface := []byte("")
	if curUser.IsAdmin {
		adminInterface = adminPanel
	}

	tmpSettings := map[string]string{
		curUser.Settings["language"]: "selected",
		curUser.Settings["theme"]:    "selected",
	}

	UserSettingsPageVars := models.PageVariables{ //store the date and time in a struct
		Theme:       curUser.Settings["theme"],
		AdminPanel:  template.HTML(adminInterface),
		Translation: translation,
		Settings:    tmpSettings,
	}
	err = t.Execute(w, UserSettingsPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                          // if there is an error
		logger.Error("template executing error: ", err) //log it
	}
}

func (h *Handler) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	logger := service.GetLogger()

	logger.Info("Trying to change user password")

	curUser := h.GetCurrentUser(w.Header().Get("userID"))

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		logger.Error("Unauthorized status code 401")
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	old_passwd := r.FormValue("old_password")
	new_passwd := r.FormValue("new_password")
	r_new_passwd := r.FormValue("r_new_password")

	err := bcrypt.CompareHashAndPassword([]byte(curUser.Password), []byte(old_passwd))
	if err != nil {
		logger.Error("invalid password: ", err)
	} else {
		if new_passwd != r_new_passwd {
			logger.Error("passwords must match!")
		} else {
			err := h.services.UserService.ChangeUserPassword(curUser.Id, new_passwd)
			if err != nil {
				logger.Info("unable to change password: ", err)
			} else {
				logger.Info("password was successfully changed")
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

func (h *Handler) ChangeSettings(w http.ResponseWriter, r *http.Request) {
	logger := service.GetLogger()
	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	theme := r.FormValue("theme")
	lang := r.FormValue("language")

	newSettings := curUser.Settings
	newSettings["theme"] = theme
	newSettings["language"] = lang
	logger.Info("Selected settings:", newSettings)

	err := h.services.UserService.ChangeUserSettings(curUser.Id, newSettings)
	if err != nil {
		logger.Error("Unable to change user settings: ", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
