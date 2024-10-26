// Users management
package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/lemjoe/Grapho/internal/models"
	"github.com/lemjoe/Grapho/internal/service"
)

func (h *Handler) GetUsersList(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		logger.Error("Unauthorized status code 401")
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	// Send 403 wrong user
	if !curUser.IsAdmin {
		logger.Error("Wrong user. Action forbidden: status code 403")
		h.SendCode(w, r, statusCodes[http.StatusForbidden])
		return
	}

	theme := curUser.Settings["theme"]
	lang := curUser.Settings["language"]
	translation := Localizer(localization, lang, h.bundle)

	users, err := h.services.UserService.GetUsersList()
	if err != nil {
		logger.Error(err)
	}

	usersInfo := make([][7]string, len(users))

	for i, user := range users {
		logger.Info(user)
		usersInfo[i] = [7]string{
			strconv.Itoa(i + 1),
			user.UserName,
			user.FullName,
			user.Email,
			strconv.FormatBool(user.IsWriter),
			strconv.FormatBool(user.IsAdmin),
			user.Id,
		}
	}

	AdminPanelPageVars := models.PageVariables{ //store the date and time in a struct
		Title:       translation["titleAdmUsersList"],
		Translation: translation,
		UserName:    curUser.FullName,
		Theme:       theme,
		UsersInfo:   usersInfo,
	}

	t, err := template.ParseFiles("lib/templates/users.html") //parse the html file homepage.html
	if err != nil {                                           // if there is an error
		logger.Error("Template parsing error: ", err) // log it
	}
	err = t.Execute(w, AdminPanelPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                        // if there is an error
		logger.Error("Template executing error: ", err) //log it
	}

}

func (h *Handler) ManageUser(w http.ResponseWriter, r *http.Request) {

	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		logger.Error("Unauthorized status code 401")
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	// Send 403 wrong user
	if !curUser.IsAdmin {
		logger.Error("Wrong user. Action forbidden: status code 403")
		h.SendCode(w, r, statusCodes[http.StatusForbidden])
		return
	}

	lang := curUser.Settings["language"]
	translation := Localizer(localization, lang, h.bundle)

	usrId := r.URL.Query().Get("usr")
	managedUsr, err := h.services.UserService.GetUserById(usrId)
	if err != nil {
		logger.Error("Can't load user data: ", err, managedUsr)
	}

	mngUsrIsAdm := ""
	if managedUsr.IsAdmin {
		mngUsrIsAdm = "checked"
	}
	mngUsrIsWriter := ""
	if managedUsr.IsWriter {
		mngUsrIsWriter = "checked"
	}

	ManageUserPageVars := models.PageVariables{ //store the date and time in a struct
		Title:               translation["titleAdmManageUser"],
		UserName:            curUser.FullName,
		Theme:               curUser.Settings["theme"],
		ManagedUserName:     managedUsr.UserName,
		ManagedUserFullName: managedUsr.FullName,
		ManagedUserEmail:    managedUsr.Email,
		ManagedUserIsAdmin:  mngUsrIsAdm,
		ManagedUserIsWriter: mngUsrIsWriter,
		Translation:         translation,
	}

	t, err := template.ParseFiles("lib/templates/manage-user.html") //parse the html file homepage.html
	if err != nil {                                                 // if there is an error
		logger.Error("template parsing error: ", err) // log it
	}
	err = t.Execute(w, ManageUserPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                        // if there is an error
		logger.Error("template executing error: ", err) //log it
	}
}

func (h *Handler) ChangeUser(w http.ResponseWriter, r *http.Request) {
	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	logger.Info("Current user: " + curUser.FullName)

	// Send 401 if unauthorized
	if curUser.UserName == "guest" {
		logger.Error("Unauthorized status code 401")
		h.SendCode(w, r, statusCodes[http.StatusUnauthorized])
		return
	}

	// Send 403 wrong user
	if !curUser.IsAdmin {
		logger.Error("Wrong user. Action forbidden: status code 403")
		h.SendCode(w, r, statusCodes[http.StatusForbidden])
		return
	}

	user_name := r.FormValue("user_name")
	full_name := r.FormValue("full_name")
	e_mail := r.FormValue("e_mail")
	is_admin := false
	if r.FormValue("is_admin") == "on" {
		is_admin = true
	}
	is_writer := false
	if r.FormValue("is_writer") == "on" {
		is_writer = true
	}

	managedUser, err := h.services.UserService.GetUserByName(user_name)
	if err != nil {
		logger.Error("Can't load user data: ", err)
	}

	err = h.services.UserService.UpdateUserData(managedUser.Id, full_name, e_mail, is_admin, is_writer)
	if err != nil {
		logger.Error("Something went wrong: ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
