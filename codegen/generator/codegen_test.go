package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImportAliases(t *testing.T) {
	codegen := CodeGenerator{}
	codegen.ResetImports()
	codegen.AddImportPaths("a/c", "b/c", "d/c", "d/e")
	assert.Equal(t, "c3", codegen.ImportPackage("d/c"))
}
