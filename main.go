package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

type Data struct {
	filePath     string
	newName      string
	newDirectory string
	merging      []string
}

func main() {
	filePath := "/home/leoyenet/Pictures"
	data := &Data{}
	directory := huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Title("Which folder do you want to join in").
				Description("It will have this parent").
				CurrentDirectory(filePath).
				ShowHidden(true).
				ShowPermissions(false).
				FileAllowed(false).
				DirAllowed(true).
				Value(&data.filePath),
		),
	).WithHeight(10)

	if err := directory.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Println(GetOnlyFolder(data.filePath))
	folders, err2 := GetOnlyFolder(data.filePath)
	if err2 != nil {
		os.Exit(2)
	}
	var options []huh.Option[string]
	for _, value := range folders {
		options = append(options, huh.NewOption[string](filepath.Base(value), value))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("which folders to merge").
				Options(huh.NewOptions(folders...)...).
				Value(&data.merging),
			huh.NewInput().
				Title("What should new Name of files be").
				Value(&data.newName),
			huh.NewInput().
				Title("What should directory name be").
				Value(&data.newDirectory),
		),
	).WithTheme(huh.ThemeCatppuccin())
	err := form.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
	allFiles := []string{}

	for _, value := range data.merging {
		files, err2 := GetOnlyFilesInFolder(value)
		if err2 != nil {
			continue
		}
		allFiles = append(allFiles, files...)
	}
	fmt.Println(allFiles)
	fmt.Println(len(allFiles))
	newPath := filepath.Join(data.filePath, data.newDirectory)
	if e := os.Mkdir(newPath, 0775); e != nil {
		fmt.Println(e)
		os.Exit(3)
	}
	for i, oldFileName := range allFiles {
		newFilePath := filepath.Join(newPath, fmt.Sprintf("%v_%d_%v", data.newName, i, filepath.Ext(oldFileName)))
		fmt.Printf("\ni: %d\nOld: %v\nNew: %v\n", i, oldFileName, newFilePath)
		os.Rename(oldFileName, newFilePath)
	}

	for _, value := range folders {
		os.Remove(value)
	}

}

func GetOnlyFilesInFolder(folderPath string) ([]string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", folderPath, err)
	}

	var filePaths []string
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			fullPath := filepath.Join(folderPath, entry.Name())
			filePaths = append(filePaths, fullPath)
		}
	}

	return filePaths, nil
}

func GetOnlyFolder(folderPath string) ([]string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", folderPath, err)
	}

	var filePaths []string
	for _, entry := range entries {
		if entry.Type().IsDir() {
			fullPath := filepath.Join(folderPath, entry.Name())
			filePaths = append(filePaths, fullPath)
		}
	}

	return filePaths, nil
}
