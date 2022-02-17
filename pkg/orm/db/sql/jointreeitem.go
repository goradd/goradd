package sql

import "github.com/goradd/goradd/pkg/orm/query"

// JoinTreeItem is used to build the join tree. The join tree creates a hierarchy of joined nodes that let us
// generate aliases, serialize the query, and afterwards unpack the results.
type JoinTreeItem struct {
	Node            query.NodeI
	Parent          *JoinTreeItem
	ChildReferences []*JoinTreeItem // TableNodeI objects
	Leafs         []*JoinTreeItem
	JoinCondition query.NodeI
	Alias    string
	Expanded bool
	IsPK     bool
}

// addChildItem attempts to add the given child item. If the item was previously found, it will NOT be
// added, but the found item will be returned.
func (j *JoinTreeItem) addChildItem(child *JoinTreeItem) (added bool, match *JoinTreeItem) {
	if _, ok := child.Node.(query.TableNodeI); ok {
		for _, j2 := range j.ChildReferences {
			if j2.Node.Equals(child.Node) {
				// The node was already here
				return false, j2
			}
		}
		child.Parent = j
		j.ChildReferences = append(j.ChildReferences, child)
	} else {
		for _, j2 := range j.Leafs {
			if j2.Node.Equals(child.Node) {
				// Leaf item was found, just skip it, but save node reference
				return false, j2
			}
		}
		child.Parent = j
		if child.IsPK {
			// PKs go to the front
			j.Leafs = append([]*JoinTreeItem{child}, j.Leafs...)
		} else {
			j.Leafs = append(j.Leafs, child)
		}

	}
	return true, child
}

// pk will return the primary key join tree item attached to this item, or nil if none exists
func (j *JoinTreeItem) pk() *JoinTreeItem {
	if j.Leafs != nil &&
		j.Leafs[0].IsPK {
		return j.Leafs[0]
	} else {
		return nil
	}
}

