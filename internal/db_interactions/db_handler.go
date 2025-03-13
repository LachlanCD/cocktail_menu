package db_interactions

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

func ReadRecipeJson() (*[]models.Recipe, error) {
	var recipeCollection []models.Recipe

	file, err := os.ReadFile("internal/data/recipes.json")
	if err != nil {
		return nil, errors.New("Could not read recipes")
	}

	if err := json.Unmarshal(file, &recipeCollection); err != nil {
		return nil, errors.New("Error parsing recipes")
	}

  return &recipeCollection, nil
}
