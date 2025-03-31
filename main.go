package main

import (
	"log"
	"net/http"

	"github.com/lachlancd/cocktail_menu/internal/db_interactions"
	"github.com/lachlancd/cocktail_menu/internal/handlers"
)

func main() {

  db_interactions.InitDB()

	// Serve static assets (CSS, JS, images)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", handlers.GetHomeHandler)
	http.HandleFunc("/recipe/{id}", handlers.GetRecipeHandler)

	log.Fatal(http.ListenAndServe(":6969", nil))
}
