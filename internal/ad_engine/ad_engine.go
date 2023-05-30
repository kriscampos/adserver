package ad_engine

import (
	"time"

	"github.com/kriscampos/adserver/internal/ad_engine/ordered_multi_list"
	"github.com/kriscampos/adserver/internal/campaign"
)

// AdEngine produces relevant campaigns from a body of campaigns and keywords.
type AdEngine struct {
	updateTicker        *time.Ticker
	updateFunctions     map[int64][]func()
	closeUpdater        chan bool
	campaignManager     *ordered_multi_list.OrderedMultiList
	impressionURLToNode map[string]*ordered_multi_list.Node
}

func NewAdEngine() *AdEngine {
	return &AdEngine{
		updateFunctions:     make(map[int64][]func()),
		campaignManager:     ordered_multi_list.NewOrderedMultiList(),
		impressionURLToNode: make(map[string]*ordered_multi_list.Node),
	}
}

// Begins activation / deactivation management for campaigns.
func (a *AdEngine) Start() {
	a.updateTicker = time.NewTicker(time.Second)
	a.closeUpdater = make(chan bool)
	go func() {
		for {
			select {
			case t := <-a.updateTicker.C:
				updateFunctions, ok := a.updateFunctions[t.Unix()]
				if ok {
					for _, updateFunction := range updateFunctions {
						updateFunction()
					}
				}
			case <-a.closeUpdater:
				a.updateTicker.Stop()
			}
		}
	}()
}

func (a *AdEngine) Stop() {
	a.closeUpdater <- true
}

// Registers a campaign to be activated or deactivated based on its start and end timestamp.
func (a *AdEngine) RegisterCampaign(campaign *campaign.Campaign) {
	now := time.Now()

	campaignNode := ordered_multi_list.NewNode(campaign)
	a.impressionURLToNode[campaign.ImpressionURL] = campaignNode
	if now.Before(campaign.StartTimestamp) {
		a.updateFunctions[campaign.StartTimestamp.Unix()] = append(a.updateFunctions[campaign.StartTimestamp.Unix()], func() {
			a.campaignManager.Insert(campaignNode, campaign.TargetKeywords)
		})
		a.updateFunctions[campaign.EndTimestamp.Unix()] = append(a.updateFunctions[campaign.EndTimestamp.Unix()], func() {
			a.DeleteCampaign(campaign.ImpressionURL)
		})
	} else if now.After(campaign.StartTimestamp) && now.Before(campaign.EndTimestamp) {
		a.campaignManager.Insert(campaignNode, campaign.TargetKeywords)
		a.updateFunctions[campaign.EndTimestamp.Unix()] = append(a.updateFunctions[campaign.EndTimestamp.Unix()], func() {
			a.DeleteCampaign(campaign.ImpressionURL)
		})
	}
}

// Returns the highest priority ad for the given keywords.
func (a *AdEngine) RecommendCampaign(keywords []string) (*campaign.Campaign, bool) {
	var bestCampaign *campaign.Campaign = nil
	for _, keyword := range keywords {
		campaign, ok := a.campaignManager.GetFirst(keyword)
		if ok {
			if bestCampaign == nil {
				bestCampaign = campaign
			} else if bestCampaign.Compare(campaign) > 0 {
				bestCampaign = campaign
			}
		}
	}
	if bestCampaign == nil {
		return nil, false
	}
	return bestCampaign, true
}

// Removes a campaign from being recommended.
func (a *AdEngine) DeleteCampaign(impressionURL string) {
	node := a.impressionURLToNode[impressionURL]
	a.campaignManager.Delete(node)
	delete(a.impressionURLToNode, impressionURL)
}
