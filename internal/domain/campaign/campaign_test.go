package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var (
	name      = "Campaign X"
	content   = "Body Hi!"
	contacts  = []string{"email1@e.com", "email2@e.com"}
	createdBy = "teste@teste.com.br"
	fake      = faker.New()
)

func Test_NewCampaign_CreateCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
	assert.Equal(createdBy, campaign.CreatedBy)
}

func Test_NewCampaign_IDIsNotNill(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.NotNil(campaign.ID)
}

func Test_NewCampaign_MustStatusStartWithPending(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Equal(Pending, campaign.Status)
}

func Test_NewCampaign_CreatedOnMustBeNow(t *testing.T) {
	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)

	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	assert.Greater(campaign.CreatedOn, now)
}

func Test_NewCampaign_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign("", content, contacts, createdBy)

	assert.Equal("name is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(fake.Lorem().Text(30), content, contacts, createdBy)

	assert.Equal("name is required with max 24", err.Error())
}

func Test_NewCampaign_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, "", contacts, createdBy)

	assert.Equal("content is required with min 5", err.Error())
}

func Test_NewCampaign_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, fake.Lorem().Text(1040), contacts, createdBy)

	assert.Equal("content is required with max 1024", err.Error())
}

func Test_NewCampaign_MustValidateContactsMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, nil, createdBy)

	assert.Equal("contacts is required with min 1", err.Error())
}

func Test_NewCampaign_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{"email_invalid"}, createdBy)

	assert.Equal("email is invalid", err.Error())
}

func Test_NewCampaign_MustValidateCreatedBy(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, contacts, "")

	assert.Equal("createdby is invalid", err.Error())
}

func Test_Done_ChangeStatus(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	campaign.Done()

	assert.Equal(Done, campaign.Status)
}

func Test_Start_ChangeStatus(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	campaign.Started()

	assert.Equal(Started, campaign.Status)
}

func Test_Cancel_ChangeStatus(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	campaign.Cancel()

	assert.Equal(Canceled, campaign.Status)
}

func Test_Delete_ChangeStatus(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	campaign.Delete()

	assert.Equal(Deleted, campaign.Status)
}

func Test_Fail_ChangeStatus(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(name, content, contacts, createdBy)

	campaign.Fail()

	assert.Equal(Fail, campaign.Status)
}
