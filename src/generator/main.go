package main

import (
	"bytes"
	"flag"
	"fmt"
	"text/template"
)

var (
	templateFile = flag.String("c", "", "path to repo template file")
	repoName     = flag.String("n", "", "the repo name")
	pkgName      = flag.String("p", "", "the name of package")
	modelName    = flag.String("m", "", "the model name")
	tmpFileName  = flag.String("t", "repo.go.tmpl", "template name")
)

func main() {
	flag.Parse()
	vars := make(map[string]interface{})
	vars["Repo"] = *repoName
	vars["Name"] = *pkgName
	vars["Model"] = *modelName

	result := processFile(*templateFile, vars)

	fmt.Println(result)
}

func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	return tmplBytes.String()
}

func processFile(fileName string, vars interface{}) string {
	tmpl, err := template.New(*tmpFileName).Delims("[[", "]]").ParseFiles(fileName)

	if err != nil {
		panic(err)
	}
	return process(tmpl, vars)
}
