package utils

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

func GetHomePageData() (*[]models.HomePageRecipes, error) {
	var recipes []models.HomePageRecipes

  recipeCollection, err := ReadRecipeJson()
  if err != nil {
    return nil, err
  }

	for _, val := range *recipeCollection {
		var recipe = models.HomePageRecipes{
			Index:  val.Index,
			Name:   val.Name,
			Spirit: val.Spirit,
		}

		recipes = append(recipes, recipe)
	}

	return &recipes, nil
}

func GetRecipeData(index int) (*models.Recipe, error) {
  recipeCollections, err := ReadRecipeJson()
  if err != nil {
    return nil, err
  }

  if index < 1 || index > len(*recipeCollections) {
    return nil, errors.New("recipe index out of range")
  }

  return &(*recipeCollections)[index-1], nil
}
