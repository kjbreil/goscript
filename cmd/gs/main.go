package main

import (
	"path/filepath"
)

func main() {
	runDir := "/Users/kjell/dev/goscript/run"
	filename := "config.yml"

	modules, err := parseConfig("config.yml")
	if err != nil {
		panic(err)
	}

	// err = setupProject(runDir, modules)
	// if err != nil {
	// 	panic(err)
	// }

	mainFile, err := genMain(filename, modules)
	if err != nil {
		panic(err)
	}

	err = mainFile.Save(filepath.Join(runDir, "main.go"))
	if err != nil {
		panic(err)
	}
}
