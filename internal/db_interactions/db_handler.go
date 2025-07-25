package db_interactions

import (
	"database/sql"
	"log"
	"reflect"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

// initDB initializes the SQLite database and creates the table if it doesn't exist
func InitDB() *sql.DB {

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}

	enablePragma(db)

	createTables(db)

	return db
}

func ReadSpirits(db *sql.DB) ([]string, error) {
	return readUniqueSpirits(db)
}

func ReadHomePageData(db *sql.DB, filterType string, filter string) (*[]models.HomePageRecipes, error) {
	var recipeCollection []models.HomePageRecipes
	var recipes []*models.HomePageRecipes
	var err error
	recipesMap := make(map[int]*models.HomePageRecipes)

	switch filterType {
	case "spirit":
		recipes, err = filterSpirits(db, filter)
	case "search":
		recipes, err = searchRecipes(db, filter)
	default:
		recipes, err = readHomeRecipes(db)
	}

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

func AddNewRecipe(db *sql.DB, recipe *models.NewRecipe) (int, error) {
	recipeId, err := addRecipeToDB(db, recipe)
	if err != nil {
		return 0, err
	}
	return recipeId, err
}

func DeleteRecipe(db *sql.DB, recipe_id int) error {
	return deleteRecipeFromDB(db, recipe_id)
}

func EditRecipe(db *sql.DB, newRecipe *models.Recipe) error {
	oldRecipe, err := ReadRecipe(db, newRecipe.Index)
	if err != nil {
		return err
	}

	recipe, ingredients, instructions, spirits := false, false, false, false
	if reflect.DeepEqual(oldRecipe, newRecipe) {
		return nil
	}

	if oldRecipe.Name != newRecipe.Name || oldRecipe.Source != newRecipe.Source {
		recipe = true
	}
	if !reflect.DeepEqual(oldRecipe.Ingredients, newRecipe.Ingredients) {
		ingredients = true
	}
	if !reflect.DeepEqual(oldRecipe.Instructions, newRecipe.Instructions) {
		instructions = true
	}
	if !reflect.DeepEqual(oldRecipe.Spirit, newRecipe.Spirit) {
		spirits = true
	}

	return editRecipe(db, newRecipe, recipe, ingredients, instructions, spirits)
}
