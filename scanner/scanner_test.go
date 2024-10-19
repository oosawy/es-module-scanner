package scanner

import (
	"testing"
)

func TestScan(t *testing.T) {
	input := "import defaultExport from \"module-name\";\n"

	expected := struct {
		imports []string
		exports []string
	}{
		[]string{"import defaultExport from \"module-name\";"},
		[]string{},
	}

	module := Scan(input)

	for i, imp := range expected.imports {
		if imp != module.Imports[i] {
			t.Errorf("Expected `%s`, got `%s`", module.Imports[i], imp)
		}
	}

	for i, exp := range expected.exports {
		if exp != module.Exports[i] {
			t.Errorf("Expected `%s`, got `%s`", module.Exports[i], exp)
		}
	}
}
