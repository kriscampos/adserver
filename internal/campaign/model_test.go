package campaign

import (
	"testing"
	"time"
)

func TestIsActive(t *testing.T) {
	testcases := []struct {
		name     string
		input    Campaign
		expected bool
	}{
		{
			name: "Verify past campaign is not active",
			input: Campaign{
				StartTimestamp: time.Now().Add(-24 * time.Hour),
				EndTimestamp:   time.Now().Add(-12 * time.Hour),
			},
			expected: false,
		},
		{
			name: "Verify current campaign is active",
			input: Campaign{
				StartTimestamp: time.Now().Add(-12 * time.Hour),
				EndTimestamp:   time.Now().Add(12 * time.Hour),
			},
			expected: true,
		},
		{
			name: "Verify future campaign is not active",
			input: Campaign{
				StartTimestamp: time.Now().Add(12 * time.Hour),
				EndTimestamp:   time.Now().Add(24 * time.Hour),
			},
			expected: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := tc.input.isActive(); actual != tc.expected {
				t.Errorf("Expected %t but found %t for campaign: %+v\n", tc.expected, actual, tc.input)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		name     string
		inputs   []*Campaign
		expected bool
	}{
		{
			name: "Expect True",
			inputs: []*Campaign{
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
			expected: true,
		},
		{
			name: "Expect False - keywords length mismatch.",
			inputs: []*Campaign{
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
				{
					ID:              0,
					StartTimestamp:  now,
					EndTimestamp:    now.Add(24 * time.Hour),
					TargetKeywords:  []string{"cat", "dog"},
					ImpressionCount: 0,
					MaxImpression:   1,
					CPM:             2.0,
					ImpressionURL:   "ad0",
				},
			},
			expected: false,
		},
		{
			name: "Expect False - keywords mismatch.",
			inputs: []*Campaign{
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
				{
					ID:              0,
					StartTimestamp:  now,
					EndTimestamp:    now.Add(24 * time.Hour),
					TargetKeywords:  []string{"dog"},
					ImpressionCount: 0,
					MaxImpression:   1,
					CPM:             2.0,
					ImpressionURL:   "ad0",
				},
			},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := tc.inputs[0].Equal(tc.inputs[1]); actual != tc.expected {
				t.Errorf("Expected: %t Found: %t for inputs:\n%+v\n%+v", tc.expected, actual, tc.inputs[0], tc.inputs[1])
			}
		})
	}
}

func TestCompare(t *testing.T) {
	now := time.Now()
	testCases := [][]*Campaign{
		{
			{
				ID:           1,
				CPM:          12.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
			{
				ID:           2,
				CPM:          10.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
		}, {
			{
				ID:           1,
				CPM:          9.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
			{
				ID:           2,
				CPM:          10.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
		}, {
			{
				ID:           1,
				CPM:          10.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
			{
				ID:           2,
				CPM:          10.0,
				EndTimestamp: now.Add(3 * time.Hour),
			},
		}, {
			{
				ID:           1,
				CPM:          10.0,
				EndTimestamp: now.Add(3 * time.Hour),
			},
			{
				ID:           2,
				CPM:          10.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
		}, {
			{
				ID:           1,
				CPM:          10.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
			{
				ID:           2,
				CPM:          10.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
		},
		{
			{
				ID:           3,
				CPM:          10.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
			{
				ID:           2,
				CPM:          10.0,
				EndTimestamp: now.Add(2 * time.Hour),
			},
		},
	}
	expecteds := []int{
		-1, 1, -1, 1, -1, 1,
	}
	for i, testCase := range testCases {
		expected := expecteds[i]
		actual := testCase[0].Compare(testCase[1])
		if equals := expected == actual; !equals {
			t.Errorf("Expected %d Found %d for campaigns %+v and %+v", expected, actual, testCase[0], testCase[1])
		}
	}
}

func TestCompare_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic, but no panic occurred.")
		}
	}()

	campaigns := []*Campaign{
		{
			ID: 0,
		},
		{
			ID: 0,
		},
	}

	campaigns[0].Compare(campaigns[1])
}
