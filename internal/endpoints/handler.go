package endpoints

import "emailn/internal/domain/campaign"

type Handler struct {
	CampaignService campaign.Service
}
