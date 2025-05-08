package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

func main() {
	filePath := "/home/leoyenet/Pictures"
	firstDir := ""
	secondDir := ""
	fileName := ""
	DirName := ""
	deleteItems := []int{}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Title("Which folders to join").
				Description("It will have this parent").
				CurrentDirectory(filePath).
				ShowHidden(true).
				ShowPermissions(false).
				FileAllowed(false).
				DirAllowed(true).
				Value(&firstDir),
			huh.NewFilePicker().
				Title("Which folders to join").
				CurrentDirectory(filePath).
				ShowHidden(true).
				ShowPermissions(false).
				FileAllowed(false).
				DirAllowed(true).
				Value(&secondDir),
			huh.NewInput().
				Title("What should new Name of files be").
				Value(&fileName),
			huh.NewInput().
				Title("What should directory name be").
				Value(&DirName),
			huh.NewMultiSelect[int]().
				Title("Do you want to delete them").
				Options(
					huh.NewOption("Delete First", 0),
					huh.NewOption("Delete Second", 1),
				).
				Value(&deleteItems),
		),
	).WithTheme(huh.ThemeCatppuccin())
	err := form.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
	files := []string{}
	firstDirFiles, err1 := GetOnlyFilesInFolder(firstDir)
	secondDirFiles, err2 := GetOnlyFilesInFolder(secondDir)
	if err1 != nil || err2 != nil {
		fmt.Println(err1, err2)
	}

	files = append(files, firstDirFiles...)
	files = append(files, secondDirFiles...)
	newDir := filepath.Dir(firstDir)
	newPath := filepath.Join(newDir, DirName)
	err3 := os.MkdirAll(newPath, os.ModePerm)
	if err3 != nil {
		os.Exit(2)
	}
	err4 := error(nil)
	for i, oldFileName := range files {

		newFileName := filepath.Join(newPath, fmt.Sprintf("%s_%d%s", fileName, i+1, filepath.Ext(oldFileName)))
		fmt.Printf("Old: %v\nNew: %v\n", oldFileName, newFileName)
		err4 = os.Rename(oldFileName, newFileName)
		if err4 != nil {
			continue
		}

	}
	for _, value := range deleteItems {
		if value == 0 {
			os.Remove(firstDir)
		}
		if value == 1 {
			os.Remove(secondDir)
		}
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
