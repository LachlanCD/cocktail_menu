package handlers

import (
	"math/rand"

	"github.com/lachlancd/cocktail_menu/internal/models"
)

func getRandomId(recipes *[]models.HomePageRecipes) int {
  randRecipe := (*recipes)[rand.Intn(len(*recipes))]
  return randRecipe.Index
}
