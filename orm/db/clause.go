package db

import (
)

type ClauseI interface {
	AddToBuilder(b QueryBuilderI)
}

type LimitClause struct {
	maxRowCount int64
	offset int64
}

func Limit(maxRowCount int64, offset int64) *LimitClause {
	return &LimitClause{maxRowCount, offset}
}

func (c *LimitClause) AddToBuilder(b QueryBuilderI) {
	b.Limit(c.maxRowCount, c.offset)
}

type ExpandClause struct {
	node NodeI
	conditions []NodeI
}

func Expand(node NodeI, conditions... NodeI) *ExpandClause {
	return &ExpandClause{node, conditions}
}

func (c *ExpandClause) AddToBuilder(b QueryBuilderI) {
	b.Join(c.node, c.conditions...)
}

type OrderByClause struct {
	nodes []NodeI
}

func OrderBy(nodes... NodeI) *OrderByClause {
	return &OrderByClause{nodes}
}

func (c *OrderByClause) AddToBuilder(b QueryBuilderI) {
	b.OrderBy(c.nodes...)
}

type SelectClause struct {
	nodes []NodeI
}

func Select(nodes... NodeI) *SelectClause {
	return &SelectClause{nodes}
}

func (c *SelectClause) AddToBuilder(b QueryBuilderI) {
	b.Select(c.nodes...)
}


