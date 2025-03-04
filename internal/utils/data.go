package utils

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

func home() ([]models.HomePageRecipes, error) {
	var recipes []models.HomePageRecipes
	var recipeCollection []models.RecipeCollection

	file, err := os.ReadFile("internal/data/recipes.json")
	if err != nil {
		return nil, errors.New("Could not read recipes")
	}

	if err := json.Unmarshal(file, &recipeCollection); err != nil {
		return nil, errors.New("Error parsing recipes")
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

	return recipes, nil
}
