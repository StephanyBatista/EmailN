package campaign_test

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	internalerrors "emailn/internal/internal-errors"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	newCampaign = contract.NewCampaign{
		Name:      "Test Y",
		Content:   "Body Hi!",
		Emails:    []string{"teste1@test.com"},
		CreatedBy: "teste@teste.com.br",
	}
	campaignPending *campaign.Campaign
	campaignStarted *campaign.Campaign
	repositoryMock  *internalmock.CampaignRepositoryMock
	service         = campaign.ServiceImp{}
)

func setUp() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
	campaignPending, _ = campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	campaignStarted = &campaign.Campaign{ID: "1", Status: campaign.Started}
}

func setUpGetByIdRepositoryBy(campaign *campaign.Campaign) {
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
}

func setUpUpdateRepository(err error) {
	repositoryMock.On("Update", mock.Anything).Return(err)
}

func setUpSendEmail(err error) {
	sendMail := func(campaign *campaign.Campaign) error {
		return err
	}
	service.SendMail = sendMail
}

func Test_Create_RequestIsValid_IdIsNotNil(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)

	assert.NotNil(t, id)
	assert.Nil(t, err)
}

func Test_Create_RequestIsNotValid_ErrInternal(t *testing.T) {
	setUp()

	_, err := service.Create(contract.NewCampaign{})

	assert.False(t, errors.Is(internalerrors.ErrInternal, err))
}

func Test_Create_RequestIsValid_CallRepository(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}

		return true
	})).Return(nil)

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ErrorOnRepository_ErrInternal(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to save on database"))

	_, err := service.Create(newCampaign)

	assert.True(t, errors.Is(internalerrors.ErrInternal, err))
}

func Test_GetById_CampaignExists_CampaignSaved(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaignPending.ID
	})).Return(campaignPending, nil)

	campaignReturned, _ := service.GetBy(campaignPending.ID)

	assert.Equal(t, campaignPending.ID, campaignReturned.ID)
	assert.Equal(t, campaignPending.Name, campaignReturned.Name)
	assert.Equal(t, campaignPending.Content, campaignReturned.Content)
	assert.Equal(t, campaignPending.Status, campaignReturned.Status)
	assert.Equal(t, campaignPending.CreatedBy, campaignReturned.CreatedBy)
}

func Test_GetById_ErrorOnRepository_ErrInternal(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong'"))

	_, err := service.GetBy("invalid_campaign")

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_CampaignWasNotFound_ErrRecordNotFound(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete("invalid_campaign")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_CampaignIsNotPending_Err(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignStarted)

	err := service.Delete(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Delete_ErrorOnRepository_ErrInternal(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))

	err := service.Delete(campaignPending.ID)

	assert.Equal(t, internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_CampaignWasDeleted_Nil(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPending == campaign
	})).Return(nil)

	err := service.Delete(campaignPending.ID)

	assert.Nil(t, err)
}

func Test_Start_CamapaignWasNotFound_ErrRecordNotFound(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("invalid_campaign")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Start_CampaignIsNotPending_Err(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignStarted)
	service.Repository = repositoryMock

	err := service.Start(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Start_CampaignWasUpdated_StatusIsStarted(t *testing.T) {
	setUp()
	setUpSendEmail(nil)
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Started
	})).Return(nil)

	service.Start(campaignPending.ID)

	assert.Equal(t, campaign.Started, campaignPending.Status)
}

func Test_SendEmailAndUpdateStatus_WhenFail_StatusIsFail(t *testing.T) {
	setUp()
	setUpSendEmail(errors.New("error to send email"))
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Fail
	})).Return(nil)

	service.SendEmailAndUpdateStatus(campaignPending)

	repositoryMock.AssertExpectations(t)
}

func Test_SendEmailAndUpdateStatus_WhenSuccess_StatusIsDone(t *testing.T) {
	setUp()
	setUpSendEmail(nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)

	service.SendEmailAndUpdateStatus(campaignPending)

	repositoryMock.AssertExpectations(t)
}
