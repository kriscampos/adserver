package ordered_multi_list

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kriscampos/adserver/internal/campaign"
)

func TestNewOrderedLinkedLists(t *testing.T) {
	errorMsg := "%s was not initialized"
	lists := NewOrderedMultiList()
	if lists.lists == nil {
		t.Fatalf(errorMsg, "lists")
	}
}

func TestGetFirst(t *testing.T) {
	testcases := []struct {
		name            string
		GetList         func() *OrderedMultiList
		list_name       string
		expected_status bool
		expected_val    *campaign.Campaign
	}{
		{
			name: "Get on empty list.",
			GetList: func() *OrderedMultiList {
				l := NewOrderedMultiList()
				return l
			},
			list_name:       "a",
			expected_status: false,
			expected_val:    nil,
		},
		{
			name: "Get on list with one element.",
			GetList: func() *OrderedMultiList {
				l := NewOrderedMultiList()
				l.insertAtHead(NewNode(&campaign.Campaign{ID: 0}), "a")
				return l
			},
			list_name:       "a",
			expected_status: true,
			expected_val:    &campaign.Campaign{ID: 0},
		},
		{
			name: "Get on second list where elem is second in master list.",
			GetList: func() *OrderedMultiList {
				l := NewOrderedMultiList()
				l.insertAtHead(NewNode(&campaign.Campaign{ID: 0}), "a")
				l.insertAtHead(NewNode(&campaign.Campaign{ID: 1}), "b")
				return l
			},
			list_name:       "b",
			expected_status: true,
			expected_val:    &campaign.Campaign{ID: 1},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			l := tc.GetList()
			actual_val, actual_status := l.GetFirst(tc.list_name)
			if actual_status != tc.expected_status {
				t.Errorf("Status mismatch. Expected: %t but Found: %t\n", tc.expected_status, actual_status)
			}
			if actual_status && !actual_val.Equal(tc.expected_val) {
				t.Errorf("Value mismatch. Expected: %+v but Found: %+v\n", tc.expected_val, actual_val)
			}
		})
	}
}

func TestGetList(t *testing.T) {
	// set up nodes
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID: 1,
		}),
		NewNode(&campaign.Campaign{
			ID: 2,
		}),
		NewNode(&campaign.Campaign{
			ID: 3,
		}),
		NewNode(&campaign.Campaign{
			ID: 4,
		}),
	}
	// Set up Next references
	nodes[0].Next[""] = nodes[1]
	nodes[1].Next[""] = nodes[2]
	nodes[2].Next[""] = nodes[3]

	// Set up Prev references
	nodes[3].Prev[""] = nodes[2]
	nodes[2].Prev[""] = nodes[1]
	nodes[1].Prev[""] = nodes[0]

	// Set up ordered linked lists.
	lists := &OrderedMultiList{
		lists: map[string]*Node{
			"": nodes[0],
		},
	}

	actual := lists.getList("")
	expected := []int{1, 2, 3, 4}

	if equal := cmp.Equal(actual, expected); !equal {
		t.Errorf("actual == expected: %t", equal)
	}
}

func TestInsert_Empty(t *testing.T) {
	lists := NewOrderedMultiList()
	listNames := []string{}
	lists.Insert(NewNode(&campaign.Campaign{
		ID:           1,
		CPM:          5.0,
		EndTimestamp: time.Now().Add(24 * time.Hour),
	}), listNames)
	expected := []int{1}
	actual := lists.getList("")
	if equals := cmp.Equal(expected, actual); !equals {
		t.Errorf("Expected: %+v Found: %+v", expected, actual)
	}
}

func TestInsert_AtHead(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          5.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
	}
	listNames := []string{}
	for _, Node := range nodes {
		lists.Insert(Node, listNames)
	}
	expected := []int{2, 1}
	actual := lists.getList("")
	if equals := cmp.Equal(expected, actual); !equals {
		t.Errorf("Expected: %+v Found: %+v", expected, actual)
	}
}

func TestInsert_AtEnd(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          5.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          4.5,
			EndTimestamp: time.Now().Add(25 * time.Hour),
		}),
	}
	listNames := []string{}
	for _, Node := range nodes {
		lists.Insert(Node, listNames)
	}
	expected := []int{2, 1, 3}
	actual := lists.getList("")
	if equals := cmp.Equal(expected, actual); !equals {
		t.Errorf("Expected: %+v Found: %+v", expected, actual)
	}
}

func TestInsert_BetweenNodes(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          5.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          4.5,
			EndTimestamp: time.Now().Add(25 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           4,
			CPM:          4.7,
			EndTimestamp: time.Now().Add(25 * time.Hour),
		}),
	}
	listNames := []string{}
	for _, Node := range nodes {
		lists.Insert(Node, listNames)
	}
	expected := []int{2, 1, 4, 3}
	actual := lists.getList("")
	if equals := cmp.Equal(expected, actual); !equals {
		t.Errorf("Expected: %+v Found: %+v", expected, actual)
	}
}

func TestInsert_MultiList_TwoNewLists(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          5.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          6.0,
			EndTimestamp: time.Now().Add(25 * time.Hour),
		}),
	}
	listNames := [][]string{
		{"dog"},
		{"cat"},
		{"cat"},
	}
	for i, Node := range nodes {
		lists.Insert(Node, listNames[i])
	}
	expected := map[string][]int{
		"dog": {1},
		"cat": {3, 2},
		"":    {3, 2, 1},
	}
	for listName := range lists.lists {
		actual := lists.getList(listName)
		if equals := cmp.Equal(expected[listName], actual); !equals {
			t.Errorf("Expected: %+v Found: %+v", expected[listName], actual)
		}
	}
}

func TestInsert_MultiList_AtHead(t *testing.T) {
	lists := NewOrderedMultiList()
	listNames := [][]string{
		{"dog"},
		{"cat"},
	}
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          5.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
	}
	for i, Node := range nodes {
		lists.Insert(Node, listNames[i])
	}
	expected := map[string][]int{
		"dog": {1},
		"cat": {2},
		"":    {2, 1},
	}
	for listName := range lists.lists {
		actual := lists.getList(listName)
		if equals := cmp.Equal(expected[listName], actual); !equals {
			t.Errorf("Expected: %+v Found: %+v", expected[listName], actual)
		}
	}
}

func TestInsert_MultiList_EndInsert(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          5.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          2.0,
			EndTimestamp: time.Now().Add(25 * time.Hour),
		}),
	}
	listNames := [][]string{
		{"dog"},
		{"cat"},
		{"dog"},
	}
	for i, Node := range nodes {
		lists.Insert(Node, listNames[i])
	}
	expected := map[string][]int{
		"dog": {1, 3},
		"cat": {2},
		"":    {2, 1, 3},
	}
	for listName := range lists.lists {
		actual := lists.getList(listName)
		if equals := cmp.Equal(expected[listName], actual); !equals {
			t.Errorf("Expected: %+v Found: %+v", expected[listName], actual)
		}
	}
}

func TestInsert_MultiList_MiddleInsert(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          5.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          4.0,
			EndTimestamp: time.Now().Add(25 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           4,
			CPM:          4.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           5,
			CPM:          4.6,
			EndTimestamp: time.Now().Add(25 * time.Hour),
		}),
	}
	listNames := [][]string{
		{"dog"},
		{"cat"},
		{"dog"},
		{"cat"},
		{"dog"},
	}
	for i, Node := range nodes {
		lists.Insert(Node, listNames[i])
	}
	expected := map[string][]int{
		"dog": {1, 5, 3},
		"cat": {2, 4},
		"":    {2, 1, 5, 4, 3},
	}
	for listName := range lists.lists {
		actual := lists.getList(listName)
		if equals := cmp.Equal(expected[listName], actual); !equals {
			t.Errorf("Expected: %+v Found: %+v", expected[listName], actual)
		}
	}
}

func TestInsert_MultiList_InsertSeveralLists(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          5.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          4.0,
			EndTimestamp: time.Now().Add(25 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           4,
			CPM:          4.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           5,
			CPM:          4.6,
			EndTimestamp: time.Now().Add(25 * time.Hour),
		}),
	}
	listNames := [][]string{
		{"dog"},
		{"cat"},
		{"dog"},
		{"cat"},
		{"dog", "cat"},
	}
	for i, Node := range nodes {
		lists.Insert(Node, listNames[i])
	}
	expected := map[string][]int{
		"dog": {1, 5, 3},
		"cat": {2, 5, 4},
		"":    {2, 1, 5, 4, 3},
	}
	for listName := range lists.lists {
		actual := lists.getList(listName)
		if equals := cmp.Equal(expected[listName], actual); !equals {
			t.Errorf("Expected: %+v Found: %+v", expected[listName], actual)
		}
	}
}

func TestDelete_EmptyList(t *testing.T) {
	lists := NewOrderedMultiList()
	n := NewNode(&campaign.Campaign{
		ID:           1,
		CPM:          1.2,
		EndTimestamp: time.Now().Add(24 * time.Hour),
	})
	lists.Delete(n)
	expected := []int{}
	actual := lists.getList("")
	if equals := cmp.Equal(expected, actual); !equals {
		t.Error("Deleting element from empty list failed.")
	}
}

func TestDelete_Head(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          6.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
	}
	listNames := []string{}
	for _, Node := range nodes {
		lists.Insert(Node, listNames)
	}
	lists.Delete(nodes[0])
	expected := []int{2}
	actual := lists.getList("")
	if equals := cmp.Equal(expected, actual); !equals {
		t.Errorf("Expected: %+v Found: %+v", expected, actual)
	}
}

func TestDelete_Middle(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          6.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          4.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
	}
	listNames := []string{}
	for _, Node := range nodes {
		lists.Insert(Node, listNames)
	}
	lists.Delete(nodes[1])
	expected := []int{1, 3}
	actual := lists.getList("")
	if equals := cmp.Equal(expected, actual); !equals {
		t.Errorf("Expected: %+v Found: %+v", expected, actual)
	}
}

func TestDelete_End(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          6.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          4.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
	}
	listNames := []string{}
	for _, Node := range nodes {
		lists.Insert(Node, listNames)
	}
	lists.Delete(nodes[2])
	expected := []int{1, 2}
	actual := lists.getList("")
	if equals := cmp.Equal(expected, actual); !equals {
		t.Errorf("Expected: %+v Found: %+v", expected, actual)
	}
}

func TestDelete_MemberOfSeveralLists(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          5.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          4.0,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           4,
			CPM:          4.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           5,
			CPM:          4.6,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
	}
	listNames := [][]string{
		{"dog"},
		{"cat"},
		{"dog"},
		{"cat"},
		{"dog", "cat"},
	}
	for i, Node := range nodes {
		lists.Insert(Node, listNames[i])
	}
	lists.Delete(nodes[4])
	expected := map[string][]int{
		"dog": {1, 3},
		"cat": {2, 4},
		"":    {2, 1, 4, 3},
	}
	for listName := range lists.lists {
		actual := lists.getList(listName)
		if equals := cmp.Equal(expected[listName], actual); !equals {
			t.Errorf("Expected: %+v Found: %+v", expected[listName], actual)
		}
	}
}

func TestDelete_NonMemberNode(t *testing.T) {
	lists := NewOrderedMultiList()
	nodes := []*Node{
		NewNode(&campaign.Campaign{
			ID:           1,
			CPM:          6.0,
			EndTimestamp: time.Now().Add(24 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           2,
			CPM:          5.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
		NewNode(&campaign.Campaign{
			ID:           3,
			CPM:          4.5,
			EndTimestamp: time.Now().Add(23 * time.Hour),
		}),
	}
	listNames := []string{}
	for _, Node := range nodes {
		lists.Insert(Node, listNames)
	}
	lists.Delete(NewNode(&campaign.Campaign{
		ID:           4,
		CPM:          3.0,
		EndTimestamp: time.Now().Add(23 * time.Hour),
	}))
	expected := []int{1, 2, 3}
	actual := lists.getList("")
	if equals := cmp.Equal(expected, actual); !equals {
		t.Errorf("Expected: %+v Found: %+v", expected, actual)
	}
}
