package ordered_multi_list

import "github.com/kriscampos/adserver/internal/campaign"

// A collection of ordered linked lists where lists share nodes.
//
// Conceptually, this can also be thought of as a graph with
// edge types where following a single type would produce an
// ordered sequence of values.
//
// The benefits of this structure are:
//
//	1.) Fast O(1) retrieval of largest element in each list.
//	2.) Shared nodes amongst lists without duplication of data.
//	3.) Once an element is found in one list, we know where it is in every list
//		it belongs to.
type OrderedMultiList struct {
	lists map[string]*Node
}

func NewOrderedMultiList() *OrderedMultiList {
	return &OrderedMultiList{
		lists: make(map[string]*Node),
	}
}

func (o *OrderedMultiList) GetFirst(listName string) (*campaign.Campaign, bool) {
	n, ok := o.lists[listName]
	if !ok {
		return nil, ok
	}
	return n.Data, ok
}

// Inserts Node into lists.
func (o *OrderedMultiList) Insert(n *Node, listNames []string) {
	listNames = append(listNames, "")
	remainingNamesToPrev := map[string]*Node{}
	for _, listName := range listNames {
		if inserted := o.insertAtHead(n, listName); !inserted {
			remainingNamesToPrev[listName] = nil
		}
	}
	if len(remainingNamesToPrev) > 0 {
		// find insertion point and collect prev references
		var prev *Node = nil
		current, i, ok := o.lists[""], 0, true
		for ok && current.Data.Compare(n.Data) < 0 {
			prev = current
			current, ok = current.Next[""]
			i++
			for name := range remainingNamesToPrev {
				if namePrev, ok := prev.Prev[name]; ok {
					remainingNamesToPrev[name] = namePrev
				}
			}
		}
		remainingNamesToPrev[""] = prev
		for name, prev := range remainingNamesToPrev {
			if prev == nil {
				prev = o.lists[name]
			}
			o.insertAfterNode(n, prev, name)
		}
	}
}

// Attempts to insert at head of a list if sort order is not compromised.
func (o *OrderedMultiList) insertAtHead(n *Node, listName string) bool {
	head, ok := o.lists[listName]
	if !ok { // empty list
		o.lists[listName] = n
		delete(n.Next, listName)
		return true
	} else if n.Data.Compare(head.Data) < 0 { // insert at 0th index
		n.Next[listName] = head
		head.Prev[listName] = n
		o.lists[listName] = n
		return true
	}
	return false
}

func (o *OrderedMultiList) insertAfterNode(n *Node, prev *Node, listName string) {
	next, ok := prev.Next[listName]
	prev.Next[listName] = n
	n.Prev[listName] = prev
	if ok && next != nil {
		n.Next[listName] = next
		next.Prev[listName] = n
	}
}

// removes Node n from all lists.
func (o *OrderedMultiList) Delete(n *Node) {
	// Delete connections where n is head or middle of list.
	for listName := range n.Next {
		next := n.Next[listName]
		prev, ok := n.Prev[listName]
		if ok { // n exists between two nodes
			next.Prev[listName] = prev
			prev.Next[listName] = next
			delete(n.Prev, listName)
		} else { // n is head of list
			o.lists[listName] = n.Next[listName]
			delete(o.lists[listName].Prev, listName)
		}
		delete(n.Next, listName)
	}
	// Delete connections where n is the end of the list.
	for listName := range n.Prev {
		prev := n.Prev[listName]
		delete(prev.Next, listName)
		delete(n.Prev, listName)
	}
	// If node was last member of list, remove list.
	for listName := range o.lists {
		if o.lists[listName].Data.Equal(n.Data) {
			delete(o.lists, listName)
		}
	}
}

// Returns a list of campaign IDs as an array. Only meant to be used in testing.
func (o *OrderedMultiList) getList(listName string) []int {
	list := make([]int, 0)
	for current, ok := o.lists[listName]; ok; current, ok = current.Next[listName] {
		list = append(list, current.Data.ID)
	}
	return list
}
