package handlers

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/utils"
)

func (h *Handlers) GetRecipeHandler(w http.ResponseWriter, r *http.Request) {

	templ := template.Must(template.ParseFiles(
		"internal/templates/index.html",
		"internal/templates/nav.html", 
		"internal/templates/add_form.html",
		"internal/templates/recipe.html",
    ))

  id, err := strconv.Atoi(r.PathValue("id"))
  if err != nil {
    http.Error(w, "Id must be a number", http.StatusInternalServerError)
    return
  }

  recipe, err := utils.GetRecipeData(id, h.DB)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

	if r.Header.Get("HX-Request") != "" {
		err := templ.ExecuteTemplate(w, "content", recipe)
		if err != nil {
			http.Error(w, "Could not load recipe page", http.StatusInternalServerError)
		}
		return
	}

  err = templ.Execute(w, recipe)
	if err != nil {
		http.Error(w, "Could not load recipe page", http.StatusInternalServerError)
		return
	}

}
