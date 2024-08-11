package api

import (
	"net/http"
	"text/template"
)

type ErrorData struct {
	Msg1       string
	Msg2       string
	StatusCode int
}

// Parse and execute error html page depending on error type.
func ErrorHandler(w http.ResponseWriter, statusCode int, msg1, msg2 string) {
	Error := ErrorData{
		Msg1:       msg1,
		Msg2:       msg2,
		StatusCode: statusCode,
	}

	w.WriteHeader(Error.StatusCode)

	tmpl, err := template.ParseFS(TemplatesFs, "templates/error.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, Error); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
