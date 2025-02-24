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
}

type RecipeCollection struct {
	Index int
	Name  string
	Types []Recipe
}
