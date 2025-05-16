package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/litetable/litetable-cli/internal/dir"
	"net/http"
)

func (h *handler) getFamilies(w http.ResponseWriter, r *http.Request) {
	families, err := dir.GetFamilies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get families",
		})
		return
	}

	fmt.Println("Families:", families)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(families); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to encode families",
		})
		return
	}
}
