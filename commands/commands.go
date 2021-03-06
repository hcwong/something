package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	mdtty "github.com/hcwong/golang-md-tty"
)

// Deploy deploys the notes page unto netlify
func Deploy() {
	dir, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
	if pathErr != nil {
		log.Printf("Deployment failed due to %s\n", pathErr)
		return
	}

	buildCmd := exec.Command("hugo", "--gc", "--minify")
	buildCmd.Dir = dir
	buildCmd.Stdout = os.Stdout
	if err := buildCmd.Run(); err != nil {
		log.Printf("Failed to build files for deployment because: %s\n", err)
		log.Println("Did you remember to install hugo?")
		return
	}

	deployCmd := exec.Command("netlify", "deploy", "--prod")
	deployCmd.Dir = dir
	deployCmd.Stdout = os.Stdout
	if err := deployCmd.Run(); err != nil {
		log.Printf("Failed to deploy because: %s\n", err)
		log.Println("Did you remember to install netlify?")
	}
}

// DeletePage delete files at the given path
func DeletePage(name string, fileType string) error {
	dir, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
	if pathErr != nil {
		return pathErr
	}

	var fileToDelete string
	if fileType == "notes" {
		fileToDelete = fmt.Sprintf("%s/%s/notes_all/%s.md", dir, fileType, name)
	} else {
		fileToDelete = fmt.Sprintf("%s/%s/%s.md", dir, fileType, name)
	}

	isPageExists, err := isFileExists(fileToDelete)
	if err != nil {
		return err
	}
	if !isPageExists {
		log.Printf("The %s page you want to delete does not exist, please first create it.\n", fileType)
		return nil
	}

	if err := delete(fileToDelete); err != nil {
		log.Println("Failed to delete page")
		return err
	}
	return nil
}

// GenPage generates a new page based on the one in templates, or creates a new one if the page is missing
func GenPage(name string, fileType string) error {
	dir, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
	if pathErr != nil {
		return pathErr
	}

	var newFile string
	if fileType == "notes" {
		newFile = fmt.Sprintf("%s/%s/notes_all/%s.md", dir, fileType, name)
	} else {
		newFile = fmt.Sprintf("%s/%s/%s.md", dir, fileType, name)
	}

	isPageExists, err := isFileExists(newFile)
	if err != nil {
		return err
	}
	if isPageExists {
		log.Printf("The %s page you want to create already exists. Edit it or delete it.\n", fileType)
		return nil
	}

	isTemplateExist, err := isFileExists(fmt.Sprintf("%s/templates/%s_template.md", dir, fileType))
	if err != nil {
		return err
	}

	if err := createFile(newFile); err != nil {
		return err
	}

	if isTemplateExist {
		// Create a new directory for notes. Basically notes got one extra step
		copy(fmt.Sprintf("%s/templates/%s_template.md", dir, fileType), newFile)
	}

	if vimErr := openFile(newFile); vimErr != nil {
		delete(newFile)
		return err
	}

	return nil
}

// EditPage allows us to edit a page
func EditPage(name string, fileType string) error {
	dir, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
	if pathErr != nil {
		return pathErr
	}

	var fileToEdit string
	if fileType == "notes" {
		fileToEdit = fmt.Sprintf("%s/%s/notes_all/%s.md", dir, fileType, name)
	} else {
		fileToEdit = fmt.Sprintf("%s/%s/%s.md", dir, fileType, name)
	}
	isPageExists, err := isFileExists(fileToEdit)
	if err != nil {
		return err
	}
	if !isPageExists {
		log.Printf("The %s page you want to edit does not exist, please first create it.\n", fileType)
		return nil
	}

	if vimErr := openFile(fileToEdit); vimErr != nil {
		delete(fileToEdit)
		return err
	}

	return nil
}

// Ls lists all the pages available
func Ls(fileType string) {
	dir, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
	if pathErr != nil {
		log.Println("ls command failed")
	}

	var path string
	if fileType == "notes" {
		path = "/notes_all"
	} else {
		path = ""
	}

	findPath := fmt.Sprintf("%s/%s%s", dir, fileType, path)

	cmd := exec.Command("bash", "-c", "ls | grep -v '_index.md'")
	cmd.Dir = findPath
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Println("ls command failed")
	}
}

// ViewMan enables outputs the man page on the terminal
// TODO: Add the option to view the page in less also if it is supported
func ViewMan(fileName string) error {
	dir, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
	if pathErr != nil {
		return pathErr
	}

	filePath := fmt.Sprintf("%s/man/%s.md", dir, fileName)
	if exists, err := isFileExists(filePath); err != nil || !exists {
		log.Println("The file you are looking for does not exist or an error occured")
		return err
	}
	mdtty.Convert(filePath)
	return nil
}

func createFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// Taken from: https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// This file overwrites any existing file in dst. File in dst must exist
func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func delete(filePath string) error {
	if err := os.RemoveAll(filePath); err != nil {
		return err
	}
	return nil
}

// This returns whether the given file or directory exists
func isFileExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// Assume that vim is present.
func openFile(path string) error {
	vimCmd := exec.Command("/usr/bin/vim", path)
	vimCmd.Stdin = os.Stdin
	vimCmd.Stderr = os.Stderr
	vimCmd.Stdout = os.Stdout
	if err := vimCmd.Run(); err != nil {
		log.Println("Error while opening the file in vim. Do you have vim installed?")
		return err
	}
	return nil
}

// Function checks if a certain path is a directory or file
func isDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		log.Printf("Error occured while checking that path at %s is a directory\n", path)
		return false
	}
	return file.Mode().IsDir()
}

func mkdir(path string) error {
	if err := os.Mkdir(path, os.FileMode(0755)); err != nil {
		log.Println("Failed to create new directory")
	}
	return nil
}

// Used to symlink the Hugo directory to the notes directory
func symlinkDir(srcPath string, destPath string) {
	os.Symlink(srcPath, destPath)
}
