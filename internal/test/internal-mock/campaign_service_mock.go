package internalmock

import (
	"emailn/internal/domain/campaign"

	"github.com/stretchr/testify/mock"
)

type CampaignServiceMock struct {
	mock.Mock
}

func (r *CampaignServiceMock) Create(newCampaign campaign.NewCampaignRequest) (string, error) {
	args := r.Called(newCampaign)
	return args.String(0), args.Error(1)
}

func (r *CampaignServiceMock) GetBy(id string) (*campaign.CampaignResponse, error) {
	args := r.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*campaign.CampaignResponse), args.Error(1)
}

func (r *CampaignServiceMock) Delete(id string) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *CampaignServiceMock) Start(id string) error {
	args := r.Called(id)
	return args.Error(0)
}
