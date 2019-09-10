package generator

import (
	"github.com/goradd/goradd/pkg/page"
)

type ControlType int

type ConnectorParam struct {
	Name        string
	Description string
	Typ         ControlType
	Template    string
	DoFunc      func(c page.ControlI, val interface{})
}

type ImportPath struct {
	Alias string
	Path  string
}

type ControlGenerator interface {
	Type() string
	Imports() []ImportPath
	SupportsColumn(col *ColumnType) bool
	GenerateCreator(col *ColumnType) string
	GenerateRefresh(col *ColumnType) string
	GenerateUpdate(col *ColumnType) string
}

type ProviderGenerator interface {
	GenerateProvider(col *ColumnType) string
}

type ControlGeneratorRegistryKey struct {
	imp string
	typ string
}

var controlGeneratorRegistry map[ControlGeneratorRegistryKey]ControlGenerator

func RegisterControlGenerator(c ControlGenerator) {
	if controlGeneratorRegistry == nil {
		controlGeneratorRegistry = make(map[ControlGeneratorRegistryKey]ControlGenerator)
	}

	i := c.Imports()
	e := ControlGeneratorRegistryKey{i[0].Path, c.Type()}
	controlGeneratorRegistry[e] = c
}

func GetControlGenerator(imp string, typ string) ControlGenerator {
	e := ControlGeneratorRegistryKey{imp, typ}

	d, _ := controlGeneratorRegistry[e]
	return d
}
