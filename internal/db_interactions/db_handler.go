package db_interactions

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"slices"
  "fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

// Database file path
const dbPath = "data/recipes.db"

// initDB initializes the SQLite database and creates the table if it doesn't exist
func InitDB() *sql.DB {
	// Open (or create) the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Create the recipes table if it doesn't exist
	query := `CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		ingredients TEXT NOT NULL,
		instructions TEXT NOT NULL
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Database initialized successfully.")
	return db
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
