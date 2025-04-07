package handlers

import (
	"errors"
	"fmt"
	"net/http"

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


  if err := utils.AddNewRecipe(h.DB, recipe); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  fmt.Println(recipe)
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
	// Return a fresh set of ingredient input fields
	html := `
    <div>
      <input name="ingredient_name" placeholder="Name">
      <input name="ingredient_quantity" placeholder="Quantity">
      <button type="button" onclick="this.parentElement.remove()">Remove</button>
    </div>`
	fmt.Fprint(w, html)
}

func AddInstructionFieldHandler(w http.ResponseWriter, r *http.Request) {
	html := `
    <div>
      <input name="instruction" placeholder="Step">
      <button type="button" onclick="this.parentElement.remove()">Remove</button>
    </div>`
	fmt.Fprint(w, html)
}

func AddSpiritFieldHandler(w http.ResponseWriter, r *http.Request) {
	html := `
    <div>
      <input name="spirit" placeholder="Name">
      <button type="button" onclick="this.parentElement.remove()">Remove</button>
    </div>`
	fmt.Fprint(w, html)
}
