package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/utils"
)

type Handlers struct {
  DB *sql.DB
}

func (h *Handlers) GetHomeHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles(
		"internal/templates/index.html",
		"internal/templates/nav.html",
		"internal/templates/add_form.html",
		"internal/templates/home.html"))

	recipes, err := utils.GetHomePageData(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") != "" {
    log.Println("htmx")
		err := templ.ExecuteTemplate(w, "content", recipes)
		if err != nil {
			http.Error(w, "Could not load home page", http.StatusInternalServerError)
		}
		return
	}

	err = templ.ExecuteTemplate(w, "base", recipes)
	if err != nil {
		http.Error(w, "Could not load home page", http.StatusInternalServerError)
		return
	}
}
