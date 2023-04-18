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
	SupportsColumn(ref interface{}) bool
	GenerateCreator(ref interface{}, desc *ControlDescription) string
	GenerateRefresh(ref interface{}, desc *ControlDescription) string
	GenerateUpdate(ref interface{}, desc *ControlDescription) string
	GenerateModifies(ref interface{}, desc *ControlDescription) string
}

type ProviderGenerator interface {
	GenerateProvider(ref interface{}, desc *ControlDescription) string
}

var controlGeneratorRegistry map[string]ControlGenerator

func RegisterControlGenerator(c ControlGenerator, path string) {
	if controlGeneratorRegistry == nil {
		controlGeneratorRegistry = make(map[string]ControlGenerator)
	}

	controlGeneratorRegistry[path] = c
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

func GetControlGenerator(controlPath string) ControlGenerator {
	d, _ := controlGeneratorRegistry[controlPath]
	return d
}
