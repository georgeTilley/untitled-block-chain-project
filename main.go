package main

import (
	"embed"
	"encoding/gob"
	"io/ioutil"
	"log"
	"text/template"

	"github.com/gorilla/sessions"
)

var (
	store     *sessions.CookieStore
	templates *template.Template
)

func init() {
	store = sessions.NewCookieStore([]byte("your-secret-key"))
	template, err := template.ParseFS(embedFiles, "static/*.html")
	templates = template
	if err != nil {
		log.Fatal(err)
	}
	gob.Register(Blockchain{})
	gob.Register([]*Block{})

}

func main() {
	// Write embedded files to disk
	writeFile("static/index.html", "static/index.html")
	writeFile("static/result.html", "static/result.html")

	newblockchain := NewBlockChain()

	startWebServer(*newblockchain)
}

//go:embed static/index.html static/result.html*
var embedFiles embed.FS

func writeFile(filename, embedPath string) {
	content, err := embedFiles.ReadFile(embedPath)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		panic(err)
	}
}
