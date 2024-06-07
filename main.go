package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	inputPath := "./input"
	outputPath := "./output"

	files, err := ioutil.ReadDir(inputPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".epub" {
			processEpub(filepath.Join(inputPath, file.Name()), filepath.Join(outputPath, file.Name()))
		}
	}
}

func processEpub(inputFile, outputFile string) {
	r, err := zip.OpenReader(inputFile)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", inputFile, err)
		return
	}
	defer r.Close()

	newZipFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", outputFile, err)
		return
	}
	defer newZipFile.Close()

	w := zip.NewWriter(newZipFile)
	defer w.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			fmt.Println("Error opening zip content:", err)
			continue
		}

		body, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			fmt.Println("Error reading zip content:", err)
			continue
		}

		if strings.HasSuffix(f.Name, ".xhtml") {
			body = bytes.ReplaceAll(body, []byte("\u00A0"), []byte(" "))
		}

		if strings.HasSuffix(f.Name, "style.css") {
			body = updateCSS(string(body))
		}

		fw, err := w.Create(f.Name)
		if err != nil {
			fmt.Println("Error creating zip content:", err)
			continue
		}

		_, err = fw.Write(body)
		if err != nil {
			fmt.Println("Error writing zip content:", err)
			continue
		}
	}
}

func updateCSS(cssContent string) []byte {
	re := regexp.MustCompile(`(?s)p\s*\{([^}]*)\}`)
	matches := re.FindStringSubmatch(cssContent)

	if len(matches) == 1 {
		fmt.Println("Warning: empty 'p { }' block in 'style.css'")
	} else if len(matches) > 1 {
		existingStyles := strings.TrimSpace(matches[1])
		if !strings.Contains(existingStyles, "white-space:") {
			newStyles := existingStyles + "\n    white-space: pre-wrap;"
			cssContent = re.ReplaceAllString(cssContent, "p {\n"+newStyles+"\n}")
		}
	} else {
		fmt.Println("Warning: 'p { }' block not found in 'style.css'")
	}

	return []byte(cssContent)
}
