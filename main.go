package main

import (
	_ "embed"

	"bufio"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/bmaupin/go-epub"
)

func main() {
	title := flag.String("title", "SourceCodex", "Title of the EPUB")
	author := flag.String("author", "Unknown", "Author of the EPUB")
	output := flag.String("output", "book.epub", "Output file")

	flag.Parse()

	log.Printf("Title: %s", *title)
	log.Printf("Author: %s", *author)
	log.Printf("Output: %s", *output)

	filenames := make([]string, 0)

	// Read the files to include from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		filenames = append(filenames, line)
	}

	// Create epub
	book := mustBuildEpub(*title, *author, filenames)
	err := book.Write(*output)
	if err != nil {
		log.Printf("Error building epub: %s", err)
	}
}

func mustBuildEpub(title string, author string, filenames []string) *epub.Epub {
	book := epub.NewEpub(title)
	book.SetAuthor(author)

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
	Code     template.HTML
}

func mustGetXHTMLContent(filename string) string {
	b := new(strings.Builder)
	rawCode := mustReadFile(filename)
	hlCode := template.HTML(mustHighlightCode(filename, rawCode))
	err := tpl.Execute(b, templateData{
		Filename: filename,
		Code:     hlCode,
	})
	if err != nil {
		log.Fatalf("Error executing template: %s", err)
	}
	return b.String()
}

func mustHighlightCode(filename, code string) string {
	lexer := lexers.Match(filename)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	style := styles.Get("bw")
	formatter := formatters.Get("html")
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		log.Fatalf("Error in lexer: %s", err)
	}
	b := new(strings.Builder)
	err = formatter.Format(b, style, iterator)
	if err != nil {
		log.Fatalf("Error in formatter: %s", err)
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
