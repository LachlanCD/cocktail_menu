package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lachlancd/cocktail_menu/internal/db_interactions"
	"github.com/lachlancd/cocktail_menu/internal/handlers"
)

func main() {

	db := db_interactions.InitDB()

	defer db.Close()

	fmt.Println("initialised")

	h := &handlers.Handlers{DB: db}

	// Serve static assets (CSS, JS, images)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", h.GetHomeHandler)
	http.HandleFunc("/spirit-filter/{spirit}", h.GetFilteredSpiritsHandler)
	http.HandleFunc("/search", h.GetSearchResultsHandler)
	http.HandleFunc("/recipe/{id}", h.GetRecipeHandler)
	http.HandleFunc("/add-recipe", h.AddRecipeHandler)
  http.HandleFunc("/remove-recipe/{id}", h.RemoveRecipeHandler)
  http.HandleFunc("/edit-recipe-form/{id}", h.EditRecipeFormHandler)
  http.HandleFunc("/edit-recipe/{id}", h.EditRecipeResponseHandler)
	http.HandleFunc("/add-ingredient-field", handlers.AddIngredientFieldHandler)
	http.HandleFunc("/add-instruction-field", handlers.AddInstructionFieldHandler)
	http.HandleFunc("/add-spirit-field", handlers.AddSpiritFieldHandler)

	fmt.Println("running")

	log.Fatal(http.ListenAndServe(":4000", nil))
}
