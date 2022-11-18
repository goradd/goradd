package sql

import (
	"fmt"
	. "github.com/goradd/goradd/pkg/orm/query"
	"strings"
)

type OperationSqler interface {
	OperationSql(op Operator, operandStrings []string) string
}

type DeleteUsesAliaser interface {
	DeleteUsesAlias() bool
}

// Generator is an aid to generating various sql statements.
// SQL dialects are similar, but have small variations. This object
// attempts to handle the major issues, while allowing individual
// implementations of SQL to do their own tweaks.
type Generator struct {
	b       *Builder
	argList []any
}

func NewGenerator(builder *Builder) *Generator {
	return &Generator{b: builder}
}

func (g *Generator) iq(v string) string {
	return g.b.db.QuoteIdentifier(v)
}

func (g *Generator) addArg(v any) string {
	g.argList = append(g.argList, v)
	return g.b.db.FormatArgument(len(g.argList))
}

func (g *Generator) generateSelectSql() (sql string) {
	if g.b.IsDistinct {
		sql = "SELECT DISTINCT\n"
	} else {
		sql = "SELECT\n"
	}

	sql += g.generateColumnListWithAliases()
	sql += g.generateFromSql()
	sql += g.generateWhereSql()
	sql += g.generateGroupBySql()
	sql += g.generateHaving()
	sql += g.generateOrderBySql()
	sql += g.generateLimitSql()
	return
}

func (g *Generator) generateDeleteSql() (sql string) {
	if t, ok := g.b.db.(DeleteUsesAliaser); ok && t.DeleteUsesAlias() {
		j := g.b.RootJoinTreeItem
		alias := g.iq(j.Alias)
		sql = "DELETE " + alias + "\n"
	} else {
		sql = "DELETE\n"
	}

	sql += g.generateFromSql()
	sql += g.generateWhereSql()
	sql += g.generateOrderBySql()
	sql += g.generateLimitSql()

	return
}

func (g *Generator) generateColumnListWithAliases() (sql string) {
	g.b.ColumnAliases.Range(func(key string, j *JoinTreeItem) bool {
		sql += g.generateColumnNodeSql(j.Parent.Alias, j.Node) + " AS " + g.iq(key) + ",\n"
		return true
	})

	if g.b.AliasNodes != nil {
		g.b.AliasNodes.Range(func(key string, v Aliaser) bool {
			node := v.(NodeI)
			aliaser := v.(Aliaser)
			sql += g.generateNodeSql(node, false)
			alias := aliaser.GetAlias()
			if alias != "" {
				// This happens in a subquery
				sql += " AS " + g.iq(alias)
			}
			sql += ",\n"
			return true
		})
	}

	sql = strings.TrimSuffix(sql, ",\n")
	sql += "\n"
	return
}

// Generate the column node sql.
func (g *Generator) generateColumnNodeSql(parentAlias string, node NodeI) (sql string) {
	return g.iq(parentAlias) + "." + g.iq(ColumnNodeDbName(node.(*ColumnNode)))
}

func (g *Generator) generateNodeSql(n NodeI, useAlias bool) (sql string) {
	switch node := n.(type) {
	case *ValueNode:
		v := ValueNodeGetValue(node)
		if a, ok := v.([]NodeI); ok {
			// value is actually a list of nodes
			var l []string
			for _, o := range a {
				l = append(l, g.generateNodeSql(o, useAlias))
			}
			return strings.Join(l, ",")
		} else {
			return g.addArg(v)
		}
	case *OperationNode:
		return g.generateOperationSql(node, useAlias)
	case *ColumnNode:
		item := g.b.GetItemFromNode(node)
		if useAlias {
			sql = g.generateAlias(item.Alias)
		} else {
			sql = g.generateColumnNodeSql(item.Parent.Alias, node)
		}
	case *AliasNode:
		sql = g.iq(node.GetAlias())
	case *SubqueryNode:
		sql = g.generateSubquerySql(node)
	case TableNodeI:
		tj := g.b.GetItemFromNode(node)
		sql = g.generateColumnNodeSql(tj.Alias, node.PrimaryKeyNode())
	default:
		panic("Can't generate sql from node type.")
	}
	return
}

func (g *Generator) generateOperationSql(n *OperationNode, useAlias bool) (sql string) {
	if useAlias && n.GetAlias() != "" {
		sql = g.iq(n.GetAlias())
		return
	}

	var operands []string
	for _, o := range OperationNodeOperands(n) {
		operands = append(operands, g.generateNodeSql(o, useAlias))
	}

	if o, ok := g.b.db.(OperationSqler); ok {
		sql = o.OperationSql(OperationNodeOperator(n), operands)
		if sql != "" {
			return sql
		}
	}

	switch OperationNodeOperator(n) {
	case OpFunc:
		if len(operands) > 0 {
			sql = strings.Join(operands, ",")
		} else {
			if OperationNodeFunction(n) == "COUNT" {
				sql = "*"
			}
		}

		if OperationNodeDistinct(n) {
			sql = "DISTINCT " + sql
		}
		sql = OperationNodeFunction(n) + "(" + sql + ") "

	case OpNull:
		fallthrough
	case OpNotNull:
		s := operands[0]
		sql = s + " IS " + OperationNodeOperator(n).String()
		sql = "(" + sql + ") "

	case OpNot:
		s := operands[0]
		sql = OperationNodeOperator(n).String() + " " + s
		sql = "(" + sql + ") "

	case OpIn:
		fallthrough
	case OpNotIn:
		s := operands[0]
		sql = s + " " + OperationNodeOperator(n).String()
		sql += " (" + operands[1] + ") "

	case OpAll:
		fallthrough
	case OpNone:
		sql = "(" + OperationNodeOperator(n).String() + ") "
	case OpStartsWith:
		// SQL supports this with a LIKE operation
		s := operands[0]
		v := operands[1]
		v += "%"
		sql = fmt.Sprintf(`(%s LIKE %s)`, s, v)
	case OpEndsWith:
		// SQL supports this with a LIKE operation
		s := operands[0]
		v := operands[1]
		v = "%" + v
		sql = fmt.Sprintf(`(%s LIKE %s)`, s, v)
	case OpContains:
		// SQL supports this with a LIKE operation
		s := operands[0]
		v := operands[1]
		v = "%" + v + "%"
		sql = fmt.Sprintf(`(%s LIKE %s)`, s, v)

	case OpDateAddSeconds:
		// Modifying a datetime in the query
		// Only works on date, datetime and timestamps. Not times.
		// This is highly SQL dialect dependent. Use the below as an example.
		// Default is to not implement it. To implement it, it must be overridden in the database implementation.

		panic("DateAddSeconds is not implemented in this database")
		/*
			s := operands[0]
			s2 := operands[1]
			sql = fmt.Sprintf(`DATE_ADD(%s, INTERVAL (%s) SECOND)`, s, s2)
		*/
	case OpXor:
		// Some sqls do not have an XOR operator, so we manually implement the code here
		// Override in the database implementation if XOR is implemented
		s := operands[0]
		s2 := operands[1]
		sql = fmt.Sprintf(`(((%[1]s) AND NOT (%[2]s)) OR (NOT (%[1]s) AND (%[2]s)))`, s, s2)

	default:
		sOp := " " + OperationNodeOperator(n).String() + " "
		sql = " (" + strings.Join(operands, sOp) + ") "
	}
	return
}

func (g *Generator) generateAlias(alias string) (sql string) {
	return g.iq(alias)
}

func (g *Generator) generateSubquerySql(node *SubqueryNode) (sql string) {
	// The copy below intentionally reuses the argList and db items
	g2 := *g
	g2.b = SubqueryBuilder(node).(*Builder)
	sql = g2.generateSelectSql()
	sql = "(" + sql + ")"
	return
}

func (g *Generator) generateFromSql() (sql string) {
	sql = "FROM\n"

	j := g.b.RootJoinTreeItem
	sql += g.iq(NodeTableName(j.Node)) + " AS " + g.iq(j.Alias) + "\n"

	for _, child := range j.ChildReferences {
		sql += g.generateJoinSql(child)
	}
	return
}

func (g *Generator) generateJoinSql(j *JoinTreeItem) (sql string) {
	var tn TableNodeI
	var ok bool

	if tn, ok = j.Node.(TableNodeI); !ok {
		return
	}

	switch node := tn.EmbeddedNode_().(type) {
	case *ReferenceNode:
		sql = "LEFT JOIN "
		sql += g.iq(ReferenceNodeRefTable(node)) + " AS " +
			g.iq(j.Alias) + " ON " + g.iq(j.Parent.Alias) + "." +
			g.iq(ReferenceNodeDbColumnName(node)) + " = " + g.iq(j.Alias) + "." + g.iq(ReferenceNodeRefColumn(node))
		if j.JoinCondition != nil {
			s := g.generateNodeSql(j.JoinCondition, false)
			sql += " AND " + s
		}
	case *ReverseReferenceNode:
		if g.b.LimitInfo != nil && ReverseReferenceNodeIsArray(node) {
			panic("We do not currently support limited queries with an array join.")
		}

		sql = "LEFT JOIN "
		sql += g.iq(ReverseReferenceNodeRefTable(node)) + " AS " +
			g.iq(j.Alias) + " ON " + g.iq(j.Parent.Alias) + "." +
			g.iq(ReverseReferenceNodeKeyColumnName(node)) + " = " + g.iq(j.Alias) + "." + g.iq(ReverseReferenceNodeRefColumn(node))
		if j.JoinCondition != nil {
			s := g.generateNodeSql(j.JoinCondition, false)
			sql += " AND " + s
		}
	case *ManyManyNode:
		if g.b.LimitInfo != nil {
			panic("We do not currently support limited queries with an array join.")
		}

		sql = "LEFT JOIN "

		sql += g.iq(ManyManyNodeDbTable(node)) + " AS " + g.iq(j.Alias+"a") + " ON " +
			g.iq(j.Parent.Alias) + "." +
			g.iq(ColumnNodeDbName(ParentNode(node).(TableNodeI).PrimaryKeyNode())) +
			" = " + g.iq(j.Alias+"a") + "." + g.iq(ManyManyNodeDbColumn(node)) + "\n"
		sql += "LEFT JOIN " + g.iq(ManyManyNodeRefTable(node)) + " AS " + g.iq(j.Alias) +
			" ON " + g.iq(j.Alias+"a") + "." + g.iq(ManyManyNodeRefColumn(node)) +
			" = " + g.iq(j.Alias) + "." + g.iq(ManyManyNodeRefPk(node))

		if j.JoinCondition != nil {
			s := g.generateNodeSql(j.JoinCondition, false)
			sql += " AND " + s
		}
	default:
		return
	}
	sql += "\n"
	for _, cj := range j.ChildReferences {
		s := g.generateJoinSql(cj)
		sql += s
	}
	return
}

func (g *Generator) generateWhereSql() (sql string) {
	if g.b.ConditionNode != nil {
		sql = "WHERE "
		var s string
		s = g.generateNodeSql(g.b.ConditionNode, false)
		sql += s + "\n"
	}
	return
}

func (g *Generator) generateGroupBySql() (sql string) {
	if g.b.GroupBys != nil && len(g.b.GroupBys) > 0 {
		sql = "GROUP BY "
		for _, n := range g.b.GroupBys {
			s := g.generateNodeSql(n, true)
			sql += s + ","
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "\n"
	}
	return
}

// Note that some SQLs (MySQL, SqlLite) allow the use of aliases in a having clause,
// and some (Postgres) do not. We might need to check for this at some point.
func (g *Generator) generateHaving() (sql string) {
	if g.b.HavingNode != nil {
		sql = "HAVING "
		var s string
		s = g.generateNodeSql(g.b.HavingNode, false)
		sql += s + "\n"
	}
	return
}

func (g *Generator) generateLimitSql() (sql string) {
	if g.b.LimitInfo == nil {
		return ""
	}

	if g.b.LimitInfo.MaxRowCount > -1 {
		sql += fmt.Sprintf("LIMIT %d ", g.b.LimitInfo.MaxRowCount)
	}

	if g.b.LimitInfo.Offset > 0 {
		sql += fmt.Sprintf("OFFSET %d ", g.b.LimitInfo.Offset)
	}
	return
}

func (g *Generator) generateOrderBySql() (sql string) {
	if g.b.OrderBys != nil && len(g.b.OrderBys) > 0 {
		sql = "ORDER BY "
		for _, n := range g.b.OrderBys {
			s := g.generateNodeSql(n, true)
			if sorter, ok := n.(NodeSorter); ok {
				if NodeSorterSortDesc(sorter) {
					s += " DESC"
				}
			}
			sql += s + ","
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "\n"
	}
	return
}
