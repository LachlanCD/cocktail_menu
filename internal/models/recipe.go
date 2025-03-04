package models

type Ingredient struct {
	Name     string
	Quantity int
}

type Recipe struct {
	Name         string
	Ingredients  []Ingredient
	Instructions []string
	Image        string
	Source       string
	Spirit       []string
}

type RecipeCollection struct {
	Index int
	Name  string
	Types []Recipe
}

type HomePageRecipes struct {
	Index  int
	Name   string
	Spirit []string
}
