// Package main generates boilerplate code for metadata document types
package main

import (
	"io/ioutil"
	"log"
	"strings"
	"text/template"
	"path/filepath"
	"os"
	"fmt"
	"bitbucket.org/yasnikoff/shuffler"
)

type TYPE struct {
	shuffler.Header
}

type TypeData struct {
	Name     string // name of the package
	TypeName string // name of the main type
	RcvName  string // name of the receiver value
}

type dummyDoc struct {
	shuffler.Doc
}

func New() shuffler.Doc {
	// dummy. just to make compiler happy. Implementations are in the real packages.
	return new(dummyDoc)
}

func main() {

	root := "../../"
	typeData := []TypeData{
		{"project", "Project", "pr"},
		{"repo", "Repo", "rp"},
		{"collection", "Collection", "c"},
	}
	paths := make([]string, len(typeData))
	for i, td := range typeData {
		paths[i] = filepath.Join(root, td.Name, fmt.Sprintf("%s_boilerplate.go", td.Name))
	}

	data, err := ioutil.ReadFile("doc_boilerplate.go")
	if err != nil {
		log.Fatal(err)
	}
	tmplStr := string(data)
	tmplStr = strings.Replace(tmplStr, "TYPE", "{{.TypeName}}", -1)
	tmplStr = strings.Replace(tmplStr, "RECEIVER", "{{.RcvName}}", -1)
	tmplStr = strings.Replace(tmplStr, "package main", "package {{.Name}}", -1)
	tmplStr = strings.Replace(tmplStr, "INIT()", "init()", -1)
	tmpl, err := template.New("doctype_boilerplate").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}

	for i, td := range typeData {
		f, err := os.Create(paths[i])
		if err != nil {
			log.Fatal(err)
		}
		if err := tmpl.Execute(f, td); err != nil {
			log.Print(err)
		}
		f.Close()
	}
	log.Println("Done.")
}
