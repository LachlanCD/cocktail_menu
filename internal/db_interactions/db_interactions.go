package db_interactions

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

// Database file path
const dbPath = "data/recipes.db"

type Execer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

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

func readDB(db *sql.DB, query string, input *string) (*sql.Rows, error) {
	if input != nil {
		return db.Query(query, input)
	}
	return db.Query(query)
}

func processHomeRecipies(db *sql.DB, query string, input *string) ([]*models.HomePageRecipes, error) {
	rows, err := readDB(db, query, input)
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

// readHomeRecipes retrieves all recipes from the database.
//
// Input:
// - db: a pointer to an open SQL database connection.
//
// Output:
// - A slice of pointers to models.HomePageRecipes, each containing the recipe's ID (Index) and name.
// - An error if the query fails or rows cannot be scanned.
func readHomeRecipes(db *sql.DB) ([]*models.HomePageRecipes, error) {
	var inputPtr *string
	query := "SELECT id, name FROM recipes ORDER BY id"
	return processHomeRecipies(db, query, inputPtr)
}

func filterSpirits(db *sql.DB, spirit string) ([]*models.HomePageRecipes, error) {
	inputPtr := &spirit
	query := "SELECT recipes.id, recipes.name FROM recipes INNER JOIN base_spirits ON recipes.id=base_spirits.recipe_id AND base_spirits.spirit=?"
	return processHomeRecipies(db, query, inputPtr)
}

func searchRecipes(db *sql.DB, search string) ([]*models.HomePageRecipes, error) {
	searchPattern := search + "%"
	inputPtr := &searchPattern
	query := "SELECT id, name FROM recipes WHERE name LIKE ?"
	return processHomeRecipies(db, query, inputPtr)
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

func readUniqueSpirits(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT spirit FROM base_spirits")
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
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO recipes (name, source) VALUES (?, ?)", recipe.Name, recipe.Source)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get the ID of the newly inserted recipe.
	last_id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	recipe_id := int(last_id)

  if err := addInstructions(tx, recipe_id, recipe.Instructions); err != nil {
    return 0, err
  }

	if err := addIngredients(tx, recipe_id, recipe.Ingredients); err != nil {
    return 0, err
  }

	if err := addSpirits(tx, recipe_id, recipe.Spirit); err != nil {
    return 0, err
  }

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return recipe_id, nil
}

func addInstructions(tx *sql.Tx, recipe_id int, instructions []string) error {
  stmt, err := tx.Prepare("INSERT INTO instructions (recipe_id, step, instruction) VALUES (?, ?, ?)")
  if err != nil {
    return err
  }
  defer stmt.Close()

	// Add instructions to the 'instructions' table.
	for step, instruction := range instructions {
		_, err := stmt.Exec(recipe_id, step+1, instruction)
		if err != nil {
			return err
		}
	}
	return nil
}

func addIngredients(tx *sql.Tx, recipe_id int, ingredients []models.Ingredient) error {
  stmt, err := tx.Prepare("INSERT INTO ingredients (recipe_id, name, quantity) VALUES (?, ?, ?)")
  if err != nil {
    return err
  }
  defer stmt.Close()

	// Add ingredients to the 'ingredients' table.
	for _, ingredient := range ingredients {
		_, err := stmt.Exec(recipe_id, ingredient.Name, ingredient.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func addSpirits(tx *sql.Tx, recipe_id int, spirits []string) error {
  stmt, err := tx.Prepare("INSERT INTO base_spirits (recipe_id, spirit) VALUES (?, ?)")
  if err != nil {
    return err
  }
  defer stmt.Close()

	// Add spirits to the 'base_spirits' table.
	for _, spirit := range spirits {
		_, err := stmt.Exec(recipe_id, spirit)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteFromDB(exec Execer, query string, id int) error {
	_, err := exec.Exec(query, id)
	return err
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
	query := "DELETE FROM recipes WHERE id=?"
	return deleteFromDB(db, query, recipe_id)
}

func editRecipe(db *sql.DB, newRecipe *models.Recipe, recipe bool, ingredients bool, instructions bool, spirits bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

  if recipe {
    if err := updateRecipe(tx, newRecipe.Index, newRecipe.Name, newRecipe.Source); err != nil {
      return err
    }
  }

  if ingredients {
    if err := updateIngredients(tx, newRecipe.Index, newRecipe.Ingredients); err != nil {
      return err
    }
  }

  if instructions {
    if err := updateInstructions(tx, newRecipe.Index, newRecipe.Instructions); err != nil {
      return err
    }
  }

  if spirits {
    if err := updateSpirits(tx, newRecipe.Index, newRecipe.Spirit); err != nil {
      return err
    }
  }

	return tx.Commit()
}

func updateRecipe(tx *sql.Tx, recipe_id int, name string, source string) error {
	query := "update recipes set name=? source=? where id=?"
	_, err := tx.Exec(query, name, source, recipe_id)
	return err
}

func updateIngredients(tx *sql.Tx, recipe_id int, ingredients []models.Ingredient) error {
	query := "DELETE FROM ingredients WHERE recipe_id=?"
  if err := deleteFromDB(tx, query, recipe_id); err != nil {
    return err
  }
  if err := addIngredients(tx, recipe_id, ingredients); err != nil {
    return err
  }
	return nil
}

func updateInstructions(tx *sql.Tx, recipe_id int, instructions []string) error {
	query := "DELETE FROM instructions WHERE recipe_id=?"
  if err := deleteFromDB(tx, query, recipe_id); err != nil {
    return err
  }
  if err := addInstructions(tx, recipe_id, instructions); err != nil {
    return err
  }
  return nil
}

func updateSpirits(tx *sql.Tx, recipe_id int, spirits []string) error {
	query := "DELETE FROM base_spirits WHERE recipe_id=?"
  if err := deleteFromDB(tx, query, recipe_id); err != nil {
    return err
  }
  if err := addSpirits(tx, recipe_id, spirits); err != nil {
    return err
  }
	return nil
}
