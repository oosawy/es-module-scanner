package main

import (
	"fmt"

	scanner "github.com/oosawy/es-module-scanner/scanner"
)

func main() {

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

	mod := scanner.Scan(input)

	fmt.Printf("Imports: %v\n", mod.Imports)
	fmt.Printf("Exports: %v\n", mod.Exports)
}
