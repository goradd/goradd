package driver
import (
	"strconv"
	"github.com/spekary/goradd/util/types"
	"context"
	"log"
	"fmt"
	"encoding/json"
	. 	"github.com/spekary/goradd/orm/db"

)

const (
	SELECT = "SELECT"
	INSERT = "INSERT"
	UPDATE = "UPDATE"
	DELETE = "DELETE"
)

const countAlias = "_count"

// Copier implements the copy interface, that returns a deep copy of an object.
type Copier interface {
	Copy() interface{}
}

type limitInfo struct {
	maxRowCount int64
	offset int64
}


// A sql builder is a helper object organize a Query object eventually into a SQL string.
// It builds the join tree and creates the aliases that will eventually be used to generate
// the sql and then unpack it into fields and objects.
type sqlBuilder struct {
	db SqlDbI
	command string
	columnAliases *types.OrderedMap
	columnAliasNumber int
	tableAliases *types.OrderedMap
	joins []NodeI
	orderBys []NodeI
	condition NodeI
	rootDbTable string
	rootNode NodeI
	distinct bool
	aliasNodes *types.OrderedMap
	// Adds a COUNT(*) to the select list
	isCount         bool
	groupBys        []NodeI
	selects         []NodeI
	limitInfo       *limitInfo
	having          NodeI
	distinctId      int	// Counter for creating fake ids when doing distinct selects
	isSubquery      bool
	subqueryCounter int
	subPrefix       string
}

/**
NewsqlBuilder creates a new sqlBuilder object.
 */
func NewSqlBuilder(db SqlDbI) *sqlBuilder {
	return &sqlBuilder{
		db: db,
		columnAliases: types.NewOrderedMap(),
		tableAliases: types.NewOrderedMap(),
		orderBys: []NodeI{},
		groupBys: []NodeI{},
		selects: []NodeI{},
		aliasNodes: types.NewOrderedMap(),
		joins: []NodeI{},
	}
}

func (b *sqlBuilder) Join(n NodeI, condition NodeI) QueryBuilderI {
	if condition != nil {
		if tn, ok := n.(TableNodeI); ok {
			if c, ok := tn.EmbeddedNode_().(conditioner); !ok {
				panic("Cannot set join a condition on this type of node")
			} else {
				c.setCondition(condition)
			}
		} else {
			panic("Cannot set join conditions on this type of node")
		}
	}
	b.joins = append(b.joins, n)

	if b.limitInfo != nil {

	}
	return b
}

// Add a node that is given a manual alias name. This is usually some kind of operation.
// We can recover this using the GetAlias() function of the result.
func (b *sqlBuilder) Alias(name string, n NodeI) QueryBuilderI {
	n.SetAlias(name)
	b.aliasNodes.Set(name, n)
	return b
}


// Expands an array type node so that it will produce individual rows instead of an array of items
func (b *sqlBuilder) Expand(n NodeI) QueryBuilderI {
	if tn, ok := n.(TableNodeI); !ok {
		panic("You can only expand a node that is a ReverseReference or ManyMany node.")
	} else {
		if en, ok := tn.EmbeddedNode_().(Expander); !ok {
			panic("You can only expand a node that is a ReverseReference or ManyMany node.")
		} else {
			en.Expand()
			b.Join(n, nil)
		}
	}

	return b
}


func (b *sqlBuilder) Condition(c NodeI) QueryBuilderI {
	b.condition = c
	return b
}

func (b *sqlBuilder) OrderBy(nodes... NodeI) QueryBuilderI {
	b.orderBys = append(b.orderBys, nodes...)
	return b
}

func (b *sqlBuilder) Limit(maxRowCount int64, offset int64) QueryBuilderI {
	if b.limitInfo != nil {
		panic("Query already has a limit")
	}
	b.limitInfo = &limitInfo{maxRowCount, offset}

	return b
}

func (b *sqlBuilder) Select(nodes... NodeI) QueryBuilderI {
	b.selects = append(b.selects, nodes...)
	return b
}

func (b *sqlBuilder) Distinct() QueryBuilderI {
	b.distinct = true
	return b
}

func (b *sqlBuilder) GroupBy(nodes... NodeI) QueryBuilderI {
	b.groupBys = append(b.groupBys, nodes...)
	return b
}

func (b *sqlBuilder) Having(node NodeI) QueryBuilderI {
	b.having = node
	return b
}

func (b *sqlBuilder) Subquery() *SubqueryNode {
	n := NewSubqueryNode(b)
	b.isSubquery = true
	return n
}


// Load terminates the builder, queries the database, and returns the results as an array of interfaces similar in structure to a json structure
func (b *sqlBuilder) Load(ctx context.Context) (result []map[string]interface{}) {
	b.buildNodeTree()

	b.makeColumnAliases()

	log.Println("Tree:")
	b.logNode(b.rootNode, 0)

	// So debugging will work, we declare variables
	var sql string
	var args []interface{}

	// Hand off the generation of sql select statements to the database, since different databases generate sql differently
	sql, args = b.db.generateSelectSql(b)

	log.Print(sql)

	rows, err := b.db.Query(ctx, sql, args...)

	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	names, err := rows.Columns()
	if err != nil {
		log.Panic(err)
	}

	columnTypes := make([]GoColumnType, len(names))
	colCount := b.columnAliases.Len()
	for i:= 0; i < colCount; i++ {
		columnTypes[i] = b.columnAliases.Get(names[i]).(*ColumnNode).goType
	}
	// add special aliases
	for i := colCount; i < len(names); i++ {
		columnTypes[i] = COL_TYPE_BYTES  // These will be unpacked when they are retrieved
	}

	result = ReceiveRows(rows, columnTypes, names)

	var result2 []map[string]interface{} = b.unpackResult(result)

	p, err := json.MarshalIndent(result2, "", "  ")

	log.Print(string(p))

	return result2
}

func (b *sqlBuilder) Delete(ctx context.Context) {
	b.buildNodeTree()

	log.Println("Tree:")
	b.logNode(b.rootNode, 0)

	// So debugging will work, we declare variables
	var sql string
	var args []interface{}

	// Hand off the generation of sql select statements to the database, since different databases generate sql differently
	sql, args = b.db.generateDeleteSql(b)

	log.Print(sql)

	_, err := b.db.Exec(ctx, sql, args...)

	if err != nil {
		log.Panic(err)
	}
}

// Count creates a query that selects one thing, a count. If distinct is specified, only distinct items will be selected.
// If no columns are specified, the count will include NULL items. Otherwise, it will not include NULL results in the count.
// You cannot include any other select items in a count. If you want to do that, you should do a normal query and add a COUNT column.
func (b *sqlBuilder) Count(ctx context.Context, distinct bool, nodes... NodeI) uint {
	var result = []map[string]interface{}{}

	b.isCount = true

	if len(b.selects) > 0 {
		panic ("Cannot count a query that also has items selected. Use an alias for a Count node instead.")
	}

	if len(b.groupBys) > 0 {
		panic ("Cannot count a query that also has group by items. Use an alias for a Count node instead.")
	}

	b.buildNodeTree()

	n := NewCountNode(nodes...)
	if distinct {
		n = n.Distinct()
	}

	b.Alias(countAlias, n)

	// Hand off the generation of sql select statements to the database, since different databases generate sql differently
	sql, args := b.db.generateSelectSql(b)

	log.Print(sql)

	rows, err := b.db.Query(ctx, sql, args...)

	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	names, err := rows.Columns()
	if err != nil {
		log.Panic(err)
	}

	columnTypes := []GoColumnType{COL_TYPE_UNSIGNED}

	result = ReceiveRows(rows, columnTypes, names)

	return result[0][countAlias].(uint)

}


// Adds the node to the node tree, which is what determines how we do joins. If parts of the node are not in the tree, the node
// will be added and assigned an alias. If the node already exists in the tree, the previously assigned alias will be copied to the
// node so that it can refer to the table or column by alias.
//
// If addColumn is true, and the node refers to a column node, the node will be also added to the list of nodes that are
// selected on. If addColumn is false, the node will be assigned an alias in case a future similar node is added to the column list,
// but the node will not be added to the column list.
func (b *sqlBuilder) addNode(n NodeI, addColumn bool) {
	var node, treeNode NodeI
	var tableName string
	var np NodeI
	var hasSubquery bool		// Turns off the check to make sure all nodes come from the same table, since subqueries might have different tables
	var cn []NodeI

	var nodes = []NodeI{}
	if nc, ok := n.(nodeContainer); ok {
		cn = nc.containedNodes()
		if cn != nil &&
			!b.isCount { // Adding contained nodes in this situation will impact how the count is calculated in some cases???
			nodes = append(nodes, cn...)
		}
	} else {
		nodes = append(nodes, n)
	}

	for _,node = range nodes {
		if sq,ok := node.(*SubqueryNode); ok {
			hasSubquery = true
			b.subqueryCounter++
			sq.b.(*sqlBuilder).subPrefix = strconv.Itoa(b.subqueryCounter) + "_"
			sq.b.(*sqlBuilder).buildNodeTree()
			continue
		}
		treeNode = node.rootNode()
		if treeNode == nil {continue} // could be value node or operation node. Aliased operation nodes are handled elsewhere.
		tableName = treeNode.(TableNodeI).tableName()

		if b.rootDbTable == "" {
			b.rootDbTable = tableName
		} else if b.rootDbTable != tableName {
			if !hasSubquery  && !b.isSubquery {
				panic("Attempting to add a node that is not starting at the table being queried.")
			} else{
				continue
			}
		}

		// walk the current node tree and find an insertion point
		if b.rootNode == nil {
			b.rootNode = treeNode
			b.assignAliases(b.rootNode, addColumn)
		} else {
			np = node.rootNode()
			if np == nil {
				np = node		// This is the case when we are adding an operation node that is going to be aliased
			}
			b.mergeNode(np, b.rootNode, addColumn)
		}
	}

}

func (b *sqlBuilder) logNode(node NodeI, level int) {
	node.log(level)
	if childNodes := node.getChildNodes(); childNodes != nil {
		for _,cn := range childNodes {
			b.logNode(cn, level + 1)
		}
	}

}

// assuming that both nodes point to a same location, merges the source node into the destination node tree
func (b *sqlBuilder) mergeNode(srcNode, destNode NodeI, addColumn bool) {
	var a string = srcNode.GetAlias()
	if !srcNode.Equals(destNode) {
		log.Fatal("mergeNode must start with equal nodes")
	}
	if srcNode.GetAlias() != "" &&
		srcNode.GetAlias() != destNode.GetAlias() &&
		srcNode.nodeType() != COLUMN_NODE {
		// Adding a pre-aliased node that is at the same level as this node, so just add it.
		b.insertNode(srcNode, destNode.getParentNode(), addColumn)
		return
	}
	_ = a

	var childNodes = srcNode.getChildNodes()
	if childNodes == nil {
		// The node already exists in the tree
		// Update information as needed.
		if tn, ok := srcNode.(TableNodeI); ok {
			if cn, ok := tn.EmbeddedNode_().(conditioner); ok &&
				cn.getCondition() != nil {

				if destCond := destNode.(TableNodeI).EmbeddedNode_().(conditioner).getCondition(); destCond == nil {
					destNode.(TableNodeI).EmbeddedNode_().(conditioner).setCondition(cn.getCondition())
				} else if !destCond.Equals(cn.getCondition()) {
					panic("Error, attempting to Join with conditions on a node which already has conditions.")
				}
			}
		}
		b.assignAliases(destNode, addColumn)	// potentially was added by SelectNodes_, and so did not get aliases
		if srcNode.GetAlias() == "" {
			srcNode.SetAlias(destNode.GetAlias()) // If src node does not get added, it still needs to know what its alias is
		} else { // alias was manually assigned, so use that one
			destNode.SetAlias(srcNode.GetAlias())
		}
		if p := srcNode.getParentNode(); p != nil && p.GetAlias() == "" {
			p.SetAlias(destNode.getParentNode().GetAlias()) // parent node generation for src node alias in case src node is not added to tree, but is still used in sql generation
		}


		if tn, ok := destNode.(TableNodeI); ok {
			if dn,ok := tn.EmbeddedNode_().(Expander); ok {
				// if we are expanding an array node, copy that to the destNode
				if srcNode.(TableNodeI).EmbeddedNode_().(Expander).isExpanded() {
					dn.Expand()
				}
			}
		}
		return
	}

	var srcChild NodeI

	if destNode.getChildNodes() == nil {
		// We have found the end of a chain, but we want to extend it
		for _,srcChild = range childNodes {
			b.insertNode(srcChild, destNode, addColumn)
		}
	} else {
		for _,srcChild = range childNodes {
			// TODO: Potentially improve speed by skipping column nodes. I suspect we will have already added those.
			// try to find the child node in the next level of the tree
			found := false
			for _,destChild:= range destNode.getChildNodes() {
				if destChild.Equals(srcChild) {
					// found a matching child node, recurse
					b.mergeNode(srcChild, destChild, addColumn)
					found = true
					break
				}
			}
			if !found {
				// Add the child node and stop
				b.insertNode(srcChild, destNode, addColumn)
				break
			}
		}
	}
}

func (b *sqlBuilder) insertNode(srcNode, parentNode NodeI, addColumn bool) {
	SetParentNode (srcNode, parentNode)
	b.assignAliases(srcNode, addColumn)
	if srcNode.nodeType() == REFERENCE_NODE {
		e := srcNode.(TableNodeI).EmbeddedNode_().(*ReferenceNode).relatedColumnNode()
		b.addNode(e, addColumn)
	}
}

// Walk DOWN the chain and assign aliases to the nodes found.
func (b *sqlBuilder) assignAliases (n NodeI, addColumn bool) {

	if n.GetAlias() == "" {
		// if it doesn't have a pre-assigned alias, give it an automated one
		if _,ok := n.(*ColumnNode); ok {
			key := "c" + b.subPrefix + strconv.Itoa(b.columnAliasNumber)
			b.columnAliasNumber++
			n.SetAlias(key)
		} else {
			key := "t" + b.subPrefix + strconv.Itoa(b.tableAliases.Len())
			n.SetAlias(key)
		}
	}

	if _,ok := n.(*ColumnNode); ok {
		if addColumn {
			b.columnAliases.Set(n.GetAlias(), n)
		}
	} else {
		b.tableAliases.Set(n.GetAlias(), n)
		if childNodes := n.getChildNodes(); childNodes != nil {
			for _, cn := range childNodes {
				b.assignAliases(cn, addColumn)
			}
		}
	}
}

// Add the columns that we will be selecting on to the node tree. Generally, we will add columns that are PKs along
// the chain of nodes so that we know how to unpack the objects.
// Specific situations that we do not automatically add the PK columns are: If its a distinct query, a count query,
// or we are a subquery.
func (b *sqlBuilder) makeColumnAliases() {
	if len(b.selects) > 0 {
		for _,n := range b.selects {
			b.addNode(n, true)
		}
		if !(b.distinct || b.isSubquery || b.isCount) {
			// Have some selects, so go through and make sure the ids are selected
			b.tableAliases.Range(func(key string, v interface{}) bool {
				node := v.(NodeI)
				n := node.(TableNodeI).PrimaryKeyNode_()
				b.addNode(n, true)
				return true
			})
		}
	} else {
		if b.isSubquery {
			// Subqueries must have specific columns called out. They might be as alias columns, so we do not panic here.
			if !(b.distinct || b.isCount) {
				// Still add id numbers so we can unpack this
				b.tableAliases.Range(func(key string, v interface{}) bool {
					node := v.(NodeI)
					n := node.(TableNodeI).PrimaryKeyNode_()
					b.addNode(n, true)
					return true
				})
			}
		} else {
			b.tableAliases.Range(func(key string, v interface{}) bool {
				node := v.(NodeI)
				selectNodes := node.(TableNodeI).SelectNodes_() // will add child nodes to node
				if selectNodes != nil {
					b.addNode(node, true) // will add the node and all child nodes
				}
				return true
			})

		}
	}

	if len(b.groupBys) > 0 {
		// SQL in general has a problem with group by items that are not selected, so we always select group by columns by implication
		// Some SQL forms have gotten aorund the problem by just choosing a random result, but modern SQL engines now consider this an error
		for _,n := range b.groupBys {
			b.addNode(n, true)
		}
	}
}


// After the intention of the query is gathered, this will add various nodes to the node tree to establish
// where the joins are. This process also adds aliases to the nodes.
func (b *sqlBuilder) buildNodeTree() {
	for _, n := range b.nodes() {
		b.addNode(n, false)
	}
}

// Returns the nodes referred to in the query. Some nodes will be container nodes, and so will have nodes
// inside them, but every node is either referred to, or contained in the returned nodes.
func (b *sqlBuilder) nodes() []NodeI {
	var nodes = []NodeI{}
	for _,n := range b.joins {
		nodes = append(nodes, n)
		if cn, ok := n.(TableNodeI).EmbeddedNode_().(conditioner); ok {
			if cn.getCondition() != nil {
				nodes = append(nodes, cn.getCondition())
			}
		}
	}
	nodes = append (nodes, b.orderBys...)

	if b.condition != nil {
		nodes = append (nodes, b.condition)
	}
	nodes = append (nodes, b.groupBys...)

	if b.having != nil {
		nodes = append (nodes, b.having)
	}
	nodes = append (nodes, b.selects...)


	b.aliasNodes.Range(func(key string, value interface{}) bool {
		nodes = append (nodes, value.(NodeI))
		return true
	})

	return nodes
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

	var oMap *types.OrderedMap = types.NewOrderedMap()
	aliasMap := types.NewOrderedMap()

	// First we create a tree structure of the data that will mirror the node structure
	for _,row := range rows {
		b.unpackObject(b.rootNode, row, oMap)
		b.unpackSpecialAliases(b.rootNode, row, aliasMap)
	}

	// We then walk the tree and create the final data structure as arrays
	oMap.Range(func(key string, value interface{}) bool {
		// Duplicate rows that are part of a join that is not an array join
		out2 := b.expandNode(b.rootNode, value.(ValueMap))
		// Add the Alias calculations specifically requested by the caller
		for _, o2 = range out2 {
			if m := aliasMap.Get(key); m != nil {
				o2[AliasResults] = m
			}
			out = append (out, o2)
		}
		return true
	})
	return out
}

// unpackObject finds the object that corresponds to parent in the row, and either adds it to the oMap, or it its
// already in the oMap, reuses the old one and adds more data to it. oMap should only contain objects of parent type.
func (b *sqlBuilder) unpackObject(parent NodeI, row ValueMap, oMap *types.OrderedMap) {
	var obj ValueMap
	var arrayKey string
	var currentArray *types.OrderedMap
	var iface interface{}
	var childNode NodeI
	var childTableNode TableNodeI
	var ok bool
	var objNode objectNode

	pk := b.makeObjectKey(parent, row)
	if pk == "" {
		// pk of object was not found in the row. This would happen for two reasons: 1) Object was not created because of a conditional expansion that failed, or
		// 2) This is a distinct select, and we are not selecting pks to avoid affecting the results of the query, in which case we will likely need to make up some ids
		if !b.distinct {
			return
		} else {
			b.distinctId++
			pk = strconv.Itoa(b.distinctId)
		}
	}

	if curObj := oMap.Get(pk); curObj != nil {
		obj = curObj.(ValueMap)
	} else {
		obj = NewValueMap()
		oMap.Set(pk, obj)
	}

	for _,childNode = range parent.getChildNodes() {
		if childTableNode,ok = childNode.(TableNodeI); ok {
			// if this is an embedded object, collect a group of objects
			objNode,_ = childTableNode.EmbeddedNode_().(objectNode)
			arrayKey = objNode.objectName()
			if 	iface,ok = obj[arrayKey]; !ok {
				// If this is the first time, create the group
				newArray := types.NewOrderedMap()
				obj[arrayKey] = newArray
				b.unpackObject(childNode, row, newArray)
			} else {
				// Already have a group, so add to the group
				currentArray = iface.(*types.OrderedMap)
				b.unpackObject(childNode, row, currentArray)
			}
		} else {
			b.unpackLeaf(childNode, row, obj)
		}
 	}
}

func (b *sqlBuilder) unpackLeaf(n NodeI, row ValueMap, obj ValueMap) {
	var key string
	var fieldName string

	switch node := n.(type) {
	case *ColumnNode:
		key = node.GetAlias()
		if b.columnAliases.Has(key) && !b.aliasNodes.Has(key) {	// could be a special alias, which we should unpack differently
			fieldName = node.dbColumn
			obj[fieldName] = row[key]
		}
	default:
		panic("Unexpected node type.")
	}
}

func (b *sqlBuilder) makeObjectKey(n NodeI, row ValueMap) string {
	var alias interface{}

	pkNode := b.getPkNode(n)
	if alias,_ = row[pkNode.GetAlias()]; alias == nil {
		return ""
	}
	pk := fmt.Sprint(alias)
	return pkNode.GetAlias() + "." + pk
}

// Returns the primary key value corresponding to the
func (b *sqlBuilder) getPkNode(n NodeI) NodeI {
	tn,ok := n.(TableNodeI)

	if !ok {
		return nil
	}
	pk := b.findMatchingChildNode(tn.PrimaryKeyNode_(), n)
	return pk
}

//
func (b *sqlBuilder) findMatchingChildNode(n NodeI, parent NodeI) (match NodeI) {
	var childNodes []NodeI

	if childNodes = parent.getChildNodes(); childNodes != nil {
		for _,cn := range childNodes {
			if cn.Equals(n) {
				return cn
			}
		}
	}
	return nil
}

// Craziness of handling situation where an array node wants to be individually expanded.
func (b *sqlBuilder) expandNode(n NodeI, nodeObject ValueMap) (outArray []ValueMap) {
	var childNode NodeI
	var item ValueMap
	var innerNodeObject ValueMap
	var copies []ValueMap
	var innerCopies []ValueMap
	var newArray []ValueMap
	var nodeCopy ValueMap

	if n.getChildNodes() == nil || len(n.getChildNodes()) == 0 {
		return
	}

	outArray = append(outArray, NewValueMap())

	for _,childNode = range n.getChildNodes()  {
		copies = []ValueMap{}
		for _, item = range outArray {
			switch node:=childNode.(type) {
			case *ColumnNode:
				item[node.dbColumn] = nodeObject[node.dbColumn]
			case TableNodeI:
				switch tableNode := node.EmbeddedNode_().(type) {
				case *ReferenceNode:
					// Should be a one or zero item array here
					om := nodeObject[tableNode.goName].(*types.OrderedMap)
					if om.Len() > 1 {
						panic ("Cannot have an array with more than one item here.")
					} else if om.Len() == 1 {
						innerNodeObject = nodeObject[tableNode.goName].(*types.OrderedMap).GetAt(0).(ValueMap)
						innerCopies = b.expandNode(childNode, innerNodeObject)
						if len(innerCopies) > 1 {
							for _, cp2 := range innerCopies {
								nodeCopy := item.Copy().(ValueMap)
								nodeCopy[tableNode.goName] = cp2
								copies = append(copies, nodeCopy)
							}
						} else {
							item[tableNode.goName] = map[string]interface{}(innerCopies[0])
						}
					}
					// else we likely were not included because of a conditional join
				case *ReverseReferenceNode:
					if !tableNode.isArray { // unique reverse or single expansion many
						newArray = []ValueMap{}
						nodeObject[tableNode.goName].(*types.OrderedMap).Range(func(key string, value interface{}) bool {
							innerNodeObject = value.(ValueMap)
							innerCopies = b.expandNode(childNode, innerNodeObject)
							for _,ic := range innerCopies {
								newArray = append(newArray, ic)
							}
							return true
						})
						for _,cp2 := range newArray {
							nodeCopy = item.Copy().(ValueMap)
							nodeCopy[tableNode.goName] = cp2
							copies = append(copies, nodeCopy)
						}
					} else {
						// From this point up, we should not be creating additional copies, since from this point down, we
						// are gathering an array
						newArray = []ValueMap{}
						nodeObject[tableNode.goName].(*types.OrderedMap).Range(func(key string, value interface{}) bool {
							innerNodeObject = value.(ValueMap)
							innerCopies = b.expandNode(childNode, innerNodeObject)
							for _,ic := range innerCopies {
								newArray = append(newArray, ic)
							}
							return true
						})
						item[tableNode.goName] = newArray
					}

				case *ManyManyNode:
					if tableNode.isTypeTable {
						intArray := []uint{}
						nodeObject[tableNode.goName].(*types.OrderedMap).Range(func(key string, value interface{}) bool {
							innerNodeObject = value.(ValueMap)
							typeKey := innerNodeObject[node.PrimaryKeyNode_().dbColumn]
							switch v := typeKey.(type) {
							case uint:
								intArray = append(intArray, v)
							case int:
								intArray = append(intArray, uint(v))
							}
							return true
						})
						if tableNode.isArray { // single expansion many
							item[tableNode.goName] = intArray
						} else {
							for _, cp2 := range intArray {
								nodeCopy = item.Copy().(ValueMap)
								nodeCopy[tableNode.goName] = []uint{cp2}
								copies = append(copies, nodeCopy)
							}
						}

					} else {
						newArray = []ValueMap{}
						nodeObject[tableNode.goName].(*types.OrderedMap).Range(func(key string, value interface{}) bool {
							innerNodeObject = value.(ValueMap)
							innerCopies = b.expandNode(childNode, innerNodeObject)
							for _, ic := range innerCopies {
								newArray = append(newArray, ic)
							}
							return true
						})
						if tableNode.isArray {
							item[tableNode.goName] = newArray
						} else {
							for _, cp2 := range newArray {
								nodeCopy = item.Copy().(ValueMap)
								nodeCopy[tableNode.goName] = []ValueMap{cp2}
								copies = append(copies, nodeCopy)
							}
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

func (b *sqlBuilder) unpackSpecialAliases(rootNode NodeI, row ValueMap, aliasMap *types.OrderedMap) {
	var obj ValueMap

	pk := b.makeObjectKey(rootNode, row)
	if pk == "" {
		return	// object was not found in the row
	}

	if curObj := aliasMap.Get(pk); curObj != nil {
		return // already added these to the row
	} else {
		obj = NewValueMap()
	}

	b.aliasNodes.Range(func(key string, value interface{}) bool {
		obj[key] = row[key]
		return true
	})

	if len(obj) > 0 {
		aliasMap.Set(pk, obj)
	}
}


type ValueMap map[string]interface{}

func NewValueMap() ValueMap {
	return make(ValueMap)
}

// Support the deep copy interface
func (m ValueMap) Copy() interface{} {
	vm := ValueMap{}
	for k,v := range m {
		if c, ok := v.(Copier); ok {
			v = c.Copy()
		}
		vm[k] = v
	}
	return vm
}