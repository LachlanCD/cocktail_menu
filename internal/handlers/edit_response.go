package handlers

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/models"
	"github.com/lachlancd/cocktail_menu/internal/utils"
)

func (h *Handlers) EditRecipeResponseHandler (w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unabbble to parse form", http.StatusInternalServerError)
		return
	}

  ingredients, err := getIngredients(r.Form["ingredient_name"], r.Form["ingredient_quantity"])
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  id, err := strconv.Atoi(r.PathValue("id"))
  if err != nil {
    http.Error(w, "Id must be a number", http.StatusInternalServerError)
    return
  }

  recipe := &models.Recipe {
    Index: id, 
    Name: r.FormValue("name"),
    Source: r.FormValue("source"),
    Ingredients: ingredients,
    Instructions: r.Form["instruction"],
    Spirit: r.Form["spirit"],
  }

  err = utils.EditRecipe(h.DB, recipe) 
  if err != nil {
  	http.Error(w, err.Error(), http.StatusInternalServerError)
  }

	templ := template.Must(template.ParseFiles(
		"internal/templates/index.html",
		"internal/templates/nav.html", 
		"internal/templates/add_form.html",
		"internal/templates/recipe.html",
		"internal/templates/remove_button.html",
		"internal/templates/edit_button.html",
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

