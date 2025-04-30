package handlers

import (
	"net/http"
	"text/template"

	"github.com/lachlancd/cocktail_menu/internal/models"
	"github.com/lachlancd/cocktail_menu/internal/utils"
)

func (h *Handlers) GetSearchResultsHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles(
		"internal/templates/index.html",
		"internal/templates/nav.html",
		"internal/templates/add_form.html",
		"internal/templates/dropdown.html",
		"internal/templates/searchbar.html",
		"internal/templates/random_button.html",
		"internal/templates/home.html"))

	search := r.URL.Query().Get("search")

	recipes, err := utils.GetRecipeSearchData(h.DB, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

  randId := getRandomId(recipes)

	spirits, err := utils.GetUniqueSpirits(h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &models.HomePageData{
    RandomId: randId,
		Spirits: spirits,
		Recipes: *recipes,
	}

	if r.Header.Get("HX-Request") != "" {
		err := templ.ExecuteTemplate(w, "content", data)
		if err != nil {
			http.Error(w, "Could not load home page", http.StatusInternalServerError)
		}
		return
	}

	err = templ.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Could not load home page", http.StatusInternalServerError)
		return
	}
}
