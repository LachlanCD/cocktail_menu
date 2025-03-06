package models

type Ingredient struct {
	Name     string
	Quantity string
}

type Recipe struct {
  Index int
	Name         string
	Ingredients  []Ingredient
	Instructions []string
	Image        string
	Source       string
	Spirit       []string
}

type HomePageRecipes struct {
	Index  int
	Name   string
	Spirit []string
}
