package endpoints

import (
	"emailn/internal/contract"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignsGetById_Campaign(t *testing.T) {
	setup()
	campaignId := "343"
	campaign := contract.CampaignResponse{
		ID:      campaignId,
		Name:    "Test",
		Content: "Hi!",
		Status:  "Pending",
	}
	service.On("GetBy", mock.Anything).Return(&campaign, nil)
	req, rr := newHttpTest("GET", "/", nil)
	req = addParameter(req, "id", campaignId)

	response, status, _ := handler.CampaignGetById(rr, req)

	assert.Equal(t, 200, status)
	assert.Equal(t, campaign.ID, response.(*contract.CampaignResponse).ID)
	assert.Equal(t, campaign.Name, response.(*contract.CampaignResponse).Name)
}

func Test_CampaignsGetById_Err(t *testing.T) {
	setup()
	errExpected := errors.New("something wrong")
	service.On("GetBy", mock.Anything).Return(nil, errExpected)
	req, rr := newHttpTest("GET", "/", nil)

	_, _, errReturned := handler.CampaignGetById(rr, req)

	assert.Equal(t, errExpected.Error(), errReturned.Error())
}
