package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
)

const Script_Extension = ".sh"
const Entry_Path = "/usr/share/applications"

var Categories = []string{
	"AudioVideo",
	"Audio",
	"Video",
	"Development",
	"Education",
	"Game",
	"Graphics",
	"Network",
	"Office",
	"Science",
	"Settings",
	"System",
	"Utility",
}

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

func format_categories(categories []string) string {
	var cats string
	for i := 0; i < len(categories); i++ {
		cats += fmt.Sprintf("%v;", categories[i])
	}
	return cats
}

func main() {
	var (
		script_selected string
		name            string
		comment         string
		version         string
		terminal        string
		categories_list []string
	)

	root, err := os.Getwd()
	if_error_exit(err)

	files, err := os.ReadDir(root)
	if_error_exit(err)

	// for i := 0; i < len(files); i++ {
	// 	fmt.Println(files[i])
	// }

	scripts, err := get_scripts(files, root)
	if_error_exit(err)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Script file *").
				Options(
					huh.NewOptions(scripts...)...,
				).
				Value(&script_selected),
			huh.NewInput().
				Title("Program name *").
				Value(&name).
				Validate(func(str string) error {
					if len(str) == 0 {
						return errors.New("Please provide a name")
					}
					return nil
				}),
			huh.NewInput().
				Title("Comment *").
				Value(&comment).
				Validate(func(str string) error {
					if len(str) == 0 {
						return errors.New("Please provide a comment")
					}
					return nil
				}),
			huh.NewInput().
				Title("Version").
				Value(&version),
			huh.NewSelect[string]().
				Title("Terminal *").
				Options(
					huh.NewOption("true", "true"),
					huh.NewOption("false", "false"),
				).
				Value(&terminal),
			huh.NewMultiSelect[string]().
				Title("Categories").
				Options(
					huh.NewOptions(Categories...)...,
				).
				Value(&categories_list),
		),
	)

	error := form.Run()
	if_error_exit(error)

	categories := format_categories(categories_list)
	file_name, _ := strings.CutSuffix(script_selected, ".sh")
	file := fmt.Sprintf("%s/gentry.%s.desktop", Entry_Path, file_name)

	fmt.Printf("Generating the following desktop entry. \n")
	fmt.Printf("%s \n\n", file)
	fmt.Printf("############### \n")
	fmt.Printf("[Desktop Entry] \n")
	fmt.Printf("Type=Application \n")
	fmt.Printf("Version=%s \n", version)
	fmt.Printf("Name=%s \n", name)
	fmt.Printf("Comment=%s \n", comment)
	fmt.Printf("Path=%s \n", root)
	fmt.Printf("Exec=%s \n", script_selected)
	fmt.Printf("Terminal=%s \n", terminal)
	fmt.Printf("Categories=%s \n", categories)
	fmt.Printf("############## \n")

}
