package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// GenPage generates a new manpage based on the one in templates, or creates a new one if the template page is missing
func GenPage(name string) error {
	dir, pathErr := filepath.Abs(filepath.Dir(os.Args[0]))
	if pathErr != nil {
		return pathErr
	}

	isManPageExists, err := fileExists(fmt.Sprintf("%s/../pages/%s", dir, name))
	if err != nil {
		return err
	}
	if isManPageExists {
		log.Println("The man page you want to create already exists. Edit it or delete it.")
		return nil
	}

	isTemplateExist, err := fileExists(fmt.Sprintf("%s../templates/man_template", dir))
	if err != nil {
		return err
	}

	newFile := fmt.Sprintf("%s/../pages/%s", dir, name)
	if isTemplateExist {
		copy(fmt.Sprintf("%s/../templates/man_template", dir), newFile)
	} else {
		file, err := os.Create(newFile)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	if vimErr := openFile(newFile); vimErr != nil {
		return err
	}

	return nil
}

// Taken from: https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// This file overwrites any existing file in dst.
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

func fileExists(filePath string) (bool, error) {
	if _, err := os.Stat(fmt.Sprintf("%s", filePath)); !os.IsNotExist(err) {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func openFile(path string) error {
	vimCmd := exec.Command("vim", path)
	if err := vimCmd.Run(); err != nil {
		log.Println("Error while opening the file in vim.")
		return err
	}
	return nil
}
