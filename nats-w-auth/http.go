package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:embed all:www
var wwwFS embed.FS
var subFS fs.FS

func init() {
	var err error
	subFS, err = fs.Sub(wwwFS, "www")
	if err != nil {
		log.Fatalf("Failed to get sub filesystem: %v", err)
	}
	templates, err = loadTemplates(subFS)
	if err != nil {
		panic(err)
	}

}

func httpServe() {

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/success", githubSuccessHandler)
	mux.HandleFunc("/", fallbackHandler)

	http.ListenAndServe(":3000", sessionManager.LoadAndSave(mux))
}

// 1. See if we have a static file with the same name. Exclude .gohtml and .go files.
// 2. See if we have an index template matching the path. E.g. "/" = "index", "/about" = "about/index"
// 3. If none of the above, 404
func fallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Static files
	p := filepath.Join("./", r.URL.Path)
	info, err := fs.Stat(subFS, p)

	// Serve static files if available
	// - Exclude gohtml files
	if err == nil && !info.IsDir() && !strings.HasSuffix(p, ".gohtml") {
		http.ServeFileFS(w, r, subFS, p)
		return
	}

	// Template data
	// TODO: Can we make this dynamic?
	pieIs := sessionManager.GetString(r.Context(), "pie")
	log.Println("debug pie", pieIs)
	if pieIs == "" {
		pieIs = "undefined"
	}

	tctx := map[string]any{}
	tctx["Env"] = map[string]string{
		// Grab from ENV
		"GH_CLIENT_ID": os.Getenv("GH_CLIENT_ID"),
	}
	tctx["Session"] = map[string]string{
		"pie": pieIs,
	}

	templateName := pathToTemplateName(r.URL.Path)
	if templateName == "" {
		http.NotFound(w, r)
		return
	}

	//fmt.Println("templates", templates.DefinedTemplates())

	err = templates.ExecuteTemplate(w, templateName, tctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func pathToTemplateName(p string) string {
	//fmt.Println("path", p)
	p, found := strings.CutPrefix(p, "/")
	if !found || p == "" {
		return "index"
	}

	searchOrder := []string{
		p,
		path.Join(p, "index"),
	}

	for _, tname := range searchOrder {
		tp := templates.Lookup(tname)
		if tp != nil {
			return tp.Name()
		}
	}

	return "" // 404
}
