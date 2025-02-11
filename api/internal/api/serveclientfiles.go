package api

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"ohmycode_api/pkg/util"
	"path"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

func serveStaticFiles(mux *http.ServeMux) {
	staticFS, _ := fs.Sub(staticFiles, "static")
	indexHtmlFound := false
	styleCssFound := false
	fileJsFound := false
	_ = fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		switch path {
		case "index.html":
			indexHtmlFound = true
		case "style.css":
			styleCssFound = true
		case "js/file.js":
			fileJsFound = true
		}
		return nil
	})
	if !indexHtmlFound || !styleCssFound || !fileJsFound {
		log.Fatal("important static http file not found")
	}
	indexHtmlData, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		log.Fatal("index.html not found")
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" && r.URL.Path != "/" && r.URL.Path != "/index.html" && !util.IsUuid(r.URL.Path[1:]) {
			requestedFile := path.Clean(r.URL.Path)
			fileToServe := fmt.Sprintf("static%s", requestedFile)
			if _, err := staticFS.Open(requestedFile[1:]); err == nil {
				if strings.HasSuffix(requestedFile, ".css") {
					w.Header().Set("Content-Type", "text/css")
				}
				data, _ := staticFiles.ReadFile(fileToServe)
				_, _ = w.Write(data)
				return
			}
		}
		_, _ = w.Write(indexHtmlData)
		return
	})
}
