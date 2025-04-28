package handlers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/models"
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
		"internal/templates/dropdown.html",
		"internal/templates/searchbar.html",
		"internal/templates/home.html"))

	recipes, err := utils.GetHomePageData(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

  spirits, err := utils.GetUniqueSpirits(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

  data := &models.HomePageData{
    Spirits: spirits,
    Recipes: *recipes,
  }

	if r.Header.Get("HX-Request") != "" {
		err := templ.ExecuteTemplate(w, "content", data)
		if err != nil {
			http.Error(w, "Could not load home page", http.StatusInternalServerError)
		}
		return
	}

	err = templ.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Could not load home page", http.StatusInternalServerError)
		return
	}
}
