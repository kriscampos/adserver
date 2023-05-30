package campaign

import (
	"time"

	"github.com/google/uuid"
)

type CampaignService struct {
	impressionUrlToCampaign map[string]*Campaign
	nextCampaignId          int
}

func NewCampaignService() *CampaignService {
	return &CampaignService{impressionUrlToCampaign: make(map[string]*Campaign)}
}

func (s *CampaignService) CreateCampaign(c *PostCampaignRequest) *Campaign {
	id := s.nextCampaignId
	s.nextCampaignId++
	newCampaign := &Campaign{
		ID:              id,
		StartTimestamp:  time.Unix(c.StartTimestamp, 0),
		EndTimestamp:    time.Unix(c.EndTimestamp, 0),
		TargetKeywords:  c.TargetKeywords,
		ImpressionCount: 0,
		MaxImpression:   c.MaxImpression,
		CPM:             c.CPM,
		ImpressionURL:   uuid.NewString(),
	}
	s.impressionUrlToCampaign[newCampaign.ImpressionURL] = newCampaign
	return newCampaign
}

// Increments impression count and returns whether the max was hit and if the
// impression url was valid.
func (s *CampaignService) IncrementImpression(impressionURL string) (bool, bool) {
	c, ok := s.impressionUrlToCampaign[impressionURL]
	if ok {
		c.ImpressionCount += 1
		if c.ImpressionCount == c.MaxImpression {
			return true, true
		}
		return false, true
	}
	return false, false
}
