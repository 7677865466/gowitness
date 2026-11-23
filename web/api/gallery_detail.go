package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sensepost/gowitness/pkg/log"
	"github.com/sensepost/gowitness/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DetailHandler returns the detail for a screenshot
//
//	@Summary		Results detail
//	@Description	Get details for a result.
//	@Tags			Results
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"The screenshot ID to load."
//	@Success		200	{object}	models.Result
//	@Failure		404	{string}	string	"result not found"
//	@Router			/results/detail/{id} [get]
func (h *ApiHandler) DetailHandler(w http.ResponseWriter, r *http.Request) {
	var response = &models.Result{}

	if err := h.DB.Model(&models.Result{}).
		Preload(clause.Associations).
		Preload("TLS.SanList").
		First(&response, chi.URLParam(r, "id")).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "result not found", http.StatusNotFound)
			return
		}

		log.Error("could not get detail for id", "err", err)
		http.Error(w, "could not get result detail", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
