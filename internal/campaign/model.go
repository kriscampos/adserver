package campaign

import "time"

// Full representation of a campaign.
type Campaign struct {
	ID              int
	StartTimestamp  time.Time
	EndTimestamp    time.Time
	TargetKeywords  []string
	ImpressionCount int
	MaxImpression   int
	CPM             float64
	ImpressionURL   string
}

// Version of campaign with information provided at request time.
//
// TODO: Validate fields when binding. e.g: end timestamp cannot be in the past.
type PostCampaignRequest struct {
	StartTimestamp int64    `json:"start_timestamp" binding:"required"`
	EndTimestamp   int64    `json:"end_timestamp" binding:"required"`
	TargetKeywords []string `json:"target_keywords" binding:"required"`
	MaxImpression  int      `json:"max_impression" binding:"required"`
	CPM            float64  `json:"cpm" binding:"required"`
}

// Determines if a campaign is active.
func (c *Campaign) isActive() bool {
	now := time.Now()
	return now.After(c.StartTimestamp) && now.Before(c.EndTimestamp)
}

// Determines if two campaigns are equal.
func (c *Campaign) Equal(other *Campaign) bool {
	if len(c.TargetKeywords) != len(other.TargetKeywords) {
		return false
	}
	for i := range c.TargetKeywords {
		if c.TargetKeywords[i] != other.TargetKeywords[i] {
			return false
		}
	}
	return c.ID == other.ID &&
		c.StartTimestamp.Equal(other.StartTimestamp) &&
		c.EndTimestamp.Equal(other.EndTimestamp) &&
		c.ImpressionCount == other.ImpressionCount &&
		c.MaxImpression == other.MaxImpression &&
		c.CPM == other.CPM &&
		c.ImpressionURL == other.ImpressionURL
}

// returns -1 when this has more priority, 0 when this and other are equal,
// and 1 when this has less priority.
func (c *Campaign) Compare(other *Campaign) int {
	// Check CPM
	cpmDiff := c.CPM - other.CPM
	if cpmDiff > 0 {
		return -1
	}
	if cpmDiff < 0 {
		return 1
	}

	// CPM is equal, Check EndDates
	now := time.Now()
	timeDiff := c.EndTimestamp.Sub(now) - other.EndTimestamp.Sub(now)
	if timeDiff < 0 {
		return -1
	}
	if timeDiff > 1 {
		return 1
	}

	// EndDates are also equal. Use ID for tiebreaker.
	idDiff := c.ID - other.ID
	if idDiff < 0 {
		return -1
	}
	if idDiff > 0 {
		return 1
	}

	// This should never happen in production code.
	panic("Found two campaigns with the same ID. This should never happen.")
}
