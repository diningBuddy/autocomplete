package handler

import (
	"net/http"

	"github.com/skku/autocomplete/app/model"
)

func HealthCheck(rd *model.AutocompleteRedis, v *model.Version, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"versions": map[string]interface{}{
			"restaurant": v.Restaurant,
		},
	})
}
