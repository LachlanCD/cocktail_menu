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
  recipeCollections, err := db_interactions.ReadRecipeJson()
  if err != nil {
    return nil, err
  }

  if index < 1 || index > len(*recipeCollections) {
    return nil, errors.New("recipe index out of range")
  }

  return &(*recipeCollections)[index-1], nil
}
