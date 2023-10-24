package handler

import (
	"log"
	"net/http"

	"github.com/lemjoe/md-blog/internal/service"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Handler struct {
	services *service.Service
	bundle   *i18n.Bundle
}

func NewHandler(services *service.Service, bundle *i18n.Bundle) *Handler {
	return &Handler{
		services: services,
		bundle:   bundle,
	}
}

func (h *Handler) Run(port string) error {
	http.HandleFunc("/show", h.ShowArticle)
	http.HandleFunc("/edit", h.Editor)
	http.HandleFunc("/delete", h.DeleteArticle)
	http.HandleFunc("/add", h.UploadArticle)
	http.HandleFunc("/upload", h.Upload)
	http.HandleFunc("/download", h.DownloadArticle)
	http.HandleFunc("/convert", h.MDConvert)
	http.HandleFunc("/save", h.SaveFile)
	http.HandleFunc("/singup", h.SingUp)
	http.HandleFunc("/singin", h.SingIn)
	http.HandleFunc("/", h.GetArticlesList)
	// http.Handle("/", deprecated.Init(h.bundle, h.services).Router())
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("lib"))))

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	log.Print("Server is running on port ", port)
	return http.ListenAndServe(port, nil)
}
