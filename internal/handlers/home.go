package handlers

import (
	"net/http"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/utils"
)

func GetHomeHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("internal/templates/index.html", "internal/templates/nav.html", "internal/templates/card.html"))

  recipes, err := utils.GetHomePageData()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  err = templ.Execute(w, recipes)
	if err != nil {
		http.Error(w, "Could not load home page", http.StatusInternalServerError)
		return
	}
}

