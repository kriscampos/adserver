package ordered_multi_list

import "github.com/kriscampos/adserver/internal/campaign"

type Node struct {
	Data *campaign.Campaign
	Next map[string]*Node
	Prev map[string]*Node
}

func NewNode(data *campaign.Campaign) *Node {
	return &Node{
		Data: data,
		Next: make(map[string]*Node),
		Prev: make(map[string]*Node),
	}
}

func (n *Node) initReferences(keywords []string) {
	for _, keyword := range keywords {
		n.Next[keyword] = nil
		n.Prev[keyword] = nil
	}
}
