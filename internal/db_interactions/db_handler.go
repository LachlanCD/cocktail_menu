package db_interactions

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

// initDB initializes the SQLite database and creates the table if it doesn't exist
func InitDB() *sql.DB {

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	createTables(db)

	return db
}

func ReadHomePageData(db *sql.DB) (*[]models.HomePageRecipes, error) {
	var recipeCollection []models.HomePageRecipes
	recipesMap := make(map[int]*models.HomePageRecipes)

	recipes, err := readHomeRecipes(db)
	if err != nil {
		return nil, err
	}

	for _, r := range recipes {
		recipesMap[r.Index] = r
	}

	if err := readHomeSpirits(db, recipesMap); err != nil {
		return nil, err
	}

	// Convert map to slice
	for _, r := range recipesMap {
		recipeCollection = append(recipeCollection, *r)
	}

	return &recipeCollection, nil
}

func ReadRecipe(db *sql.DB, recipe_id int) (*models.Recipe, error) {
	recipe, err := readRecipeByID(db, recipe_id)
	if err != nil {
		return nil, err
	}

	ingredients, err := readIngredients(db, recipe_id)
	if err != nil {
		return nil, err
	}
	recipe.Ingredients = ingredients

	instructions, err := readInstructions(db, recipe_id)
	if err != nil {
		return nil, err
	}
	recipe.Instructions = instructions

	spirits, err := readSpirits(db, recipe_id)
	if err != nil {
		return nil, err
	}
	recipe.Spirit = spirits

	return recipe, nil
}

func AddNewRecipe(db *sql.DB, recipe *models.Recipe) (int, error) {
  recipe_id, err := addRecipeToDB(db, recipe)
  if err != nil {
    return 0, err
  }
  return int(recipe_id), nil
}

func DeleteRecipe(db *sql.DB, recipe_id int) error {
  if err := deleteRecipeFromDB(db, recipe_id); err != nil {
    return err
  }
  return nil
}
