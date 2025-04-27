package db_interactions

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

// Database file path
const dbPath = "data/recipes.db"

func openDB() (*sql.DB, error) {
	// Open (or create) the SQLite database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// required for foreign key interactions
func enablePragma(db *sql.DB) {
	// Enable foreign key support
  _, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}
}

func createTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS Recipes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			source TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS Ingredients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			recipe_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			quantity TEXT NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS Instructions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			recipe_id INTEGER NOT NULL,
			step INTEGER NOT NULL,
			instruction TEXT NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS Base_Spirits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			recipe_id INTEGER NOT NULL,
			spirit TEXT NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// readRecipeByID retrieves a single recipe from the database that matches the given recipe ID.
//
// Input:
// - db: a pointer to an open SQL database connection.
// - recipe_id: the ID of the recipe to retrieve.
//
// Output:
// - A pointer to a models.Recipe containing the recipe's ID, name, and source.
// - An error if the query fails, no row is found, or scanning fails.
func readRecipeByID(db *sql.DB, recipe_id int) (*models.Recipe, error) {
	row := db.QueryRow("SELECT id, name, source FROM recipes WHERE id = ?", recipe_id)

	recipe := &models.Recipe{}

	if err := row.Scan(&recipe.Index, &recipe.Name, &recipe.Source); err != nil {
		return nil, err
	}

	return recipe, nil
}

// readIngredients retrieves all ingredients associated with a specific recipe ID.
//
// Input:
// - db: a pointer to an open SQL database connection.
// - recipe_id: the ID of the recipe whose ingredients should be fetched.
//
// Output:
// - A slice of pointers to models.Ingredient, each containing the ingredient's name and quantity.
// - An error if the query fails or row scanning encounters an issue.
func readIngredients(db *sql.DB, recipe_id int) ([]models.Ingredient, error) {
	rows, err := db.Query("SELECT name, quantity FROM ingredients WHERE recipe_id=?", recipe_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []models.Ingredient

	for rows.Next() {
		ingredient := models.Ingredient{}
		if err := rows.Scan(&ingredient.Name, &ingredient.Quantity); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

// readInstructions retrieves all instructions associated with a specific recipe ID,
// ordered by the step number.
//
// Input:
// - db: a pointer to an open SQL database connection.
// - recipe_id: the ID of the recipe whose instructions should be fetched.
//
// Output:
// - A slice of string pointers, each containing one instruction in order.
// - An error if the query fails or rows cannot be scanned.
func readInstructions(db *sql.DB, recipe_id int) ([]string, error) {
	rows, err := db.Query("SELECT instruction FROM instructions WHERE recipe_id=? ORDER BY step", recipe_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instructions []string

	for rows.Next() {
		var instruction string
		if err := rows.Scan(&instruction); err != nil {
			return nil, err
		}
		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

// readSpirits retrieves all spirits associated with a given recipe ID from the database.
//
// Input:
// - db: a pointer to an open SQL database connection.
// - recipe_id: the ID of the recipe whose associated spirits should be retrieved.
//
// Output:
// - A slice of string pointers, where each pointer represents a spirit name linked to the recipe.
// - An error if any issue occurs during the query or result scanning.
//
// It queries the 'spirit' table for rows where the 'recipe_id' matches the input,
// and returns all corresponding 'spirit' values as pointers.
func readSpirits(db *sql.DB, recipe_id int) ([]string, error) {
	rows, err := db.Query("SELECT spirit FROM Base_Spirits WHERE recipe_id=?", recipe_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spirits []string

	for rows.Next() {
		var spirit string
		if err := rows.Scan(&spirit); err != nil {
			return nil, err
		}
		spirits = append(spirits, spirit)
	}

	return spirits, nil
}

// readHomeRecipes retrieves all recipes from the database.
//
// Input:
// - db: a pointer to an open SQL database connection.
//
// Output:
// - A slice of pointers to models.HomePageRecipes, each containing the recipe's ID (Index) and name.
// - An error if the query fails or rows cannot be scanned.
func readHomeRecipes(db *sql.DB) ([]*models.HomePageRecipes, error) {
	rows, err := db.Query("SELECT id, name FROM recipes ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*models.HomePageRecipes

	for rows.Next() {
		recipe := models.HomePageRecipes{}
		if err := rows.Scan(&recipe.Index, &recipe.Name); err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}
	return recipes, nil
}

// readHomeSpirits retrieves the spirits associated with recipes from the "Base_Spirits" table
// and appends them to the corresponding recipes in the provided map.
//
// Input:
//   - db: a pointer to an open SQL database connection.
//   - recipesMap: a map where keys are recipe IDs and values are pointers to HomePageRecipes
//     that hold the spirits for each recipe.
//
// Output:
// - An error if the query fails, or if scanning rows encounters an issue. Returns nil if successful.
//
// This function queries the "Base_Spirits" table for each recipe ID and its associated spirit.
// It then checks if the recipe ID exists in the provided recipesMap. If it does, the spirit is appended
// to the recipe's Spirit slice.
func readHomeSpirits(db *sql.DB, recipesMap map[int]*models.HomePageRecipes) error {
	rows, err := db.Query("SELECT recipe_id, spirit FROM Base_Spirits")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe_id int
		var spirit string
		if err := rows.Scan(&recipe_id, &spirit); err != nil {
			return err
		}
		if recipe, exists := recipesMap[recipe_id]; exists {
			recipe.Spirit = append(recipe.Spirit, spirit)
		}
	}
	return nil
}

func filterSpirits(db *sql.DB, spirit string) ([]*models.HomePageRecipes, error) {
  rows, err := db.Query("select recipes.id, recipes.name from recipes inner join base_spirits on recipes.id=base_spirits.recipe_id and base_spirits.spirit='?'", spirit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
  
	var recipes []*models.HomePageRecipes

	for rows.Next() {
		recipe := models.HomePageRecipes{}
		if err := rows.Scan(&recipe.Index, &recipe.Name); err != nil {
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}
	return recipes, nil
}

// addRecipe adds a new recipe along with its ingredients and instructions to the database.
//
// Input:
// - db: a pointer to an open SQL database connection.
// - recipe: the recipe to be added, including its name, source, ingredients, and instructions.
//
// Output:
// - An error if the query fails or anything goes wrong during the insertion.
func addRecipeToDB(db *sql.DB, recipe *models.NewRecipe) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec("INSERT INTO recipes (name, source) VALUES (?, ?)", recipe.Name, recipe.Source)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get the ID of the newly inserted recipe.
	recipeID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Add instructions to the 'instructions' table.
	for step, instruction := range recipe.Instructions {
		_, err := tx.Exec("INSERT INTO instructions (recipe_id, step, instruction) VALUES (?, ?, ?)", recipeID, step+1, instruction)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add ingredients to the 'ingredients' table.
	for _, ingredient := range recipe.Ingredients {
		_, err := tx.Exec("INSERT INTO ingredients (recipe_id, name, quantity) VALUES (?, ?, ?)", recipeID, ingredient.Name, ingredient.Quantity)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add spirits to the 'base_spirits' table.
	for _, spirit := range recipe.Spirit {
		_, err := tx.Exec("INSERT INTO base_spirits (recipe_id, spirit) VALUES (?, ?)", recipeID, spirit)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(recipeID), nil
}

// deleteRecipe deletes a recipe from the database based on the provided recipe ID.
// The related ingredients and instructions will be automatically deleted due to
// the ON DELETE CASCADE foreign key constraints.
//
// Input:
// - db: a pointer to an open SQL database connection.
// - recipe_id: the ID of the recipe to be deleted.
//
// Output:
// - An error if the delete query fails.
func deleteRecipeFromDB(db *sql.DB, recipe_id int) error {
	_, err := db.Exec("DELETE FROM recipes WHERE id=?", recipe_id)
	if err != nil {
		return err
	}

	return nil
}
