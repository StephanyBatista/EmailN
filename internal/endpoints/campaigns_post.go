package endpoints

import (
	"emailn/internal/domain/campaign"
	"net/http"

	"github.com/go-chi/render"
)

func (h *Handler) CampaignPost(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var request campaign.NewCampaignRequest
	render.DecodeJSON(r.Body, &request)
	email := r.Context().Value("email").(string)
	request.CreatedBy = email
	id, err := h.CampaignService.Create(request)
	return map[string]string{"id": id}, 201, err
}
