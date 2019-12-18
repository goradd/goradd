package generator

import (
	"github.com/goradd/goradd/pkg/page"
	"path"
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
	SupportsColumn(ref interface{}) bool
	GenerateCreator(ref interface{}, desc *ControlDescription) string
	GenerateRefresh(ref interface{}, desc *ControlDescription) string
	GenerateUpdate(ref interface{}, desc *ControlDescription) string
}

type ProviderGenerator interface {
	GenerateProvider(ref interface{}, desc *ControlDescription) string
}

var controlGeneratorRegistry map[string]ControlGenerator

func RegisterControlGenerator(c ControlGenerator) {
	if controlGeneratorRegistry == nil {
		controlGeneratorRegistry = make(map[string]ControlGenerator)
	}

	i := c.Imports()
	e := path.Join(i[0].Path, c.Type())
	controlGeneratorRegistry[e] = c
}

/*
func RegisterControl(c page.ControlI) {
	if controlGeneratorRegistry == nil {
		controlGeneratorRegistry = make(map[ControlGeneratorRegistryKey]ControlGenerator)
	}


	i := c.Imports()
	e := ControlGeneratorRegistryKey{i[0].Path, c.Type()}
	controlGeneratorRegistry[e] = c
}
*/

func GetControlGenerator(imp string, typ string) ControlGenerator {
	e := path.Join(imp, typ)

	d, _ := controlGeneratorRegistry[e]
	return d
}
