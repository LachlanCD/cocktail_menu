package utils

import (
	"errors"

	"github.com/lachlancd/cocktail_menu/internal/models"
	"github.com/lachlancd/cocktail_menu/internal/db_interactions"
)

func GetHomePageData() (*[]models.HomePageRecipes, error) {
	var recipes []models.HomePageRecipes

  recipeCollection, err := db_interactions.ReadRecipeJson()
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
  recipeCollection, err := db_interactions.ReadRecipeJson()
  if err != nil {
    return nil, err
  }

  if index < 1 || index > len(*recipeCollection) {
    return nil, errors.New("recipe index out of range")
  }

  return &(*recipeCollection)[index-1], nil
}

func CreateNewRecipe(newRecipe models.Recipe) (error) {
  if err := db_interactions.AddRecipeJson(newRecipe); err != nil {
    return err
  }

  return nil
}
