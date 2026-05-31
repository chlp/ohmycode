package api

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"ohmycode_api/pkg/util"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	jsImportRe = regexp.MustCompile(`"(\./[^"?]+\.js)"`)
	versionRe  = regexp.MustCompile(`\?v=\w+`)
)

func computeBuildHash(staticFS fs.FS) string {
	h := sha256.New()
	_ = fs.WalkDir(staticFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := fs.ReadFile(staticFS, p)
		if err != nil {
			return nil
		}
		h.Write(data)
		return nil
	})
	return fmt.Sprintf("%x", h.Sum(nil))[:8]
}

//go:embed client/index.html client/style.css client/favicon* client/github-mark.svg client/js client/codemirror client/md
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
		if r.URL.Path != "" && r.URL.Path != "/" && r.URL.Path != "/index.html" && !util.IsValidId(r.URL.Path[1:]) {
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
	mainJsFound := false
	_ = fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		switch path {
		case "index.html":
			indexHtmlFound = true
		case "style.css":
			styleCssFound = true
		case "js/main.js":
			mainJsFound = true
		}
		return nil
	})
	if !indexHtmlFound || !styleCssFound || !mainJsFound {
		log.Fatal("important static http file not found")
	}

	hash := computeBuildHash(staticFS)
	log.Printf("Static build hash: %s", hash)

	indexHtmlData, err := staticFiles.ReadFile("client/index.html")
	if err != nil {
		log.Fatal("index.html not found")
	}
	indexHtmlData = versionRe.ReplaceAll(indexHtmlData, []byte("?v="+hash))

	// Patch relative imports in every JS file so all modules share the same versioned URLs.
	// Without this, different modules importing the same file with and without ?v= would
	// create duplicate module instances, breaking shared state.
	patchedJS := make(map[string][]byte)
	_ = fs.WalkDir(staticFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".js") {
			return err
		}
		data, readErr := staticFiles.ReadFile("client/" + p)
		if readErr != nil {
			return nil
		}
		patchedJS["/"+p] = jsImportRe.ReplaceAll(data, []byte(`"$1?v=`+hash+`"`))
		return nil
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" && r.URL.Path != "/" && r.URL.Path != "/index.html" && !util.IsValidId(r.URL.Path[1:]) {
			requestedFile := path.Clean(r.URL.Path)

			if patched, ok := patchedJS[requestedFile]; ok {
				setCacheHeadersForJS(w, r)
				w.Header().Set("Content-Type", "application/javascript")
				_, _ = w.Write(patched)
				return
			}

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
