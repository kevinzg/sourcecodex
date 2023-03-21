package main

import (
	_ "embed"

	"bufio"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/bmaupin/go-epub"
)

func main() {
	filenames := make([]string, 0)

	// Read the files to include from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		filenames = append(filenames, line)
	}

	// Create epub
	book := mustBuildEpub("Source code", "Author", filenames)
	err := book.Write("book.epub")
	if err != nil {
		log.Printf("Error building epub: %s", err)
	}
}

func mustBuildEpub(title string, author string, filenames []string) *epub.Epub {
	book := epub.NewEpub("Source code")

	for i, filename := range filenames {
		content := mustGetXHTMLContent(filename)
		// The internal filename cannot contain slashes!
		name := fmt.Sprintf("%06d_%s.xhtml", i, strings.ReplaceAll(filename, "/", "__"))
		_, err := book.AddSection(content, filename, name, "")
		if err != nil {
			log.Fatalf("Error adding section for file %s: %s", filename, err)
		}
		log.Printf("Added section for file %s", filename)
	}

	return book
}

//go:embed page.xhtml.tpl
var rawTemplate string
var tpl = template.Must(template.New("page").Parse(rawTemplate))

type templateData struct {
	Filename string
	Code     string
}

func mustGetXHTMLContent(filename string) string {
	b := new(strings.Builder)
	code := mustReadFile(filename)
	err := tpl.Execute(b, templateData{
		Filename: filename,
		Code:     code,
	})
	if err != nil {
		log.Fatalf("Error executing template: %s", err)
	}
	return b.String()
}

func mustReadFile(name string) string {
	b, err := os.ReadFile(name)
	if err != nil {
		log.Fatalf("Cannot read file %s: %s", name, err)
	}
	return string(b)
}
