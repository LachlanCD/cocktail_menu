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
			FOREIGN KEY (recipe_id) REFERENCES Recipes (id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS Instructions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			recipe_id INTEGER NOT NULL,
			step INTEGER NOT NULL,
			instruction TEXT NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES Recipes (id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS Base_Spirits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			recipe_id INTEGER NOT NULL,
			spirit TEXT NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES Recipes (id) ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// initDB initializes the SQLite database and creates the table if it doesn't exist
func InitDB() *sql.DB {

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	createTables(db)

	return db
}

func ReadHomePageData(db *sql.DB) (*[]models.Recipe, error) {
	var recipeCollection []models.Recipe
	recipesMap := make(map[int]*models.Recipe)

	recipes, err := readRecipes(db)
	if err != nil {
		return nil, err
	}

	for _, r := range recipes {
		recipesMap[r.Index] = r
	}

	err = readIngredients(db, recipesMap)
	if err != nil {
		return nil, err
	}

	err = readInstructions(db, recipesMap)
	if err != nil {
		return nil, err
	}

	err = readSpirits(db, recipesMap)
	if err != nil {
		return nil, err
	}

	// Convert map to slice
	for _, r := range recipesMap {
		recipeCollection = append(recipeCollection, *r)
	}

	return &recipeCollection, nil
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

func readRecipes(db *sql.DB) ([]*models.Recipe, error) {
	rows, err := db.Query("SELECT id, name, source FROM recipes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*models.Recipe

	for rows.Next() {
		r := &models.Recipe{}
		if err := rows.Scan(&r.Index, &r.Name, &r.Source); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, nil
}

func readIngredients(db *sql.DB, recipesMap map[int]*models.Recipe) error {
	rows, err := db.Query("SELECT recipe_id, name, quantity FROM ingredients")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var recipeID int
		ing := models.Ingredient{}
		if err := rows.Scan(&recipeID, &ing.Name, &ing.Quantity); err != nil {
			return err
		}
		if recipe, exists := recipesMap[recipeID]; exists {
			recipe.Ingredients = append(recipe.Ingredients, ing)
		}
	}

	return nil
}

func readInstructions(db *sql.DB, recipesMap map[int]*models.Recipe) error {
	rows, err := db.Query("SELECT recipe_id, instruction FROM instructions ORDER BY step")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var recipeID int
		var instruction string
		if err := rows.Scan(&recipeID, &instruction); err != nil {
			return err
		}
		if recipe, exists := recipesMap[recipeID]; exists {
			recipe.Instructions = append(recipe.Instructions, instruction)
		}
	}

	return nil
}

func readSpirits(db *sql.DB, recipesMap map[int]*models.Recipe) error {
	rows, err := db.Query("SELECT recipe_id, spirit FROM base_spirits")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var recipeID int
		var spirit string
		if err := rows.Scan(&recipeID, &spirit); err != nil {
			return err
		}
		if recipe, exists := recipesMap[recipeID]; exists {
			recipe.Spirit = append(recipe.Spirit, spirit)
		}
	}

	return nil
}
