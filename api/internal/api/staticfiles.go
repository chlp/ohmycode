package api

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"ohmycode_api/pkg/util"
	"os"
	"path"
	"path/filepath"
)

//go:embed client/*
var staticFiles embed.FS

func serveDynamicFiles(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" && r.URL.Path != "/" && r.URL.Path != "/index.html" && !util.IsUuid(r.URL.Path[1:]) {
			file := "./internal/api/client" + r.URL.Path
			if _, err := os.Stat(file); err == nil {
				http.ServeFile(w, r, file)
				return
			}
		}
		http.ServeFile(w, r, "./internal/api/client/index.html")
		return
	})
}

func serveStaticFiles(mux *http.ServeMux) {
	staticFS, _ := fs.Sub(staticFiles, "client")
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
		case "js/app.js":
			fileJsFound = true
		}
		return nil
	})
	if !indexHtmlFound || !styleCssFound || !fileJsFound {
		log.Fatal("important static http file not found")
	}
	indexHtmlData, err := staticFiles.ReadFile("client/index.html")
	if err != nil {
		log.Fatal("index.html not found")
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" && r.URL.Path != "/" && r.URL.Path != "/index.html" && !util.IsUuid(r.URL.Path[1:]) {
			requestedFile := path.Clean(r.URL.Path)
			fileToServe := fmt.Sprintf("client%s", requestedFile)
			if _, err := staticFS.Open(requestedFile[1:]); err == nil {
				w.Header().Set("Content-Type", getMimeType(requestedFile))
				data, _ := staticFiles.ReadFile(fileToServe)
				_, _ = w.Write(data)
				return
			}
		}
		_, _ = w.Write(indexHtmlData)
		return
	})
}

var mimeTypes = map[string]string{
	".html":  "text/html; charset=utf-8",
	".css":   "text/css",
	".js":    "application/javascript",
	".json":  "application/json",
	".xml":   "application/xml",
	".svg":   "image/svg+xml",
	".png":   "image/png",
	".jpg":   "image/jpeg",
	".jpeg":  "image/jpeg",
	".gif":   "image/gif",
	".ico":   "image/x-icon",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".ttf":   "font/ttf",
	".otf":   "font/otf",
	".eot":   "application/vnd.ms-fontobject",
	".mp4":   "video/mp4",
	".webm":  "video/webm",
	".ogg":   "audio/ogg",
	".mp3":   "audio/mpeg",
	".wav":   "audio/wav",
	".zip":   "application/zip",
	".pdf":   "application/pdf",
	".txt":   "text/plain; charset=utf-8",
}

func getMimeType(filename string) string {
	ext := filepath.Ext(filename)
	if mime, found := mimeTypes[ext]; found {
		return mime
	}
	return "application/octet-stream"
}
