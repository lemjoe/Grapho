package handler

import (
	"log"
	"net/http"
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
