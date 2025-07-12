package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var bootstrap *template.Template

//go:embed www
var wwwFS embed.FS
var subFS fs.FS

func init() {
	var err error
	subFS, err = fs.Sub(wwwFS, "www")
	if err != nil {
		log.Fatalf("Failed to get sub filesystem: %v", err)
	}
	fmt.Println("fs ready")
}

func httpServe() {
	var err error

	bootstrap, err = LoadTemplates(subFS)
	if err != nil {
		panic(err)
	}
	fmt.Println("templates", bootstrap.DefinedTemplates())

	fmt.Println("debug", bootstrap)

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":3000", mux)
}

// 1. See if we have a static file with the same name. Exclude .gohtml and .go files.
// 2. See if we have an index template matching the path. E.g. "/" = "index", "/about" = "about/index"
// 3. If none of the above, 404
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Static files
	p := filepath.Join("./", r.URL.Path)
	info, err := fs.Stat(subFS, p)
	if err == nil && !info.IsDir() {
		// Serve file from subFS
		http.ServeFileFS(w, r, subFS, p)
		return // stop
	}
	if err != nil {
		log.Println("Error", err.Error())
	}

	tctx := map[string]any{}
	tctx["Env"] = map[string]string{
		// Grab from ENV
		"GH_CLIENT_ID": os.Getenv("GH_CLIENT_ID"),
	}

	err = bootstrap.ExecuteTemplate(w, "index", tctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
