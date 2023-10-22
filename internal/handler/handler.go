package handler

import (
	"log"
	"net/http"

	"github.com/lemjoe/md-blog/internal/handler/deprecated"
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
	// http.HandleFunc("/show", ShowArticle)
	// http.HandleFunc("/edit", Editor)
	// http.HandleFunc("/delete", DeleteArticle)
	// http.HandleFunc("/add", UploadArticle)
	// http.HandleFunc("/upload", Upload)
	// http.HandleFunc("/download", DownloadArticle)
	// http.HandleFunc("/convert", MDConvert)
	// http.HandleFunc("/save", SaveFile)
	// http.HandleFunc("/singup", SingUp)
	// http.HandleFunc("/singin", SingIn)
	// http.HandleFunc("/", ArticleList)
	http.Handle("/", deprecated.Init(h.bundle, h.services).Router())
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("lib"))))

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	log.Print("Server is running on port ", port)
	return http.ListenAndServe(port, nil)
}
