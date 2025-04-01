package main

import (
	"log"
	"net/http"

	"github.com/lachlancd/cocktail_menu/internal/db_interactions"
	"github.com/lachlancd/cocktail_menu/internal/handlers"
)

func main() {

  db := db_interactions.InitDB()

  defer db.Close()

  h := &handlers.Handlers{DB: db}

	// Serve static assets (CSS, JS, images)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", h.GetHomeHandler)
	http.HandleFunc("/recipe/{id}", h.GetRecipeHandler)

	log.Fatal(http.ListenAndServe(":6969", nil))
}
