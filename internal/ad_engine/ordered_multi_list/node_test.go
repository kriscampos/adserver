package ordered_multi_list

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kriscampos/adserver/internal/campaign"
)

func TestNewNode(t *testing.T) {
	errorMsg := "%s was not properly set. Expected %s but Found %s"
	newCampaign := &campaign.Campaign{
		ID: 1,
	}
	n := NewNode(newCampaign)
	if equals := cmp.Equal(n.Data, newCampaign); !equals {
		t.Errorf(errorMsg, "Didn't successully set data in node. Expected: %+v Found: %+v", newCampaign, n.Data)
	}
}

func TestInitReferences(t *testing.T) {
	errorMsg := "%s is missing from %s"
	n := NewNode(&campaign.Campaign{ID: 1, CPM: 100.0, ImpressionURL: "test"})
	keywords := []string{"cat", "dog"}
	n.initReferences(keywords)
	for _, keyword := range keywords {
		if _, ok := n.Next[keyword]; !ok {
			t.Errorf(errorMsg, keyword, "Next")
		}
		if _, ok := n.Prev[keyword]; !ok {
			t.Errorf(errorMsg, keyword, "Prev")
		}
	}
}
