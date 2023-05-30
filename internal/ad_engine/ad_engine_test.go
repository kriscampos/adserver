package ad_engine

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kriscampos/adserver/internal/campaign"
)

// TODO: Refactor ad engine to inject a mock clock for testing.
// Might just be able to use a mock ticker. need to investigate.

func TestRegisterCampaign(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		name      string
		campaigns []*campaign.Campaign
		expected  *campaign.Campaign
	}{
		{
			name: "Attempt to register currently active campaign.",
			campaigns: []*campaign.Campaign{
				{
					ID:              0,
					StartTimestamp:  now,
					EndTimestamp:    now.Add(24 * time.Hour),
					TargetKeywords:  []string{"cat"},
					ImpressionCount: 0,
					MaxImpression:   1,
					CPM:             2.0,
					ImpressionURL:   "ad0",
				},
			},
			expected: &campaign.Campaign{
				ID:              0,
				StartTimestamp:  now,
				EndTimestamp:    now.Add(24 * time.Hour),
				TargetKeywords:  []string{"cat"},
				ImpressionCount: 0,
				MaxImpression:   1,
				CPM:             2.0,
				ImpressionURL:   "ad0",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			adEngine := NewAdEngine()
			adEngine.Start()
			defer adEngine.Stop()

			for _, campaign := range tc.campaigns {
				adEngine.RegisterCampaign(campaign)
			}

			topCampaign, ok := adEngine.campaignManager.GetFirst("cat")
			if !ok {
				t.Error("Expected successful recommendation but received not okay instead.")
			}
			if equals := cmp.Equal(tc.expected, topCampaign); !equals {
				t.Errorf("Was recommended incorrect Ad. Expected: %+v and Received: %+v", tc.expected, topCampaign)
			}
		})
	}
}

func TestRecommendCampaign(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		name      string
		campaigns []*campaign.Campaign
		keywords  []string
		expected  *campaign.Campaign
	}{
		{
			name: "Successfully Register New Active Campaign",
			campaigns: []*campaign.Campaign{
				{
					ID:              0,
					StartTimestamp:  now,
					EndTimestamp:    now.Add(24 * time.Hour),
					TargetKeywords:  []string{"cat"},
					ImpressionCount: 0,
					MaxImpression:   1,
					CPM:             2.0,
					ImpressionURL:   "ad0",
				},
			},
			keywords: []string{"cat"},
			expected: &campaign.Campaign{
				ID:              0,
				StartTimestamp:  now,
				EndTimestamp:    now.Add(24 * time.Hour),
				TargetKeywords:  []string{"cat"},
				ImpressionCount: 0,
				MaxImpression:   1,
				CPM:             2.0,
				ImpressionURL:   "ad0",
			},
		},
		{
			name: "Successfully Register New Upcoming Campaign",
			campaigns: []*campaign.Campaign{
				{
					ID:              0,
					StartTimestamp:  now,
					EndTimestamp:    now.Add(24 * time.Hour),
					TargetKeywords:  []string{"cat"},
					ImpressionCount: 0,
					MaxImpression:   1,
					CPM:             2.0,
					ImpressionURL:   "ad0",
				},
			},
			keywords: []string{"cat"},
			expected: &campaign.Campaign{
				ID:              0,
				StartTimestamp:  now,
				EndTimestamp:    now.Add(24 * time.Hour),
				TargetKeywords:  []string{"cat"},
				ImpressionCount: 0,
				MaxImpression:   1,
				CPM:             2.0,
				ImpressionURL:   "ad0",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			adEngine := NewAdEngine()
			adEngine.Start()
			defer adEngine.Stop()
			for _, campaign := range tc.campaigns {
				adEngine.RegisterCampaign(campaign)
			}
			recommendedCampaign, ok := adEngine.RecommendCampaign(tc.keywords)
			if !ok {
				t.Error("Expected successful recommendation but received not okay instead.")
			}
			if equals := cmp.Equal(tc.expected, recommendedCampaign); !equals {
				t.Errorf("Recommended incorrect Ad.\nExpected: %+v\nFound: %+v.", tc.expected, recommendedCampaign)
			}
		})
	}
}
