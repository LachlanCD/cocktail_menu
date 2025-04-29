package utils

import (
	"database/sql"

	"github.com/lachlancd/cocktail_menu/internal/db_interactions"
	"github.com/lachlancd/cocktail_menu/internal/models"
)

func GetHomePageData(db *sql.DB) (*[]models.HomePageRecipes, error) {
	recipes, err := db_interactions.ReadHomePageData(db, "", "")
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func GetSpiritFilterData(db *sql.DB, spirit string) (*[]models.HomePageRecipes, error) {
	recipes, err := db_interactions.ReadHomePageData(db, "spirit", spirit)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func GetRecipeSearchData(db *sql.DB, search string) (*[]models.HomePageRecipes, error) {
	recipes, err := db_interactions.ReadHomePageData(db, "search", search)
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

func AddNewRecipe(db *sql.DB, recipe *models.NewRecipe) (int, error) {
	recipeId, err := db_interactions.AddNewRecipe(db, recipe)
	if err != nil {
		return 0, err
	}

	return recipeId, nil
}

func DeleteRecipe(db *sql.DB, recipe_id int) error {
	if err := db_interactions.DeleteRecipe(db, recipe_id); err != nil {
		return err
	}

	return nil
}

func GetUniqueSpirits(db *sql.DB) ([]string, error) {
  return db_interactions.ReadSpirits(db) 
}

func EditRecipe(db *sql.DB, newRecipe *models.Recipe) error {
  return db_interactions.EditRecipe(db, newRecipe) 
}
