package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

func GetRecipesHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("internal/templates/index.html", "internal/templates/nav.html", "internal/templates/card.html"))

	file, err := os.ReadFile("internal/data/recipes.json")
	if err != nil {
		http.Error(w, "Could not read recipes", http.StatusInternalServerError)
		return
	}

	var recipes []models.HomePageRecipes
	var recipeCollection []models.RecipeCollection
	if err := json.Unmarshal(file, &recipeCollection); err != nil {
		http.Error(w, "Error parsing recipes", http.StatusInternalServerError)
		return
	}

	for _, val := range recipeCollection {
		var recipe = models.HomePageRecipes{
			Index:  val.Index,
			Name:   val.Name,
			Spirit: val.Types[0].Spirit,
			Colour: "gray",
		}

		recipes = append(recipes, recipe)
	}

  err = templ.Execute(w, recipes)
	if err != nil {
		http.Error(w, "Could not load home page", http.StatusInternalServerError)
		return
	}
}

