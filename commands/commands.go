package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// DeletePage delete files at the given path
func DeletePage(name string, fileType string) error {
	dir, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
	if pathErr != nil {
		return pathErr
	}

	newFile := fmt.Sprintf("%s/%s/%s", dir, fileType, name)
	isPageExists, err := fileExists(newFile)
	if err != nil {
		return err
	}
	if !isPageExists {
		log.Printf("The %s page you want to delete does not exist, please first create it.\n", fileType)
		return nil
	}

	if err := delete(newFile); err != nil {
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

	isPageExists, err := fileExists(fmt.Sprintf("%s/%s/%s", dir, fileType, name))
	if err != nil {
		return err
	}
	if isPageExists {
		log.Printf("The %s page you want to create already exists. Edit it or delete it.\n", fileType)
		return nil
	}

	isTemplateExist, err := fileExists(fmt.Sprintf("%s/templates/%s_template", dir, fileType))
	if err != nil {
		return err
	}

	newFile := fmt.Sprintf("%s/%s/%s", dir, fileType, name)
	if isTemplateExist {
		if err := createFile(newFile); err != nil {
			return err
		}
		copy(fmt.Sprintf("%s/templates/%s_template", dir, fileType), newFile)
	} else {
		if err := createFile(newFile); err != nil {
			return err
		}
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

	newFile := fmt.Sprintf("%s/%s/%s", dir, fileType, name)
	isPageExists, err := fileExists(newFile)
	if err != nil {
		return err
	}
	if !isPageExists {
		log.Printf("The %s page you want to edit does not exist, please first create it.\n", fileType)
		return nil
	}

	if vimErr := openFile(newFile); vimErr != nil {
		delete(newFile)
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

	cmd := exec.Command("ls", fmt.Sprintf("%s/%s", dir, fileType))
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Println("ls command failed")
	}
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
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

// This returns whether the given file or directory exists
func fileExists(filePath string) (bool, error) {
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
		log.Println("Error while opening the file in vim.")
		return err
	}
	return nil
}

// Function checks if a certain path is a directory or file
func isDir(path string) (bool, error) {
	file, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return file.Mode().IsDir(), nil
}
