package connector

import (
	"github.com/spekary/goradd/page"
	"github.com/spekary/goradd/orm/db"
	"github.com/spekary/goradd/util/types"
)

type ControlType int

const (
	ControlTypeInteger ControlType = iota + 1
)

type ConnectorParam struct {
	Name string
	Description string
	Typ ControlType
	Template string
	DoFunc func(c page.ControlI, val interface{})
}

type Describer interface {
	Type() string
	NewFunc() string
	Import() string
	SupportsColumn(col db.ColumnDescription) bool
	ConnectorParams() *types.OrderedMap
}

type ControlRegistryKey struct {
	imp string
	typ string
}

var registry map[ControlRegistryKey]Describer

func RegisterControl(c Describer) {
	if registry == nil {
		registry = make(map[ControlRegistryKey]Describer)
	}

	e := ControlRegistryKey{c.Import(), c.Type()}
	registry[e] = c
}