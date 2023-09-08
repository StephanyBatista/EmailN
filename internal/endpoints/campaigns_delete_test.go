package endpoints

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignsDelete_200(t *testing.T) {
	setup()
	campaignId := "xpto"
	service.On("Delete", mock.MatchedBy(func(id string) bool {
		return id == campaignId
	})).Return(nil)
	req, rr := newHttpTest("PATCH", "/", nil)
	req = addParameter(req, "id", campaignId)

	_, status, err := handler.CampaignDelete(rr, req)

	assert.Equal(t, 200, status)
	assert.Nil(t, err)
}

func Test_CampaignsDelete_Err(t *testing.T) {
	setup()
	errExpected := errors.New("something wrong")
	service.On("Delete", mock.Anything).Return(errExpected)
	req, rr := newHttpTest("PATCH", "/", nil)

	_, _, err := handler.CampaignDelete(rr, req)

	assert.Equal(t, errExpected, err)
}
