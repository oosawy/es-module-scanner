package scanner

import (
	"testing"
)

func TestScan(t *testing.T) {
	input := `import defaultExport from "module-name";
import * as name from "module-name";
import { export1 } from "module-name";
`

	expectedImports := []string{
		`import defaultExport from "module-name";`,
		`import * as name from "module-name";`,
		`import { export1 } from "module-name";`,
	}

	expectedExports := []string{}

	module := Scan(input)

	for i, imp := range expectedImports {
		if len(module.Imports) <= i {
			t.Errorf("Expected `%s`, got nothing", imp)
		} else if imp != module.Imports[i] {
			t.Errorf("Expected `%s`, got `%s`", module.Imports[i], imp)
		}
	}

	for i, exp := range expectedExports {
		if len(module.Exports) <= i {
			t.Errorf("Expected `%s`, got nothing", exp)
		} else if exp != module.Exports[i] {
			t.Errorf("Expected `%s`, got `%s`", module.Exports[i], exp)
		}
	}
}
