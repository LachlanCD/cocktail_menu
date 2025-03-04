package handlers

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/utils"
)

func GetRecipeHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("internal/templates/index.html", "internal/templates/nav.html", "internal/templates/card.html"))

  id, err := strconv.Atoi(r.PathValue("id"))
  if err != nil {
    http.Error(w, "Id must be a number", http.StatusInternalServerError)
  }

  recipe, err := utils.GetRecipeData(id)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  err = templ.Execute(w, recipe)
	if err != nil {
		http.Error(w, "Could not load recipe page", http.StatusInternalServerError)
		return
	}

}
