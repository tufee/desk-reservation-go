package api

import (
	"html/template"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("web/templates/*.html"))
}

func Home(w http.ResponseWriter, r *http.Request) {
	x := "Olá"
	tmpl.ExecuteTemplate(w, "base.html", x)
}
