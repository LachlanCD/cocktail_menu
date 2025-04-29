package handlers

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/utils"
)

func (h *Handlers) EditRecipeFormHandler(w http.ResponseWriter, r *http.Request) {

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

	templ := template.Must(template.ParseFiles(
		"internal/templates/index.html",
		"internal/templates/nav.html", 
		"internal/templates/edit_form.html",
    ))

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

