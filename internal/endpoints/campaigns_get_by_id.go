package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GET /campaigns/{id}
func (h *Handler) CampaignGetById(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	id := chi.URLParam(r, "id")
	campaign, err := h.CampaignService.GetBy(id)
	return campaign, 200, err
}
