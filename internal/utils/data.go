package utils

import (
	"database/sql"

	"github.com/lachlancd/cocktail_menu/internal/db_interactions"
	"github.com/lachlancd/cocktail_menu/internal/models"
)

func GetHomePageData(db *sql.DB) (*[]models.HomePageRecipes, error) {
  recipes, err := db_interactions.ReadHomePageData(db)
  if err != nil {
    return nil, err
  }

	return recipes, nil
}

func GetRecipeData(index int, db *sql.DB) (*models.Recipe, error) {
  recipe, err := db_interactions.ReadRecipe(db, index)
  if err != nil {
    return nil, err
  }

  return recipe, nil
}

func CreateNewRecipe(newRecipe models.Recipe) (error) {
  if err := db_interactions.AddRecipeJson(newRecipe); err != nil {
    return err
  }

  return nil
}
