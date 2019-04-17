package db

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/goradd/gengen/pkg/maps"
	. "github.com/goradd/goradd/pkg/orm/query"
	"strconv"
)

const countAlias = "_count"
const columnAliasPrefix = "c_"
const tableAliasPrefix = "t_"

// Copier implements the copy interface, that returns a deep copy of an object.
type Copier interface {
	Copy() interface{}
}

// A sql builder is a helper object to organize a Query object eventually into a SQL string.
// It builds the join tree and creates the aliases that will eventually be used to generate
// the sql and then unpack it into fields and objects. It implements the QueryBuilderI interface.
// It is used both as the overriding controller of a query, and the controller of a subquery, so its recursive.
// The approach is to gather up the parameters of the query first, build the nodes into a node tree without
// changing the nodes themselves, build the query, execute the query, and finally return the results.
type sqlBuilder struct {
	db SqlDbI // The sql database object

	/* The variables below are populated while defining the query */

	QueryBuilder

	/* The variables below are populated during the sql build process */

	isCount bool
	isDelete bool
	rootDbTable       string                  // The database name for the table that is the root of the query
	rootJoinTreeItem  *joinTreeItem           // The top of the join tree
	subPrefix         string                  // The prefix for sub items. If this is a sub query, this gets updated
	subqueryCounter   int                     // Helper to make unique prefixes for subqueries
	columnAliases     *joinTreeItemSliceMap   // Map to go from an alias to a joinTreeItem for columns, which can also get us to a node
	columnAliasNumber int                     // Helper to make unique generated aliases
	tableAliases      *joinTreeItemSliceMap   // Map to go from an alias to a joinTreeItem for tables
	nodeMap           map[NodeI]*joinTreeItem // A map that gets us to a joinTreeItem from a node.
	rowId             int                     // Counter for creating fake ids when doing distinct or orderby selects
	parentBuilder     *sqlBuilder             // The parent builder of a subquery
}

// NewSqlBuilder creates a new sqlBuilder object.
func NewSqlBuilder(db SqlDbI) *sqlBuilder {
	b := &sqlBuilder{
		db:            db,
		columnAliases: NewjoinTreeItemSliceMap(),
		tableAliases:  NewjoinTreeItemSliceMap(),
		nodeMap:       make(map[NodeI]*joinTreeItem),
	}
	b.QueryBuilder.Init(b)
	return b
}


// Load terminates the builder, queries the database, and returns the results as an array of interfaces similar in structure to a json structure
func (b *sqlBuilder) Load(ctx context.Context) (result []map[string]interface{}) {
	b.buildJoinTree()

	b.makeColumnAliases()

	// Hand off the generation of sql select statements to the database, since different databases generate sql differently
	sql, args := b.db.generateSelectSql(b)

	rows, err := b.db.Query(ctx, sql, args...)

	if err != nil {
		// This is possibly generating an error related to the sql itself, so put the sql in the error message.
		s := err.Error()
		s += "\nSql: " + sql

		panic(errors.New(s))
	}
	defer rows.Close()

	names, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	columnTypes := make([]GoColumnType, len(names))
	colCount := b.columnAliases.Len()
	for i := 0; i < colCount; i++ {
		columnTypes[i] = ColumnNodeGoType(b.columnAliases.Get(names[i]).node.(*ColumnNode))
	}
	// add special aliases
	for i := colCount; i < len(names); i++ {
		columnTypes[i] = ColTypeBytes // These will be unpacked when they are retrieved
	}

	result = ReceiveRows(rows, columnTypes, names)

	var result2 = b.unpackResult(result)

	return result2
}

func (b *sqlBuilder) Delete(ctx context.Context) {
	b.isDelete = true
	b.buildJoinTree()

	// Hand off the generation of sql statements to the database, since different databases generate sql differently
	sql, args := b.db.generateDeleteSql(b)

	_, err := b.db.Exec(ctx, sql, args...)

	if err != nil {
		panic(err)
	}
}

// Count creates a query that selects one thing, a count. If distinct is specified, only distinct items will be selected.
// If no columns are specified, the count will include NULL items. Otherwise, it will not include NULL results in the count.
// You cannot include any other select items in a count. If you want to do that, you should do a normal query and add a
// COUNT operation node.
func (b *sqlBuilder) Count(ctx context.Context, distinct bool, nodes ...NodeI) uint {
	var result []map[string]interface{}

	b.isCount = true

	if len(b.selects) > 0 {
		panic("cannot count a query that also has items selected. Use an alias for a Count node instead")
	}

	if len(b.groupBys) > 0 {
		panic("cannot count a query that also has group by items. Use an alias for a Count node instead")
	}

	n := NewCountNode(nodes...)
	if distinct {
		n = n.Distinct()
	}

	b.Alias(countAlias, n)

	b.buildJoinTree()

	// Hand off the generation of sql select statements to the database, since different databases generate sql differently
	sql, args := b.db.generateSelectSql(b)

	rows, err := b.db.Query(ctx, sql, args...)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	names, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	columnTypes := []GoColumnType{ColTypeUnsigned}

	result = ReceiveRows(rows, columnTypes, names)

	return result[0][countAlias].(uint)

}

// After the intention of the query is gathered, this will add the various nodes from the query
// to the node tree to establish the joins.
func (b *sqlBuilder) buildJoinTree() {
	for _, n := range b.nodes() {
		b.addNodeToJoinTree(n)
	}
	b.assignTableAliases(b.rootJoinTreeItem)
}

// Returns the nodes referred to in the query. Some nodes will be container nodes, and so will have nodes
// inside them, but every node is either referred to, or contained in the returned nodes.
func (b *sqlBuilder) nodes() []NodeI {
	var nodes []NodeI
	for _, n := range b.joins {
		nodes = append(nodes, n)
		if c := NodeCondition(n); c != nil {
			nodes = append(nodes, c)
		}
	}
	nodes = append(nodes, b.orderBys...)

	if b.condition != nil {
		nodes = append(nodes, b.condition)
	}

	for _, n := range b.groupBys {
		if NodeIsTableNodeI(n) {
			n = NodePrimaryKey(n) // Allow table nodes, but then actually have them be the pk in this context
		}
		nodes = append(nodes, n)
	}

	if b.having != nil {
		nodes = append(nodes, b.having)
	}
	nodes = append(nodes, b.selects...)

	b.aliasNodes.Range(func(key string, value interface{}) bool {
		nodes = append(nodes, value.(NodeI))
		return true
	})

	return nodes
}

// Adds the node to the join tree.
func (b *sqlBuilder) addNodeToJoinTree(n NodeI) {
	var node NodeI
	var tableName string
	var hasSubquery bool // Turns off the check to make sure all nodes come from the same table, since subqueries might have different tables

	nodes := b.gatherContainedNodes(n)

	for _, node = range nodes {
		if sq, ok := node.(*SubqueryNode); ok {
			hasSubquery = true
			b.subqueryCounter++
			b2 := SubqueryBuilder(sq).(*sqlBuilder)
			b2.subPrefix = strconv.Itoa(b.subqueryCounter) + "_"
			b2.parentBuilder = b
			b2.buildJoinTree()
			continue
		}

		rootNode := RootNode(node)
		if rootNode == nil {
			continue // An operation or value node perhaps
		}
		tableName = NodeTableName(rootNode)

		if b.rootDbTable == "" {
			b.rootDbTable = tableName
		} else if b.rootDbTable != tableName {
			if !hasSubquery && !b.isSubquery {
				panic("Attempting to add a node that is not starting at the table being queried.")
			} else {
				continue
			}
		}

		// walk the current node tree and find an insertion point
		if b.rootJoinTreeItem == nil {
			b.rootJoinTreeItem = &joinTreeItem{node: rootNode}
			b.mapNode(rootNode, b.rootJoinTreeItem)
		}

		b.mergeNode(rootNode, b.rootJoinTreeItem)
	}
}

// gatherContainedNodes will return all of the nodes "contained" by the given node, including the given
// node if it make sense. Contained nodes are nodes that need to become part of the join tree, but that
// are embedded inside operations, subqueries, etc.
func (b *sqlBuilder) gatherContainedNodes(n NodeI) (nodes []NodeI) {
	if sn, ok := n.(*SubqueryNode); ok {
		nodes = append(nodes, n) // Return the subquery node itself, because we need to do some work on it

		// must expand the returned nodes one more time
		for _, n2 := range SubqueryBuilder(sn).(*sqlBuilder).nodes() {
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
func (b *sqlBuilder) logNode(node NodeI, level int) {
	LogNode(node, level)
	if childNodes := ChildNodes(node); childNodes != nil {
		for _, cn := range childNodes {
			b.logNode(cn, level+1)
		}
	}

}
*/

// Assuming that both nodes point to the same location, merges the source node and its children into the destination node tree
func (b *sqlBuilder) mergeNode(srcNode NodeI, destJoinItem *joinTreeItem) {
	if !srcNode.Equals(destJoinItem.node) {
		panic("mergeNode must start with equal nodes")
	}
	// make sure node is mapped
	b.mapNode(srcNode, destJoinItem)

	srcAliaser, ok := srcNode.(Aliaser)
	if ok &&
		srcAliaser.GetAlias() != "" &&
		srcAliaser.GetAlias() != destJoinItem.alias {
		_, isColumnNode := srcNode.(*ColumnNode)
		if !isColumnNode {
			// Adding a pre-aliased node that is at the same level as this node, so just add it.
			b.insertNode(srcNode, destJoinItem.parent)
			return
		}
	}

	var childNode = ChildNode(srcNode)
	if childNode == nil {
		// The srcNode already exists in the tree. Since there is nothing below it, we might have additional information
		// in this version of the node, so we add any new information to our join tree.
		if prevCond := NodeCondition(srcNode); prevCond != nil {
			if destJoinItem.joinCondition == nil {
				destJoinItem.joinCondition = prevCond
			} else if !destJoinItem.joinCondition.Equals(prevCond) {
				// TODO: We need a mechanism to allow different kinds of conditional joins, perhaps through aliases so that
				// items further down the chain can be identified as to which conditional join they belong to.
				panic("Error, attempting to Join with conditions on a node which already has different conditions.")
			}
		}

		if NodeIsExpander(destJoinItem.node) {
			if NodeIsExpanded(srcNode) {
				destJoinItem.expanded = true
			}
		}

		return
	}

	if destJoinItem.childReferences == nil {
		// We have found the end of the table chain, so insert what is left
		b.insertNode(childNode, destJoinItem)
	} else {
		found := false
		for _, destChild := range destJoinItem.childReferences {
			if destChild.node.Equals(childNode) {
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
// to the joinTreeItem found that matches the node.
func (b *sqlBuilder) insertNode(srcNode NodeI, parentItem *joinTreeItem) {
	j := &joinTreeItem{
		node:          srcNode,
		isPK:          NodeIsPK(srcNode),
		expanded:      NodeIsExpanded(srcNode),
		joinCondition: NodeCondition(srcNode),
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
func (b *sqlBuilder) makeColumnAliases() {

	if len(b.groupBys) > 0 {
		// SQL in general has a problem with group by items that are not selected, so we always select group by columns by implication
		// Some SQL forms have gotten aorund the problem by just choosing a random result, but modern SQL engines now consider this an error
		for _, n := range b.groupBys {
			b.assignAlias(b.getItemFromNode(n))
		}
	} else if len(b.selects) > 0 {
		for _, n := range b.selects {
			b.assignAlias(b.getItemFromNode(n))
		}
		// We must also select on orderby's, or we cannot actually order by them
		for _, n := range b.orderBys {
			b.assignAlias(b.getItemFromNode(n))
		}

		if !(b.distinct || b.isSubquery || b.isCount) {
			// Have some selects, so go through and make sure all primary keys in the chain are selected on
			b.assignPrimaryKeyAliases(b.rootJoinTreeItem)
		}
	} else {
		if b.isSubquery {
			// Subqueries must have specific columns selected. They might be as alias columns, so we do not panic here.
			if !(b.distinct || b.isCount) {
				// Still add pks so we can unpack this
				b.assignPrimaryKeyAliases(b.rootJoinTreeItem)
			}
		} else {
			b.assignAllColumnAliases(b.rootJoinTreeItem)
		}
	}
}

// assignTableAliases will assign aliases to the item and all children that are tables. Call this with the
// root to assign all the table aliases.
func (b *sqlBuilder) assignTableAliases(item *joinTreeItem) {
	b.assignAlias(item)
	for _, item2 := range item.childReferences {
		b.assignTableAliases(item2)
	}
}

// assign aliases to all primary keys in join tree. We do this to make sure we can unpack the linked records even
// when specific tables are not called out in selects.
func (b *sqlBuilder) assignPrimaryKeyAliases(item *joinTreeItem) {
	if item.leafs == nil || !item.leafs[0].isPK {
		b.addNodeToJoinTree(item.node.(TableNodeI).PrimaryKeyNode_())
	}

	if !item.leafs[0].isPK {
		panic("pk was not added")
	}

	b.assignAlias(item.leafs[0])

	for _, item2 := range item.childReferences {
		b.assignPrimaryKeyAliases(item2)
	}
}

// assignAllColumnAliases will add every column in the given table.
// This is the default on queries that have no Select clauses just to make it easier to build queries during
// development. After a product matures, Select statements can be added to streamline the database accesses.
func (b *sqlBuilder) assignAllColumnAliases(item *joinTreeItem) {
	if tn, ok := item.node.(TableNodeI); ok {
		for _, sn := range tn.SelectNodes_() {
			b.addNodeToJoinTree(sn)
			b.assignAlias(b.getItemFromNode(sn))
		}
	}
	for _,item2 := range item.childReferences {
		b.assignAllColumnAliases(item2)
	}
}

// assignAlias assigns an alias to the item given.
func (b *sqlBuilder) assignAlias(item *joinTreeItem) {
	_, isColumnNode := item.node.(*ColumnNode)

	if item.alias == "" {
		// if it doesn't have a pre-assigned alias, give it an automated one
		if a, ok := item.node.(Aliaser); ok && a.GetAlias() != "" {
			// This node has been assigned an alias by the developer, so use it
			item.alias = a.GetAlias()
		} else if isColumnNode {
			item.alias = columnAliasPrefix + b.subPrefix + strconv.Itoa(b.columnAliasNumber)
			b.columnAliasNumber++
		} else {
			item.alias = tableAliasPrefix + b.subPrefix + strconv.Itoa(b.tableAliases.Len())
		}
	}

	if isColumnNode {
		b.columnAliases.Set(item.alias, item)
	} else {
		b.tableAliases.Set(item.alias, item)
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
func (b *sqlBuilder) unpackResult(rows []map[string]interface{}) (out []map[string]interface{}) {
	var o2 ValueMap

	oMap := maps.NewSliceMap()
	aliasMap := maps.NewSliceMap()

	// First we create a tree structure of the data that will mirror the node structure
	for _, row := range rows {
		rowId := b.unpackObject(b.rootJoinTreeItem, row, oMap)
		b.unpackSpecialAliases(b.rootJoinTreeItem, rowId, row, aliasMap)
	}

	// We then walk the tree and create the final data structure as arrays
	oMap.Range(func(key string, value interface{}) bool {
		// Duplicate rows that are part of a join that is not an array join
		out2 := b.expandNode(b.rootJoinTreeItem, value.(ValueMap))
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
func (b *sqlBuilder) unpackObject(parent *joinTreeItem, row ValueMap, oMap *maps.SliceMap) (rowId string) {
	var obj ValueMap
	var arrayKey string
	var currentArray *maps.SliceMap

	if b.distinct || b.groupBys != nil {
		// We are not identifying the row by a PK because of one of the following:
		// 1) This is a distinct select, and we are not selecting pks to avoid affecting the results of the query, in which case we will likely need to make up some ids
		// 2) This is a groupby clause, which forces us to select only the groupby items and we cannot add a PK to the row
		// We will therefore make up a unique key to identify the row
		b.rowId++
		rowId = strconv.Itoa(b.rowId)
	} else {
		rowId = b.makeObjectKey(parent, row)
		if rowId == "" {
			// Object was not created because of a conditional expansion that failed
			return
		}
	}

	if curObj := oMap.Get(rowId); curObj != nil {
		obj = curObj.(ValueMap)
	} else {
		obj = NewValueMap()
		oMap.Set(rowId, obj)
	}

	for _, childItem := range parent.childReferences {
		if !NodeIsTableNodeI(childItem.node) {
			panic("leaf node put in the table nodes")
		} else {
			arrayKey = NodeGoName(childItem.node)
		}
		// if this is an embedded object, collect a group of objects
		if i, ok := obj[arrayKey]; !ok {
			// If this is the first time, create the group
			newArray := maps.NewSliceMap()
			obj[arrayKey] = newArray
			b.unpackObject(childItem, row, newArray)
		} else {
			// Already have a group, so add to the group
			currentArray = i.(*maps.SliceMap)
			b.unpackObject(childItem, row, currentArray)
		}
	}
	for _, leafItem := range parent.leafs {
		b.unpackLeaf(leafItem, row, obj)
	}
	return
}

func (b *sqlBuilder) unpackLeaf(j *joinTreeItem, row ValueMap, obj ValueMap) {
	var key string
	var fieldName string

	switch node := j.node.(type) {
	case *ColumnNode:
		key = j.alias
		if b.columnAliases.Has(key) && !b.aliasNodes.Has(key) { // could be a special alias, which we should unpack differently
			fieldName = ColumnNodeDbName(node)
			obj[fieldName] = row[key]
		}
	default:
		panic("Unexpected node type.")
	}
}

// makeObjectKey makes the key for the object of the row. This should be called only once per row.
// the key is used in subsequent calls to determine what row joined data belongs to.
func (b *sqlBuilder) makeObjectKey(j *joinTreeItem, row ValueMap) string {
	var alias interface{}

	pkItem := b.getPkJoinTreeItem(j)
	if pkItem == nil {
		return "" // Primary key was not selected on, which could happen if this is distinct, count, failed expansion, etc.
	}
	if alias, _ = row[pkItem.alias]; alias == nil {
		return ""
	}
	pk := fmt.Sprint(alias) // primary keys are usually either integers or strings in most databases. We standardize on strings.
	return pkItem.alias + "." + pk
}

// getPkJoinTreeItem returns the join item corresponding to the primary key contained in the given node,
// meaning the node given should be a table node.
func (b *sqlBuilder) getPkJoinTreeItem(j *joinTreeItem) *joinTreeItem {
	if j.leafs != nil && j.leafs[0].isPK {
		return j.leafs[0]
	}
	return nil
}

// findChildJoinItem returns the joinTreeItem matching the given node
func (b *sqlBuilder) findChildJoinItem(childNode NodeI, parent *joinTreeItem) (match *joinTreeItem) {
	if _, ok := childNode.(TableNodeI); ok {
		for _, cj := range parent.childReferences {
			if cj.node.Equals(childNode) {
				return cj
			}
		}
	} else {
		for _, cj := range parent.leafs {
			if cj.node.Equals(childNode) {
				return cj
			}
		}
	}
	return nil
}

// findChildJoinItemRecursive recursively finds a join item
func (b *sqlBuilder) findChildJoinItemRecursive(n NodeI, joinItem *joinTreeItem) (match *joinTreeItem) {
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

// findJoinItem starts the process of finding a joinTreeItem corresponding to a particular node
func (b *sqlBuilder) findJoinItem(n NodeI) *joinTreeItem {
	return b.findChildJoinItemRecursive(RootNode(n), b.rootJoinTreeItem)
}

func (b *sqlBuilder) mapNode(node NodeI, item *joinTreeItem) {
	if item == nil {
		panic("linking nil item")
	}
	if item.node == nil {
		panic("linking nil node in item")
	}
	b.nodeMap[node] = item
}

func (b *sqlBuilder) getItemFromNode(node NodeI) *joinTreeItem {
	j := b.nodeMap[node]
	if j == nil && b.parentBuilder != nil { // if we are in a subquery, ask the parent query for the item
		j = b.parentBuilder.getItemFromNode(node)
	}
	if j == nil {
		// For some reason, saving in the nodeMap does not always work. Lets double-check and walk the tree to find it.
		j = b.findJoinItem(node)
	}

	return j
}

// Craziness of handling situation where an array node wants to be individually expanded.
func (b *sqlBuilder) expandNode(j *joinTreeItem, nodeObject ValueMap) (outArray []ValueMap) {
	var item ValueMap
	var innerNodeObject ValueMap
	var copies []ValueMap
	var innerCopies []ValueMap
	var newArray []ValueMap
	var nodeCopy ValueMap

	outArray = append(outArray, NewValueMap())

	// order of reference or leaf processing is not important
	for _, childItem := range j.leafs {
		for _, item = range outArray {
			if cn, ok := childItem.node.(*ColumnNode); ok {
				dbName := ColumnNodeDbName(cn)
				item[dbName] = nodeObject[dbName]
			}
		}
	}

	for _, childItem := range append(j.childReferences) {
		copies = []ValueMap{}
		tableGoName := NodeGoName(childItem.node)

		for _, item = range outArray {
			switch NodeGetType(childItem.node) {
			case ReferenceNodeType:
				// Should be a one or zero item array here
				om := nodeObject[tableGoName].(*maps.SliceMap)
				if om.Len() > 1 {
					panic("Cannot have an array with more than one item here.")
				} else if om.Len() == 1 {
					innerNodeObject = nodeObject[tableGoName].(*maps.SliceMap).GetAt(0).(ValueMap)
					innerCopies = b.expandNode(childItem, innerNodeObject)
					if len(innerCopies) > 1 {
						for _, cp2 := range innerCopies {
							nodeCopy := item.Copy().(ValueMap)
							nodeCopy[tableGoName] = cp2
							copies = append(copies, nodeCopy)
						}
					} else {
						item[tableGoName] = map[string]interface{}(innerCopies[0])
					}
				}
				// else we likely were not included because of a conditional join
			case ReverseReferenceNodeType:
				if childItem.expanded { // unique reverse or single expansion many
					newArray = []ValueMap{}
					nodeObject[tableGoName].(*maps.SliceMap).Range(func(key string, value interface{}) bool {
						innerNodeObject = value.(ValueMap)
						innerCopies = b.expandNode(childItem, innerNodeObject)
						for _, ic := range innerCopies {
							newArray = append(newArray, ic)
						}
						return true
					})
					for _, cp2 := range newArray {
						nodeCopy = item.Copy().(ValueMap)
						nodeCopy[tableGoName] = cp2
						copies = append(copies, nodeCopy)
					}
				} else {
					// From this point up, we should not be creating additional copies, since from this point down, we
					// are gathering an array
					newArray = []ValueMap{}
					nodeObject[tableGoName].(*maps.SliceMap).Range(func(key string, value interface{}) bool {
						innerNodeObject = value.(ValueMap)
						innerCopies = b.expandNode(childItem, innerNodeObject)
						for _, ic := range innerCopies {
							newArray = append(newArray, ic)
						}
						return true
					})
					item[tableGoName] = newArray
				}

			case ManyManyNodeType:
				if ManyManyNodeIsTypeTable(childItem.node.(TableNodeI).EmbeddedNode_().(*ManyManyNode)) {
					var intArray []uint
					nodeObject[tableGoName].(*maps.SliceMap).Range(func(key string, value interface{}) bool {
						innerNodeObject = value.(ValueMap)
						typeKey := innerNodeObject[ColumnNodeDbName(childItem.node.(TableNodeI).PrimaryKeyNode_())]
						switch v := typeKey.(type) {
						case uint:
							intArray = append(intArray, v)
						case int:
							intArray = append(intArray, uint(v))
						}
						return true
					})
					if !childItem.expanded { // single expansion many
						item[tableGoName] = intArray
					} else {
						for _, cp2 := range intArray {
							nodeCopy = item.Copy().(ValueMap)
							nodeCopy[tableGoName] = []uint{cp2}
							copies = append(copies, nodeCopy)
						}
					}

				} else {
					newArray = []ValueMap{}
					nodeObject[tableGoName].(*maps.SliceMap).Range(func(key string, value interface{}) bool {
						innerNodeObject = value.(ValueMap)
						innerCopies = b.expandNode(childItem, innerNodeObject)
						for _, ic := range innerCopies {
							newArray = append(newArray, ic)
						}
						return true
					})
					if !childItem.expanded {
						item[tableGoName] = newArray
					} else {
						for _, cp2 := range newArray {
							nodeCopy = item.Copy().(ValueMap)
							nodeCopy[tableGoName] = []ValueMap{cp2}
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
func (b *sqlBuilder) unpackSpecialAliases(rootItem *joinTreeItem, rowId string, row ValueMap, aliasMap *maps.SliceMap) {
	var obj ValueMap

	if curObj := aliasMap.Get(rowId); curObj != nil {
		return // already added these to the row
	} else {
		obj = NewValueMap()
	}

	b.aliasNodes.Range(func(key string, value interface{}) bool {
		obj[key] = row[key]
		return true
	})

	if len(obj) > 0 {
		aliasMap.Set(rowId, obj)
	}
}

type ValueMap map[string]interface{}

func NewValueMap() ValueMap {
	return make(ValueMap)
}

// Support the deep copy interface
func (m ValueMap) Copy() interface{} {
	vm := ValueMap{}
	for k, v := range m {
		if c, ok := v.(Copier); ok {
			v = c.Copy()
		}
		vm[k] = v
	}
	return vm
}

func init() {
	gob.Register(&ValueMap{})
}

// joinTreeItem is used to build the join tree. The join tree creates a hierarchy of joined nodes that let us
// generate aliases, serialize the query, and afterwards unpack the results.
type joinTreeItem struct {
	node            NodeI
	parent          *joinTreeItem
	childReferences []*joinTreeItem // TableNodeI objects
	leafs           []*joinTreeItem
	joinCondition   NodeI
	alias           string
	expanded        bool
	isPK            bool
}

// addChildItem attempts to add the given child item. If the item was previously found, it will NOT be
// added, but the found item will be returned.
func (j *joinTreeItem) addChildItem(child *joinTreeItem) (added bool, match *joinTreeItem) {
	if _, ok := child.node.(TableNodeI); ok {
		for _, j2 := range j.childReferences {
			if j2.node.Equals(child.node) {
				// The node was already here
				return false, j2
			}
		}
		child.parent = j
		j.childReferences = append(j.childReferences, child)
	} else {
		for _, j2 := range j.leafs {
			if j2.node.Equals(child.node) {
				// Leaf item was found, just skip it, but save node reference
				return false, j2
			}
		}
		child.parent = j
		if child.isPK {
			// PKs go to the front
			j.leafs = append([]*joinTreeItem{child}, j.leafs...)
		} else {
			j.leafs = append(j.leafs, child)
		}

	}
	return true, child
}

// pk will return the primary key join tree item attached to this item, or nil if none exists
func (j *joinTreeItem) pk() *joinTreeItem {
	if j.leafs != nil &&
		j.leafs[0].isPK {
		return j.leafs[0]
	} else {
		return nil
	}
}
