package campaign

type CampaignResponse struct {
	ID                   string
	Name                 string
	Content              string
	Status               string
	AmountOfEmailsToSend int
	CreatedBy            string
}
