package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/lemjoe/Grapho/internal/models"
	"github.com/lemjoe/Grapho/internal/service"
)

type StatusCode struct {
	Code        string
	Title       string
	Description string
}

type AlertMessage struct {
	Title       string
	Description string
}

func (h *Handler) SendCode(w http.ResponseWriter, r *http.Request, status StatusCode) {
	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	lang := curUser.Settings["language"]
	translation := Localizer(localization, lang, h.bundle)

	intCode, err := strconv.Atoi(status.Code)
	if err != nil { // if there is an error
		logger.Error("status conversion error: ", err) // log it
		fmt.Fprintln(w, "status conversion error: ", err)
		return
	}
	w.WriteHeader(intCode)

	t, err := template.ParseFiles("lib/templates/status.html") //parse the html file homepage.html
	if err != nil {                                            // if there is an error
		logger.Error("template parsing error: ", err) // log it
		fmt.Fprintln(w, "template parsing error: ", err)
		return
	}
	StatusPageVars := models.PageVariables{ //store the date and time in a struct
		Title:        status.Code + " - " + status.Title,
		BodyLoudText: status.Title,
		BodyText:     status.Description,
		Theme:        curUser.Settings["theme"],
		Translation:  translation,
	}
	err = t.Execute(w, StatusPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                    // if there is an error
		logger.Error("template executing error: ", err) //log it
		fmt.Fprintln(w, "template executing error: ", err)
		return
	}
}

func (h *Handler) SendAlert(w http.ResponseWriter, r *http.Request, alert AlertMessage) {
	curUser := h.GetCurrentUser(w.Header().Get("userID"))
	logger := service.GetLogger()

	lang := curUser.Settings["language"]
	translation := Localizer(localization, lang, h.bundle)

	t, err := template.ParseFiles("lib/templates/status.html") //parse the html file homepage.html
	if err != nil {                                            // if there is an error
		logger.Error("template parsing error: ", err) // log it
		fmt.Fprintln(w, "template parsing error: ", err)
		return
	}
	AlertPageVars := models.PageVariables{ //store the date and time in a struct
		Title:        alert.Title,
		BodyLoudText: alert.Title,
		BodyText:     alert.Description,
		Theme:        curUser.Settings["theme"],
		Translation:  translation,
	}
	err = t.Execute(w, AlertPageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                   // if there is an error
		logger.Error("template executing error: ", err) //log it
		fmt.Fprintln(w, "template executing error: ", err)
		return
	}
}

func (h *Handler) PageNotFound(w http.ResponseWriter, r *http.Request) {
	// Send 404 if not found
	h.SendCode(w, r, statusCodes[http.StatusNotFound])
}
