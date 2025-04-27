package handlers

import (
	"errors"
	"net/http"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/models"
	"github.com/lachlancd/cocktail_menu/internal/utils"
)

func (h *Handlers) AddRecipeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unabbble to parse form", http.StatusInternalServerError)
		return
	}

  ingredients, err := getIngredients(r.Form["ingredient_name"], r.Form["ingredient_quantity"])
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  recipe := &models.NewRecipe {
    Name: r.FormValue("name"),
    Source: r.FormValue("source"),
    Ingredients: ingredients,
    Instructions: r.Form["instruction"],
    Spirit: r.Form["spirit"],
  }


  recipeId, err := utils.AddNewRecipe(h.DB, recipe) 
  if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
  }


	templ := template.Must(template.ParseFiles("internal/templates/add_new_response.html"))

  recipeData := models.HomePageRecipes {
    Index: recipeId,
    Name: recipe.Name,
    Spirit: recipe.Spirit,
  }

	err = templ.ExecuteTemplate(w, "add_new_response.html", recipeData)
	if err != nil {
		http.Error(w, "Could not load home page", http.StatusInternalServerError)
		return
	}
}

func getIngredients(names []string, quantities []string) ([]models.Ingredient, error) {
	var ingredients []models.Ingredient

	if len(names) != len(quantities) {
		return nil, errors.New("Ingredient names and quantities do not match")
	}

	for i := range names {
		ingredient := models.Ingredient{Name: names[i], Quantity: quantities[i]}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func AddIngredientFieldHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("internal/templates/add_ingredient.html"))
  err := templ.Execute(w, "base")
	if err != nil {
		http.Error(w, "Could not load ingredient", http.StatusInternalServerError)
		return
	}
}

func AddInstructionFieldHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("internal/templates/add_instruction.html"))
  err := templ.Execute(w, "base")
	if err != nil {
		http.Error(w, "Could not load instruction", http.StatusInternalServerError)
		return
	}
}

func AddSpiritFieldHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("internal/templates/add_spirit.html"))
  err := templ.Execute(w, "base")
	if err != nil {
		http.Error(w, "Could not load instruction", http.StatusInternalServerError)
		return
	}
}
