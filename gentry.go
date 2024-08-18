package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

const Script_Extension = ".sh"

func manage_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func get_scripts(files []fs.DirEntry, path string) ([]string, error) {
	var scripts []string
	for _, v := range files {
		if !v.IsDir() {
			name := v.Name()
			extension := filepath.Ext(fmt.Sprintf("%s/%s", path, name))
			if extension == Script_Extension {
				scripts = append(scripts, name)
			}
		}
	}
	if len(scripts) == 0 {
		return scripts, errors.New("No script found")
	}
	return scripts, nil
}

func main() {
	var script_selected string

	root, err := os.Getwd()
	manage_error(err)

	files, err := os.ReadDir(root)
	manage_error(err)

	scripts, err := get_scripts(files, root)
	manage_error(err)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose the script file").
				Options(
					huh.NewOptions[string](scripts...)...,
				).
				Value(&script_selected),
		),
	)

	error := form.Run()
	manage_error(error)

	fmt.Printf("Your selection is %s", script_selected)
}
