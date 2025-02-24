package handlers

import (
	"encoding/json"
  "os"
	"net/http"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

func GetRecipesHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("internal/data/recipes.json")
	if err != nil {
		http.Error(w, "Could not read recipes", http.StatusInternalServerError)
		return
	}

	var recipes []models.RecipeCollection
	if err := json.Unmarshal(file, &recipes); err != nil {
		http.Error(w, "Error parsing recipes", http.StatusInternalServerError)
		return
	}

	// Generate recipe HTML dynamically
	for _, recipe := range recipes {
    w.Write([]byte("<div><h3>" + recipe.Name + "</h3><p>" + recipe.Types[0].Name + "</p></div>"))
	}
}

