package main

import scanner "github.com/oosawy/es-module-scanner/scanner"

func main() {

	input := `import defaultExport from "module-name";`

	scanner.Scan(input)
}
