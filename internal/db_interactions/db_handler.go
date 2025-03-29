package db_interactions

import (
	"encoding/json"
	"errors"
	"os"
	"slices"

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

func AddRecipeJson(newRecipe models.Recipe) error {
	var recipeCollection, err = ReadRecipeJson()
	if err != nil {
		return err
	}

	*recipeCollection = append(*recipeCollection, newRecipe)

	data, err := json.Marshal(recipeCollection)
	if err != nil {
		return errors.New("Error adding new recipe to list")
	}

	if err := os.WriteFile("internal/data/recipes.json", data, 0666); err != nil {
		return errors.New("Error adding new recipe to list")
	}

	return nil
}

func DeleteRecipeJson(recipe models.Recipe) error {
	var recipeCollection, err = ReadRecipeJson()
	if err != nil {
		return err
	}

	if recipe.Index > len(*recipeCollection) {
		return errors.New("recipe index out of range")
	}

  *recipeCollection = slices.Delete(*recipeCollection, recipe.Index, recipe.Index + 1)

  for _, r := range (*recipeCollection)[recipe.Index:] {
    r.Index -= 1
  }

	return nil
}
