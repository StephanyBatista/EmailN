package endpoints

import (
	"emailn/internal/domain/campaign"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	createdByExpected = "teste1@teste.com.br"
	body              = campaign.NewCampaignRequest{
		Name:    "teste",
		Content: "Hi everyone",
		Emails:  []string{"teste@teste.com"},
	}
)

func Test_CampaignsPost_201(t *testing.T) {
	setup()
	service.On("Create", mock.MatchedBy(func(request campaign.NewCampaignRequest) bool {
		if request.Name == body.Name &&
			request.Content == body.Content &&
			request.CreatedBy == createdByExpected {
			return true
		} else {
			return false
		}
	})).Return("34x", nil)
	req, rr := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)

	_, status, err := handler.CampaignPost(rr, req)

	assert.Equal(t, 201, status)
	assert.Nil(t, err)
}

func Test_CampaignsPost_Err(t *testing.T) {
	setup()
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	req, rr := newHttpTest("POST", "/", body)
	req = addContext(req, "email", createdByExpected)

	_, _, err := handler.CampaignPost(rr, req)

	assert.NotNil(t, err)
}
