package main

import (
	"log"
	"net/http"

	"github.com/lachlancd/cocktail_menu/internal/handlers"
)

func main() {

	// Serve static assets (CSS, JS, images)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", handlers.GetHomeHandler)

	log.Fatal(http.ListenAndServe(":6969", nil))
}
