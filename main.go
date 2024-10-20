package main

import (
	"fmt"

	scanner "github.com/oosawy/es-module-scanner/scanner"
)

func main() {

	input := `import "module-name";
import defaultExport from "module-name";
import * as name from "module-name";
import { export1 } from "module-name";
`

	mod := scanner.Scan(input)

	fmt.Printf("Imports: %v\n", mod.Imports)
	fmt.Printf("Exports: %v\n", mod.Exports)
}
