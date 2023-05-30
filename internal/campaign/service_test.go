package campaign

import (
	"testing"
	"time"
)

func TestCreateCampaign(t *testing.T) {
	s := NewCampaignService()
	postCampaignRequest := &PostCampaignRequest{
		StartTimestamp: 1684616602, // Saturday, May 20, 2023 9:03:22 PM GMT
		EndTimestamp:   1687295002, // Tuesday, June 20, 2023 9:03:22 PM GMT
		TargetKeywords: []string{"dog"},
		MaxImpression:  10,
		CPM:            5.0,
	}
	campaignModel := s.CreateCampaign(postCampaignRequest)

	// Verify underlying storage was updated.
	if equals := len(s.impressionUrlToCampaign) == 1; !equals {
		t.Errorf("Campaign was not properly stored. Storage: %+v", s.impressionUrlToCampaign)
	}

	if equals := campaignModel.StartTimestamp.Equal(time.Unix(postCampaignRequest.StartTimestamp, 0)); !equals {
		t.Errorf("StartTimestamp was not converted properly. Expected: %s Found %s",
			time.Unix(postCampaignRequest.StartTimestamp, 0).String(),
			campaignModel.StartTimestamp.String())
	}
}

func TestIncrementImpression(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		name           string
		getServiceFunc func() (*CampaignService, string)
		expectOK       bool
		expected       bool
	}{
		{
			name: "Increment to less than max",
			getServiceFunc: func() (*CampaignService, string) {
				s := NewCampaignService()
				c := s.CreateCampaign(&PostCampaignRequest{
					StartTimestamp: now.Unix(),
					EndTimestamp:   now.Add(3 * time.Hour).Unix(),
					TargetKeywords: []string{"dog"},
					MaxImpression:  100,
					CPM:            2.4,
				})
				return s, c.ImpressionURL
			},
			expectOK: true,
			expected: false,
		},
		{
			name: "Increment to max",
			getServiceFunc: func() (*CampaignService, string) {
				s := NewCampaignService()
				c := s.CreateCampaign(&PostCampaignRequest{
					StartTimestamp: now.Unix(),
					EndTimestamp:   now.Add(3 * time.Hour).Unix(),
					TargetKeywords: []string{"dog"},
					MaxImpression:  1,
					CPM:            2.4,
				})
				return s, c.ImpressionURL
			},
			expectOK: true,
			expected: true,
		},
		{
			name: "Ensure bad url returns not okay",
			getServiceFunc: func() (*CampaignService, string) {
				s := NewCampaignService()
				s.CreateCampaign(&PostCampaignRequest{
					StartTimestamp: now.Unix(),
					EndTimestamp:   now.Add(3 * time.Hour).Unix(),
					TargetKeywords: []string{"dog"},
					MaxImpression:  1,
					CPM:            2.4,
				})
				return s, "badURL"
			},
			expectOK: false,
			expected: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s, impressionURL := tc.getServiceFunc()
			reachedMax, ok := s.IncrementImpression(impressionURL)
			if ok != tc.expectOK {
				t.Errorf("OK: Expected %t but Found %t\n", tc.expectOK, ok)
			}
			if ok && reachedMax != tc.expected {
				t.Errorf("Max Impressions Reached: Expected %t but Found %t\n", tc.expected, reachedMax)
			}
		})
	}
}
