package main

import (
	"io"
	"time"

	"github.com/gomarkdown/markdown"

	"html/template"
	"log"
	"net/http"
	"os"
)

type PageVariables struct {
	Md        string
	MDArticle template.HTML
}

var toTheTop = []byte("\n<a href=\"#top\"><i>back to top</i></a>")

func init() {
	http.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("lib"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	// http.Handle("/articles/", http.StripPrefix("/articles/", http.FileServer(http.Dir("articles"))))
}

func main() {
	http.HandleFunc("/show", ShowArticle)
	http.HandleFunc("/edit", Editor)
	http.HandleFunc("/delete", DeleteArticle)
	http.HandleFunc("/add", UploadArticle)
	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/", ArticleList)
	log.Print("Server is running on port 4007")
	log.Fatal(http.ListenAndServe(":4007", nil))
}

func ShowArticle(w http.ResponseWriter, r *http.Request) {

	artclPath := r.URL.Query().Get("md")
	md, err := os.ReadFile("articles/" + artclPath) // just pass the file name
	if err != nil {
		log.Print("MD file open error: ", err, artclPath)
	}
	// always normalize newlines!
	md = markdown.NormalizeNewlines(md)
	html := append(markdown.ToHTML(md, nil, nil), toTheTop[:]...)

	HomePageVars := PageVariables{ //store the date and time in a struct
		MDArticle: template.HTML(html),
	}

	t, err := template.ParseFiles("template.html") //parse the html file homepage.html
	if err != nil {                                // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func Editor(w http.ResponseWriter, r *http.Request) {

	artclPath := r.URL.Query().Get("md")
	md, err := os.ReadFile("articles/" + artclPath) // just pass the file name
	if err != nil {
		log.Print("MD file open error: ", err)
	}
	// always normalize newlines!
	md = markdown.NormalizeNewlines(md)
	html := markdown.ToHTML(md, nil, nil)

	HomePageVars := PageVariables{ //store the date and time in a struct
		Md:        string(md),
		MDArticle: template.HTML(html),
	}

	t, err := template.ParseFiles("markdown_editor.html") //parse the html file homepage.html
	if err != nil {                                       // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {

	artclPath := r.URL.Query().Get("md")
	err := os.Remove("articles/" + artclPath)
	if err != nil {
		log.Print("MD file delete error: ", err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("Successfully Deleted File")
}

func UploadArticle(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("upload.html") //parse the html file homepage.html
	if err != nil {                              // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, time.Now()) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	log.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}

	ftype := handler.Filename[len(handler.Filename)-3:]
	if ftype != ".md" {
		log.Println("File type must be text/markdown" + ftype)
		return
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern

	tempFile, err := os.Create("articles/" + handler.Filename)
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("Successfully Uploaded File")
}

func ArticleList(w http.ResponseWriter, r *http.Request) {

	f, err := os.Open("articles")
	if err != nil {
		log.Print("Articles directory open error: ", err)
		return
	}

	files, err := f.Readdir(0)
	if err != nil {
		log.Print("Articles read error: ", err)
		return
	}

	html := "<h2>List of available articles:</h2><ul>"

	for _, v := range files {
		if !v.IsDir() {
			html += "<li>" + "<a href='show?md=" + v.Name() + "'>" + v.Name() + "</a><i> (Last modification: " + v.ModTime().Format("2006-Jan-02") + ") </i><a href='edit?md=" + v.Name() + "'><i>edit</i></a> | <a href='delete?md=" + v.Name() + "'><i>delete</i></a></li>"
		}
	}

	html += "</ul>"

	HomePageVars := PageVariables{ //store the date and time in a struct
		MDArticle: template.HTML(html),
	}

	t, err := template.ParseFiles("template.html") //parse the html file homepage.html
	if err != nil {                                // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
