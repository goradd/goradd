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

type Generator interface {
	Type() string
	NewFunc() string
	Import() string
	SupportsColumn(col *db.ColumnDescription) bool
	ConnectorParams() *types.OrderedMap
	GenerateCreate(col *db.ColumnDescription) string
	GenerateLoad(ctrlName string, objName string, col *db.ColumnDescription) string
	GenerateSave(ctrlName string, objName string, col *db.ColumnDescription) string
}

type ControlRegistryKey struct {
	imp string
	typ string
}

var registry map[ControlRegistryKey]Generator

func RegisterGenerator(c Generator) {
	if registry == nil {
		registry = make(map[ControlRegistryKey]Generator)
	}

	e := ControlRegistryKey{c.Import(), c.Type()}
	registry[e] = c
}

func GetGenerator(imp string, typ string) Generator {
	e := ControlRegistryKey{imp, typ}

	d,_ := registry[e]
	return d
}