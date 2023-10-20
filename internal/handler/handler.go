package handler

import (
	"log"
	"net/http"

	"github.com/lemjoe/md-blog/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
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
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("lib"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	log.Print("Server is running on port ", port)
	return http.ListenAndServe(port, nil)
}
