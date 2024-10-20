package scanner

import (
	"testing"
)

func TestScan(t *testing.T) {
	input := `import defaultExport from "module-name";
import * as name from "module-name";
import { export1 } from "module-name";
import { export1 as alias1 } from "module-name";
import { default as alias } from "module-name";
import { export1, export2 } from "module-name";
import { export1, export2 as alias2, /* … */ } from "module-name";
import { "string name" as alias } from "module-name";
import defaultExport, { export1, /* … */ } from "module-name";
import defaultExport, * as name from "module-name";
import "module-name";`

	expectedImports := []string{
		`import defaultExport from "module-name";`,
		`import * as name from "module-name";`,
		`import { export1 } from "module-name";`,
		`import { export1 as alias1 } from "module-name";`,
		`import { default as alias } from "module-name";`,
		`import { export1, export2 } from "module-name";`,
		`import { export1, export2 as alias2, /* … */ } from "module-name";`,
		`import { "string name" as alias } from "module-name";`,
		`import defaultExport, { export1, /* … */ } from "module-name";`,
		`import defaultExport, * as name from "module-name";`,
		`import "module-name";`,
	}

	expectedExports := []string{}

	module := Scan(input)

	for i, imp := range expectedImports {
		if len(module.Imports) <= i {
			t.Errorf("Expected `%s`, got nothing", imp)
		} else if imp != module.Imports[i] {
			t.Errorf("Expected `%s`, got `%s`", imp, module.Imports[i])
		}
	}

	for i, exp := range expectedExports {
		if len(module.Exports) <= i {
			t.Errorf("Expected `%s`, got nothing", exp)
		} else if exp != module.Exports[i] {
			t.Errorf("Expected `%s`, got `%s`", exp, module.Exports[i])
		}
	}
}
