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
	"strings"
)

//go:embed client/*
var staticFiles embed.FS

func setCacheHeadersForJS(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".js") {
		return
	}
	// Manual cache control:
	// - If URL includes ?v=... treat as versioned asset → cache "forever".
	// - Otherwise allow caching but force revalidation (no cache-bust needed, and no module duplication).
	if _, ok := r.URL.Query()["v"]; ok {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
	}
}

func setNoCacheForIndex(w http.ResponseWriter) {
	// index.html should be revalidated so changes to the entrypoint URL (?v=...) are picked up reliably.
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func serveDynamicFiles(mux *http.ServeMux) {
	const diskRoot = "./internal/api/client"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" && r.URL.Path != "/" && r.URL.Path != "/index.html" && !util.IsUuid(r.URL.Path[1:]) {
			cleaned := path.Clean("/" + strings.TrimPrefix(r.URL.Path, "/"))
			// prevent traversal like /../../etc/passwd
			if strings.HasPrefix(cleaned, "/..") {
				http.NotFound(w, r)
				return
			}
			rel := strings.TrimPrefix(cleaned, "/")
			file := filepath.Join(diskRoot, filepath.FromSlash(rel))
			if st, err := os.Stat(file); err == nil && !st.IsDir() {
				setCacheHeadersForJS(w, r)
				http.ServeFile(w, r, file)
				return
			}
		}
		setNoCacheForIndex(w)
		http.ServeFile(w, r, filepath.Join(diskRoot, "index.html"))
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
			if f, err := staticFS.Open(requestedFile[1:]); err == nil {
				_ = f.Close()
				setCacheHeadersForJS(w, r)
				w.Header().Set("Content-Type", getMimeType(requestedFile))
				data, err := staticFiles.ReadFile(fileToServe)
				if err != nil {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				_, _ = w.Write(data)
				return
			}
		}
		setNoCacheForIndex(w)
		_, _ = w.Write(indexHtmlData)
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
