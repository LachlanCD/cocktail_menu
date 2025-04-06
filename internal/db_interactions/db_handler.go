package db_interactions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"slices"

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

	*recipeCollection = slices.Delete(*recipeCollection, recipe.Index, recipe.Index+1)

	for _, r := range (*recipeCollection)[recipe.Index:] {
		r.Index -= 1
	}

	return nil
}
