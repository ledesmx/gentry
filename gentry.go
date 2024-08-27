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

func if_error_exit(err error) {
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
	var (
		script_selected string
		name            string
		comment         string
		// version         float64
		// exec            string
		// terminal        bool
		// program_type    string
		// categories      string
	)

	root, err := os.Getwd()
	if_error_exit(err)

	files, err := os.ReadDir(root)
	if_error_exit(err)

	scripts, err := get_scripts(files, root)
	if_error_exit(err)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose the script file").
				Options(
					huh.NewOptions[string](scripts...)...,
				).
				Value(&script_selected),
			huh.NewInput().
				Title("Program name").
				Value(&name).
				Validate(func(str string) error {
					if len(str) == 0 {
						return errors.New("Please provide a name")
					}
					return nil
				}),
			huh.NewInput().
				Title("Comment").
				Value(&comment).
				Validate(func(str string) error {
					if len(str) == 0 {
						return errors.New("Please provide a comment")
					}
					return nil
				}),
			// huh.NewInput().
			// 	Title("Version").
			// 	Value(&version).
			// 	Validate(),
		),
	)

	error := form.Run()
	if_error_exit(error)

	fmt.Printf("Your selection is %s\n", script_selected)
	fmt.Printf("Program name %s\n", name)
	fmt.Printf("Comment %s\n", comment)
}
