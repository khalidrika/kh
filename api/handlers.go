package api

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"asciiArt/asciiart"
)

var TemplatesFs embed.FS

type WebPageData struct {
	Text   string
	Banner string
	Art    string
	Fonts  []string
}

// Handle Home path and write the index.html template.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, 404, "Look like you're lost!", "The page you are looking for is not available!")
		return
	}

	if r.Method != http.MethodGet {
		ErrorHandler(w, 405, http.StatusText(http.StatusMethodNotAllowed), "Only GET method is allowed!")
		return
	}

	tmpl, err := template.ParseFS(TemplatesFs, "templates/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Something seems wrong, try again later!", "Internal Server Error!")
		return
	}

	data := WebPageData{
		Text:   "Hello World!",
		Banner: "standard",
		Art:    "",
	}
	data.ReadFonts()
	data.ReadUserFonts()
	if data.Fonts == nil {
		ErrorHandler(w, http.StatusInternalServerError, "Something seems wrong, try again later!", "Internal Server Error!")
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Something seems wrong, try again later!", "Internal Server Error!")
		return
	}
}

// Handle the asciiart path.
func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, 405, http.StatusText(http.StatusMethodNotAllowed), "Only POST method is allowed!")
		return
	}

	text := "\r\n" + r.FormValue("text")
	banner := r.FormValue("banner")

	if len(text) > 5004 {
		ErrorHandler(w, 400, "Your text is beyond our limits :}", "Bad Request!")
		return
	}

	if banner == "" {
		ErrorHandler(w, 400, "Make sure your input is correct!", "Bad Request!")
		return
	}

	art, err := asciiart.ASCIIArt(text, banner)
	if err != nil {
		ErrorHandler(w, 400, "Make sure your input is correct!", "Bad Request!")
		return
	}

	data := WebPageData{
		Text:   text,
		Banner: banner,
		Art:    art,
	}
	data.ReadFonts()
	data.ReadUserFonts()

	tmpl, err := template.ParseFS(TemplatesFs, "templates/index.html")
	if err != nil {
		ErrorHandler(w, 500, "Something seems wrong, try again later!", "Internal Server Error!")
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, 500, "Something seems wrong, try again later!", "Internal Server Error!")
		return
	}
}

func (d *WebPageData) ReadFonts() {
	BannersFS := asciiart.Banners
	entries, err := BannersFS.ReadDir("banners")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
			d.Fonts = append(d.Fonts, strings.TrimSuffix(entry.Name(), ".txt"))
		}
	}
}

func (d *WebPageData) ReadUserFonts() {
	files, err := os.ReadDir("banners")
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			d.Fonts = append(d.Fonts, strings.TrimSuffix(file.Name(), ".txt"))
		}
	}
}

func Style(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/css/" {
		ErrorHandler(w, 404, "Look like you're lost!", "The page you are looking for is not available!")
		return
	}
	http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/css"))).ServeHTTP(w, r)
}
