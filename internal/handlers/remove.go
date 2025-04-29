package handlers

import (
	"net/http"
	"strconv"

	"github.com/lachlancd/cocktail_menu/internal/utils"
)

func (h *Handlers) RemoveRecipeHandler(w http.ResponseWriter, r *http.Request) {

  id, err := strconv.Atoi(r.PathValue("id"))
  if err != nil {
    http.Error(w, "Id must be a number", http.StatusInternalServerError)
    return
  }

  err = utils.DeleteRecipe(h.DB, id)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  http.Redirect(w, r, "/", http.StatusFound) 
}
