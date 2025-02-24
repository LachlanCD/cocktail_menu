package handlers

import (
	"net/http"
)

func ServeHomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/templates/index.html")
}

