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

  enablePragma(db)

	createTables(db)

	return db
}

func ReadHomePageData(db *sql.DB, filterType string, filter string) (*[]models.HomePageRecipes, error) {
	var recipeCollection *[]models.HomePageRecipes
  var err error
  switch filterType {
    case "spirit":
      recipeCollection, err = readSpiritFilterHomeData(db, filter)
      if err != nil {
        return nil, err
      }
    default:
      recipeCollection, err = readDefaultHomeData(db)
      if err != nil {
        return nil, err
      }
  }
  return recipeCollection, nil
}

func readDefaultHomeData(db *sql.DB) (*[]models.HomePageRecipes, error) {
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

func readSpiritFilterHomeData(db *sql.DB, spirit string) (*[]models.HomePageRecipes, error) {
	var recipeCollection []models.HomePageRecipes
	recipesMap := make(map[int]*models.HomePageRecipes)

	recipes, err := filterSpirits(db, spirit)
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
	if err := deleteRecipeFromDB(db, recipe_id); err != nil {
		return err
	}
	return nil
}
