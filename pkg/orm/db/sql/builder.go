package sql

import (
	"context"
	"errors"
	"fmt"
	db2 "github.com/goradd/goradd/pkg/orm/db"
	. "github.com/goradd/goradd/pkg/orm/query"
	"github.com/goradd/maps"
	"strconv"
)

const countAlias = "_count"
const columnAliasPrefix = "c_"
const tableAliasPrefix = "t_"

type objectMapType = maps.SliceMap[string, any]
type aliasMapType = maps.SliceMap[string, any]
type JoinTreeItemSliceMap = maps.SliceMap[string, *JoinTreeItem]

// A Builder is a helper object to organize a Query object eventually into a SQL string.
// It builds the join tree and creates the aliases that will eventually be used to generate
// the sql and then unpack it into fields and objects. It implements the QueryBuilderI interface.
// It is used both as the overriding controller of a query, and the controller of a subquery, so its recursive.
// The approach is to gather up the parameters of the query first, build the nodes into a node tree without
// changing the nodes themselves, build the query, execute the query, and finally return the results.
type Builder struct {
	db2.QueryBuilder

	db DbI // The sql database object

	/* The variables below are populated during the sql build process */

	IsCount           bool
	IsDelete          bool
	RootDbTable       string                  // The database name for the table that is the root of the query
	RootJoinTreeItem  *JoinTreeItem           // The top of the join tree
	SubPrefix         string                  // The prefix for sub items. If this is a sub query, this gets updated
	SubqueryCounter   int                     // Helper to make unique prefixes for subqueries
	ColumnAliases     *JoinTreeItemSliceMap   // StdMap to go from an alias to a JoinTreeItem for columns, which can also get us to a node
	ColumnAliasNumber int                     // Helper to make unique generated aliases
	TableAliases      *JoinTreeItemSliceMap   // StdMap to go from an alias to a JoinTreeItem for tables
	NodeMap           map[NodeI]*JoinTreeItem // A map that gets us to a JoinTreeItem from a node.
	RowId             int                     // Counter for creating fake ids when doing distinct or orderby selects
	ParentBuilder     *Builder                // The parent builder of a subquery
}

// NewSqlBuilder creates a new Builder object.
func NewSqlBuilder(ctx context.Context, db DbI) *Builder {
	b := &Builder{
		db:            db,
		ColumnAliases: new(JoinTreeItemSliceMap),
		TableAliases:  new(JoinTreeItemSliceMap),
		NodeMap:       make(map[NodeI]*JoinTreeItem),
	}
	b.QueryBuilder.Init(ctx)
	return b
}

// Subquery adds a subquery node, which is like a mini query builder that should result in a single value.
func (b *Builder) Subquery() *SubqueryNode {
	n := NewSubqueryNode(b)
	b.IsSubquery = true
	return n
}

// Load terminates the builder, queries the database, and returns the results as an array of interfaces similar in structure to a json structure
func (b *Builder) Load() (result []map[string]interface{}) {
	b.buildJoinTree()

	b.makeColumnAliases()

	sql, args := b.generateSelectSql()
	rows, err := b.db.Query(b.Ctx, sql, args...)

	if err != nil {
		// This is possibly generating an error related to the sql itself, so put the sql in the error message.
		s := err.Error()
		s += "\nSql: " + sql

		panic(errors.New(s))
	}

	names, _ := rows.Columns()

	columnTypes := make([]GoColumnType, len(names))
	colCount := b.ColumnAliases.Len()
	for i := 0; i < colCount; i++ {
		columnTypes[i] = ColumnNodeGoType(b.ColumnAliases.Get(names[i]).Node.(*ColumnNode))
	}
	// add special aliases
	for i := colCount; i < len(names); i++ {
		columnTypes[i] = ColTypeBytes // These will be unpacked when they are retrieved
	}

	result = SqlReceiveRows(rows, columnTypes, names, b)

	return result
}

// LoadCursor terminates the builder, queries the database, and returns a cursor that can be used to step through
// the results.
//
// LoadCursor is helpful when loading a large set of data that you want to output in chunks.
// You cannot use this with Joins that create multiple connections to other objects,
// like reverse FKs or Multi-multi relationships
func (b *Builder) LoadCursor() CursorI {
	for _, n := range db2.Nodes(b.QueryBuilder) {
		if NodeIsExpander(n) && !NodeIsExpanded(n) {
			panic("You cannot use a database cursor with a multiple relationship like a reverse relationship or multi-multi relationship.")
		}
	}

	b.buildJoinTree()

	b.makeColumnAliases()

	// Hand off the generation of sql select statements to the database, since different databases generate sql differently
	sql, args := b.generateSelectSql()

	rows, err := b.db.Query(b.Ctx, sql, args...)

	if err != nil {
		// This is possibly generating an error related to the sql itself, so put the sql in the error message.
		s := err.Error()
		s += "\nSql: " + sql

		panic(errors.New(s))
	}

	names, _ := rows.Columns()
	columnTypes := make([]GoColumnType, len(names))
	colCount := b.ColumnAliases.Len()
	for i := 0; i < colCount; i++ {
		columnTypes[i] = ColumnNodeGoType(b.ColumnAliases.Get(names[i]).Node.(*ColumnNode))
	}
	// add special aliases
	for i := colCount; i < len(names); i++ {
		columnTypes[i] = ColTypeBytes // These will be unpacked when they are retrieved
	}
	return NewSqlCursor(rows, columnTypes, nil, b)
}

func (b *Builder) Delete() {
	b.IsDelete = true
	b.buildJoinTree()
	sql, args := b.generateDeleteSql()
	_, err := b.db.Exec(b.Ctx, sql, args...)
	if err != nil {
		panic(err)
	}
}

// Count creates a query that selects one thing, a count. If distinct is specified, only distinct items will be selected.
// If no columns are specified, the count will include NULL items. Otherwise, it will not include NULL results in the count.
// You cannot include any other select items in a count. If you want to do that, you should do a normal query and add a
// COUNT operation node.
func (b *Builder) Count(distinct bool, nodes ...NodeI) uint {
	var result []map[string]interface{}

	b.IsCount = true

	if len(b.Selects) > 0 {
		panic("cannot count a query that also has items selected. Use an alias for a Count node instead")
	}
	if len(b.GroupBys) > 0 {
		panic("cannot count a query that also has group by items. Use an alias for a Count node instead")
	}

	n := NewCountNode(nodes...)
	if distinct {
		n = n.Distinct()
	}

	b.Alias(countAlias, n)
	b.buildJoinTree()

	sql, args := b.generateSelectSql()
	rows, err := b.db.Query(b.Ctx, sql, args...)

	if err != nil {
		panic(err)
	}

	names, _ := rows.Columns()
	columnTypes := []GoColumnType{ColTypeUnsigned}
	result = SqlReceiveRows(rows, columnTypes, names, nil)

	return result[0][countAlias].(uint)
}

// After the intention of the query is gathered, this will add the various nodes from the query
// to the node tree to establish the joins.
func (b *Builder) buildJoinTree() {
	nodes := db2.Nodes(b.QueryBuilder)
	for _, n := range nodes {
		b.addNodeToJoinTree(n)
	}
	b.assignTableAliases(b.RootJoinTreeItem)
}

// Adds the node to the join tree.
func (b *Builder) addNodeToJoinTree(n NodeI) {
	var node NodeI
	var tableName string
	var hasSubquery bool // Turns off the check to make sure all nodes come from the same table, since subqueries might have different tables

	nodes := b.gatherContainedNodes(n)

	for _, node = range nodes {
		if sq, ok := node.(*SubqueryNode); ok {
			hasSubquery = true
			b.SubqueryCounter++
			b2 := SubqueryBuilder(sq).(*Builder)
			b2.SubPrefix = strconv.Itoa(b.SubqueryCounter) + "_"
			b2.ParentBuilder = b
			b2.buildJoinTree()
			continue
		}

		rootNode := RootNode(node)
		if rootNode == nil {
			continue // An operation or value node perhaps
		}
		tableName = NodeTableName(rootNode)

		if b.RootDbTable == "" {
			b.RootDbTable = tableName
		} else if b.RootDbTable != tableName {
			if !hasSubquery && !b.IsSubquery {
				panic("Attempting to add a node that is not starting at the table being queried.")
			} else {
				continue
			}
		}

		// walk the current node tree and find an insertion point
		if b.RootJoinTreeItem == nil {
			b.RootJoinTreeItem = &JoinTreeItem{Node: rootNode}
			b.mapNode(rootNode, b.RootJoinTreeItem)
		}

		b.mergeNode(rootNode, b.RootJoinTreeItem)
	}
}

// gatherContainedNodes will return the given node and all the nodes "contained" by the given node, including the given
// node if it makes sense. Contained nodes are nodes that need to become part of the join tree, but that
// are embedded inside operations, subqueries, etc.
func (b *Builder) gatherContainedNodes(n NodeI) (nodes []NodeI) {
	if sn, ok := n.(*SubqueryNode); ok {
		nodes = append(nodes, n) // Return the subquery node itself, because we need to do some work on it

		// must expand the returned nodes one more time
		for _, n2 := range db2.Nodes(SubqueryBuilder(sn).(*Builder).QueryBuilder) {
			nodes = append(nodes, b.gatherContainedNodes(n2)...)
		}
	} else if cn := ContainedNodes(n); cn != nil {
		nodes = append(nodes, cn...)
	} else {
		nodes = append(nodes, n)
	}
	return
}

/*
func (b *Builder) logNode(node NodeI, level int) {
	LogNode(node, level)
	if childNodes := ChildNodes(node); childNodes != nil {
		for _, cn := range childNodes {
			b.logNode(cn, level+1)
		}
	}

}
*/

// Assuming that both nodes point to the same location, merges the source node and its children into the destination node tree
func (b *Builder) mergeNode(srcNode NodeI, destJoinItem *JoinTreeItem) {
	if !srcNode.Equals(destJoinItem.Node) {
		panic("mergeNode must start with equal nodes")
	}
	// make sure node is mapped
	b.mapNode(srcNode, destJoinItem)

	srcAliaser, ok := srcNode.(Aliaser)
	if ok &&
		srcAliaser.GetAlias() != "" &&
		srcAliaser.GetAlias() != destJoinItem.Alias {
		_, isColumnNode := srcNode.(*ColumnNode)
		if !isColumnNode {
			// Adding a pre-aliased node that is at the same level as this node, so just add it.
			b.insertNode(srcNode, destJoinItem.Parent)
			return
		}
	}

	var childNode = ChildNode(srcNode)
	if childNode == nil {
		// The srcNode already exists in the tree. Since there is nothing below it, we might have additional information
		// in this version of the node, so we add any new information to our join tree.
		if prevCond := NodeCondition(srcNode); prevCond != nil {
			if destJoinItem.JoinCondition == nil {
				destJoinItem.JoinCondition = prevCond
			} else if !destJoinItem.JoinCondition.Equals(prevCond) {
				// TODO: We need a mechanism to allow different kinds of conditional joins, perhaps through aliases so that
				// items further down the chain can be identified as to which conditional join they belong to.
				panic("Error, attempting to Join with conditions on a node which already has different conditions.")
			}
		}

		if NodeIsExpander(destJoinItem.Node) {
			if NodeIsExpanded(srcNode) {
				destJoinItem.Expanded = true
			}
		}

		return
	}

	if destJoinItem.ChildReferences == nil {
		// We have found the end of the table chain, so insert what is left
		b.insertNode(childNode, destJoinItem)
	} else {
		found := false
		for _, destChild := range destJoinItem.ChildReferences {
			if destChild.Node.Equals(childNode) {
				// found a matching child node, recurse
				b.mergeNode(childNode, destChild)
				found = true
				break
			}
		}
		if !found {
			b.insertNode(childNode, destJoinItem)
		}
	}
	return
}

// insertNode inserts the node into the join tree, adding an item into the join tree.
// If the node is already in the joinTree, it will not add it, but it WILL map the node
// to the JoinTreeItem found that matches the node.
func (b *Builder) insertNode(srcNode NodeI, parentItem *JoinTreeItem) {
	j := &JoinTreeItem{
		Node:          srcNode,
		IsPK:          NodeIsPK(srcNode),
		Expanded:      NodeIsExpanded(srcNode),
		JoinCondition: NodeCondition(srcNode),
	}

	added, matchingItem := parentItem.addChildItem(j)
	b.mapNode(srcNode, matchingItem)
	if !added {
		return
	}

	if rn := RelatedColumnNode(srcNode); rn != nil {
		b.addNodeToJoinTree(rn)
	}
	if cn := ChildNode(srcNode); cn != nil {
		b.insertNode(cn, j)
	}
}

// makeColumnAliases will build the column and table maps as it assigns aliases to the join tree. These maps determine
// what columns will appear in the select clause, and how they will be aliased.
// Generally, we will always add columns that are PKs along the chain of nodes so that we know how to unpack the objects.
// Specific situations that we do not automatically add the PK columns are: If its a distinct query, a count query,
// or we are a subquery.
// We will automatically add columns for all tables when a query does not specifically have Select statements.
// We will treat GroupBy clauses as Select statements, since most SQL drivers require that only they be selected on, and
// some even require aliases.
func (b *Builder) makeColumnAliases() {

	if len(b.GroupBys) > 0 {
		// SQL in general has a problem with group by items that are not selected, so we always select group by columns by implication
		// Some SQL forms have gotten around the problem by just choosing a random result, but modern SQL engines now consider this an error
		for _, n := range b.GroupBys {
			_, isAlias := n.(*AliasNode)

			if !isAlias {
				b.assignAlias(b.GetItemFromNode(n))
			}
		}
	} else if len(b.Selects) > 0 {
		for _, n := range b.Selects {
			b.assignAlias(b.GetItemFromNode(n))
		}
		// We must also select on orderby's, or we cannot actually order by them
		for _, n := range b.OrderBys {
			b.assignAlias(b.GetItemFromNode(n))
		}

		if !(b.IsDistinct || b.IsSubquery || b.IsCount) {
			// Have some selects, so go through and make sure all primary keys in the chain are selected on
			b.assignPrimaryKeyAliases(b.RootJoinTreeItem)
		}
	} else {
		if b.IsSubquery {
			// Subqueries must have specific columns selected. They might be as alias columns, so we do not panic here.
			if !(b.IsDistinct || b.IsCount) {
				// Still add pks so we can unpack this
				b.assignPrimaryKeyAliases(b.RootJoinTreeItem)
			}
		} else {
			b.assignAllColumnAliases(b.RootJoinTreeItem)
		}
	}
}

// assignTableAliases will assign aliases to the item and all children that are tables. Call this with the
// root to assign all the table aliases.
func (b *Builder) assignTableAliases(item *JoinTreeItem) {
	b.assignAlias(item)
	for _, item2 := range item.ChildReferences {
		b.assignTableAliases(item2)
	}
}

// assign aliases to all primary keys in join tree. We do this to make sure we can unpack the linked records even
// when specific tables are not called out in selects.
func (b *Builder) assignPrimaryKeyAliases(item *JoinTreeItem) {
	if item.Leafs == nil || !item.Leafs[0].IsPK {
		b.addNodeToJoinTree(item.Node.(TableNodeI).PrimaryKeyNode())
	}

	if !item.Leafs[0].IsPK {
		panic("pk was not added")
	}

	b.assignAlias(item.Leafs[0]) // Assign the primary key alias

	// If this has a related column node, assign its alias too.
	if rn := RelatedColumnNode(item.Node); rn != nil {
		i2 := b.findJoinItem(rn)
		b.assignAlias(i2)
	}

	for _, item2 := range item.ChildReferences {
		b.assignPrimaryKeyAliases(item2)
	}
}

// assignAllColumnAliases will add every column in the given table.
// This is the default on queries that have no Select clauses just to make it easier to build queries during
// development. After a product matures, Select statements can be added to streamline the database accesses.
func (b *Builder) assignAllColumnAliases(item *JoinTreeItem) {
	if tn, ok := item.Node.(TableNodeI); ok {
		for _, sn := range tn.SelectNodes_() {
			b.addNodeToJoinTree(sn)
			b.assignAlias(b.GetItemFromNode(sn))
		}
	}
	for _, item2 := range item.ChildReferences {
		b.assignAllColumnAliases(item2)
	}
}

// assignAlias assigns an alias to the item given.
func (b *Builder) assignAlias(item *JoinTreeItem) {
	_, isColumnNode := item.Node.(*ColumnNode)

	if item.Alias == "" {
		if isColumnNode {
			item.Alias = columnAliasPrefix + b.SubPrefix + strconv.Itoa(b.ColumnAliasNumber)
			b.ColumnAliasNumber++
		} else {
			item.Alias = tableAliasPrefix + b.SubPrefix + strconv.Itoa(b.TableAliases.Len())
		}
	}

	if isColumnNode {
		b.ColumnAliases.Set(item.Alias, item)
	} else {
		b.TableAliases.Set(item.Alias, item)
	}
}

/*
Notes on the unpacking process:
This is quite tricky. Depending on the node structure, you may get repeated branches, or repeated entire structures with
individual differences.

After getting sql rows full of aliases for individual columns, we let the node structure direct how to unpack it.
We are going to do it in steps:
1) Create objects keyed by join table alias and id number. Foreign keys and Unique Reverse Fks will be a key to an object.
Reverse FKs and ManyMany relationships will be an ordered map of keys.
2) Walk the node map, assembling the structure
	a) If we arrive at a toMany relationship that is specified not to assemble as an array, we will duplicate the entire
	   structure each time.
	b) If we arrive at a toMany relationship that is arrayed, we pull in the individual items and keep walking
3) Return the assembled structure

Note that the order matters, so we put the whole thing in an OrderedMap so we can walk the whole thing in the order
that each object arrives, but then look for items in order.
*/

// unpackResult takes a flattened result set from the database that is a series of values keyed by alias, and turns them
// into a hierarchical result set that is keyed by join table alias and key.
func (b *Builder) unpackResult(rows []map[string]interface{}) (out []map[string]interface{}) {
	var o2 db2.ValueMap

	oMap := new(objectMapType)
	aliasMap := new(aliasMapType)

	// First we create a tree structure of the data that will mirror the node structure
	for _, row := range rows {
		rowId := b.unpackObject(b.RootJoinTreeItem, row, oMap)
		b.unpackSpecialAliases(rowId, row, aliasMap)
	}

	// We then walk the tree and create the final data structure as arrays
	oMap.Range(func(key string, value any) bool {
		// Duplicate rows that are part of a join that is not an array join
		out2 := b.expandNode(b.RootJoinTreeItem, value.(db2.ValueMap))
		// Add the Alias calculations specifically requested by the caller
		for _, o2 = range out2 {
			if m := aliasMap.Get(key); m != nil {
				o2[AliasResults] = m
			}
			out = append(out, o2)
		}
		return true
	})
	return out
}

// unpackObject finds the object that corresponds to parent in the row, and either adds it to the oMap, or if its
// already in the oMap, reuses the old one and adds more data to it. oMap should only contain objects of parent type.
// Returns the row id to use to refer to the row later.
func (b *Builder) unpackObject(parent *JoinTreeItem, row db2.ValueMap, oMap *objectMapType) (rowId string) {
	var obj db2.ValueMap
	var arrayKey string
	var currentArray *objectMapType

	if b.IsDistinct || b.GroupBys != nil {
		// We are not identifying the row by a PK because of one of the following:
		// 1) This is a distinct select, and we are not selecting pks to avoid affecting the results of the query, in which case we will likely need to make up some ids
		// 2) This is a groupby clause, which forces us to select only the groupby items and we cannot add a PK to the row
		// We will therefore make up a unique key to identify the row
		b.RowId++
		rowId = strconv.Itoa(b.RowId)
	} else {
		rowId = b.makeObjectKey(parent, row)
		if rowId == "" {
			// Object was not created because of a conditional expansion that failed
			return
		}
	}

	if curObj := oMap.Get(rowId); curObj != nil {
		obj = curObj.(db2.ValueMap)
	} else {
		obj = db2.NewValueMap()
		oMap.Set(rowId, obj)
	}

	for _, childItem := range parent.ChildReferences {
		if !NodeIsTableNodeI(childItem.Node) {
			panic("leaf node put in the table nodes")
		} else {
			arrayKey = NodeGoName(childItem.Node)
		}
		// if this is an embedded object, collect a group of objects
		if i, ok := obj[arrayKey]; !ok {
			// If this is the first time, create the group
			newArray := new(objectMapType)
			obj[arrayKey] = newArray
			b.unpackObject(childItem, row, newArray)
		} else {
			// Already have a group, so add to the group
			currentArray = i.(*objectMapType)
			b.unpackObject(childItem, row, currentArray)
		}
	}
	for _, leafItem := range parent.Leafs {
		b.unpackLeaf(leafItem, row, obj)
	}
	return
}

func (b *Builder) unpackLeaf(j *JoinTreeItem, row db2.ValueMap, obj db2.ValueMap) {
	var key string
	var fieldName string

	switch node := j.Node.(type) {
	case *ColumnNode:
		key = j.Alias
		if b.ColumnAliases != nil &&
			b.ColumnAliases.Has(key) &&
			(b.AliasNodes == nil ||
				!b.AliasNodes.Has(key)) { // could be a special alias, which we should unpack differently

			fieldName = ColumnNodeDbName(node)
			obj[fieldName] = row[key]
		}
	default:
		panic("Unexpected node type.")
	}
}

// makeObjectKey makes the key for the object of the row. This should be called only once per row.
// the key is used in subsequent calls to determine what row joined data belongs to.
func (b *Builder) makeObjectKey(j *JoinTreeItem, row db2.ValueMap) string {
	var alias interface{}

	pkItem := b.getPkJoinTreeItem(j)
	if pkItem == nil {
		return "" // Primary key was not selected on, which could happen if this is distinct, count, failed expansion, etc.
	}
	if alias, _ = row[pkItem.Alias]; alias == nil {
		return ""
	}
	pk := fmt.Sprint(alias) // primary keys are usually either integers or strings in most databases. We standardize on strings.
	return pkItem.Alias + "." + pk
}

// getPkJoinTreeItem returns the join item corresponding to the primary key contained in the given node,
// meaning the node given should be a table node.
func (b *Builder) getPkJoinTreeItem(j *JoinTreeItem) *JoinTreeItem {
	if j.Leafs != nil && j.Leafs[0].IsPK {
		return j.Leafs[0]
	}
	return nil
}

// findChildJoinItem returns the JoinTreeItem matching the given node
func (b *Builder) findChildJoinItem(childNode NodeI, parent *JoinTreeItem) (match *JoinTreeItem) {
	if _, ok := childNode.(TableNodeI); ok {
		for _, cj := range parent.ChildReferences {
			if cj.Node.Equals(childNode) {
				return cj
			}
		}
	} else {
		for _, cj := range parent.Leafs {
			if cj.Node.Equals(childNode) {
				return cj
			}
		}
	}
	return nil
}

// findForeignKeyItem will find the matching leaf node for a forward referencing foreignKey.
/*func (b *Builder) findForeignKeyItem(item *JoinTreeItem) (match *JoinTreeItem) {
	parent := item.parent
	if parent == nil {
		panic("Trying to find a foreign key item on a top-level item")
	}
	refNode,ok := item.node.(*ReferenceNode)
	if !ok {
		panic("Can only find a foreign key item on a reference node")
	}
	for _,leaf := range parent.leafs {
		col := leaf.node.(*ColumnNode).name()
		if col.
	}
}*/

// findChildJoinItemRecursive recursively finds a join item
func (b *Builder) findChildJoinItemRecursive(n NodeI, joinItem *JoinTreeItem) *JoinTreeItem {
	childNode := ChildNode(n)

	if childNode == nil {
		return joinItem
	} else {
		match := b.findChildJoinItem(childNode, joinItem)

		if match != nil {
			return b.findChildJoinItemRecursive(childNode, match)
		}
	}

	return nil
}

// findJoinItem starts the process of finding a JoinTreeItem corresponding to a particular node
func (b *Builder) findJoinItem(n NodeI) *JoinTreeItem {
	return b.findChildJoinItemRecursive(RootNode(n), b.RootJoinTreeItem)
}

func (b *Builder) mapNode(node NodeI, item *JoinTreeItem) {
	if item == nil {
		panic("linking nil item")
	}
	if item.Node == nil {
		panic("linking nil node in item")
	}
	b.NodeMap[node] = item
}

func (b *Builder) GetItemFromNode(node NodeI) *JoinTreeItem {
	j := b.NodeMap[node]
	if j == nil && b.ParentBuilder != nil { // if we are in a subquery, ask the parent query for the item
		j = b.ParentBuilder.GetItemFromNode(node)
	}
	if j == nil {
		// For some reason, saving in the nodeMap does not always work. Lets double-check and walk the tree to find it.
		j = b.findJoinItem(node)
	}

	return j
}

func (b *Builder) GetAliasedNode(node AliasNodeI) NodeI {
	j := b.AliasNodes.Get(node.GetAlias())
	return j.(NodeI)
}

// Craziness of handling situation where an array node wants to be individually expanded.
func (b *Builder) expandNode(j *JoinTreeItem, nodeObject db2.ValueMap) (outArray []db2.ValueMap) {
	var item db2.ValueMap
	var innerNodeObject db2.ValueMap
	var copies []db2.ValueMap
	var innerCopies []db2.ValueMap
	var newArray []db2.ValueMap

	outArray = append(outArray, db2.NewValueMap())

	// order of reference or leaf processing is not important
	for _, childItem := range j.Leafs {
		for _, item = range outArray {
			if cn, ok := childItem.Node.(*ColumnNode); ok {
				dbName := ColumnNodeDbName(cn)
				item[dbName] = nodeObject[dbName]
			}
		}
	}

	for _, childItem := range append(j.ChildReferences) {
		copies = []db2.ValueMap{}
		tableGoName := NodeGoName(childItem.Node)

		for _, item = range outArray {
			switch NodeGetType(childItem.Node) {
			case ReferenceNodeType:
				// Should be a one or zero item array here
				om := nodeObject[tableGoName].(*objectMapType)
				if om.Len() > 1 {
					panic("Cannot have an array with more than one item here.")
				} else if om.Len() == 1 {
					innerNodeObject = nodeObject[tableGoName].(*objectMapType).GetAt(0).(db2.ValueMap)
					innerCopies = b.expandNode(childItem, innerNodeObject)
					if len(innerCopies) > 1 {
						for _, cp2 := range innerCopies {
							nodeCopy := item.Copy().(db2.ValueMap)
							nodeCopy[tableGoName] = cp2
							copies = append(copies, nodeCopy)
						}
					} else {
						item[tableGoName] = map[string]interface{}(innerCopies[0])
					}
				}
				// else we likely were not included because of a conditional join
			case ReverseReferenceNodeType:
				if childItem.Expanded { // unique reverse or single expansion many
					newArray = []db2.ValueMap{}
					nodeObject[tableGoName].(*objectMapType).Range(func(key string, value interface{}) bool {
						innerNodeObject = value.(db2.ValueMap)
						innerCopies = b.expandNode(childItem, innerNodeObject)
						for _, ic := range innerCopies {
							newArray = append(newArray, ic)
						}
						return true
					})
					for _, cp2 := range newArray {
						nodeCopy := item.Copy().(db2.ValueMap)
						nodeCopy[tableGoName] = cp2
						copies = append(copies, nodeCopy)
					}
				} else {
					// From this point up, we should not be creating additional copies, since from this point down, we
					// are gathering an array
					newArray = []db2.ValueMap{}
					nodeObject[tableGoName].(*objectMapType).Range(func(key string, value interface{}) bool {
						innerNodeObject = value.(db2.ValueMap)
						innerCopies = b.expandNode(childItem, innerNodeObject)
						for _, ic := range innerCopies {
							newArray = append(newArray, ic)
						}
						return true
					})
					item[tableGoName] = newArray
				}

			case ManyManyNodeType:
				if ManyManyNodeIsTypeTable(childItem.Node.(TableNodeI).EmbeddedNode_().(*ManyManyNode)) {
					var intArray []uint
					nodeObject[tableGoName].(*objectMapType).Range(func(key string, value interface{}) bool {
						innerNodeObject = value.(db2.ValueMap)
						typeKey := innerNodeObject[ColumnNodeDbName(childItem.Node.(TableNodeI).PrimaryKeyNode())]
						switch v := typeKey.(type) {
						case uint:
							intArray = append(intArray, v)
						case int:
							intArray = append(intArray, uint(v))
						case int64:
							intArray = append(intArray, uint(v))
						}
						return true
					})
					if !childItem.Expanded { // single expansion many
						item[tableGoName] = intArray
					} else {
						for _, cp2 := range intArray {
							nodeCopy := item.Copy().(db2.ValueMap)
							nodeCopy[tableGoName] = []uint{cp2}
							copies = append(copies, nodeCopy)
						}
					}

				} else {
					newArray = []db2.ValueMap{}
					nodeObject[tableGoName].(*objectMapType).Range(func(key string, value interface{}) bool {
						innerNodeObject = value.(db2.ValueMap)
						innerCopies = b.expandNode(childItem, innerNodeObject)
						for _, ic := range innerCopies {
							newArray = append(newArray, ic)
						}
						return true
					})
					if !childItem.Expanded {
						item[tableGoName] = newArray
					} else {
						for _, cp2 := range newArray {
							nodeCopy := item.Copy().(db2.ValueMap)
							nodeCopy[tableGoName] = []db2.ValueMap{cp2}
							copies = append(copies, nodeCopy)
						}
					}
				}
			}

		}
		if len(copies) > 0 {
			outArray = copies
		}
	}

	return
}

// unpack the manually aliased items from the result
func (b *Builder) unpackSpecialAliases(rowId string, row db2.ValueMap, aliasMap *aliasMapType) {
	var obj db2.ValueMap

	if curObj := aliasMap.Get(rowId); curObj != nil {
		return // already added these to the row
	} else {
		obj = db2.NewValueMap()
	}

	if b.AliasNodes != nil {
		b.AliasNodes.Range(func(key string, value Aliaser) bool {
			obj[key] = row[key]
			return true
		})
	}

	if len(obj) > 0 {
		aliasMap.Set(rowId, obj)
	}
}

func (b *Builder) generateSelectSql() (sql string, args []any) {
	g := newSelectGenerator(b)
	sql = g.generateSelectSql()
	args = g.argList
	return
}

func (b *Builder) generateDeleteSql() (sql string, args []any) {
	g := newSelectGenerator(b)
	sql = g.generateDeleteSql()
	args = g.argList
	return
}
