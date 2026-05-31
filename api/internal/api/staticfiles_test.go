package api

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func makeTestFS(files map[string]string) fs.FS {
	m := make(fstest.MapFS)
	for name, content := range files {
		m[name] = &fstest.MapFile{Data: []byte(content)}
	}
	return m
}

// --- computeBuildHash ---

func TestComputeBuildHash_Consistent(t *testing.T) {
	testFS := makeTestFS(map[string]string{
		"js/main.js": `import "./status.js";`,
		"style.css":  `body { margin: 0; }`,
	})
	if computeBuildHash(testFS) != computeBuildHash(testFS) {
		t.Error("hash should be deterministic for same content")
	}
}

func TestComputeBuildHash_LengthIsEight(t *testing.T) {
	testFS := makeTestFS(map[string]string{"a.js": "x"})
	if h := computeBuildHash(testFS); len(h) != 8 {
		t.Errorf("expected 8-char hash, got len=%d %q", len(h), h)
	}
}

func TestComputeBuildHash_DifferentContentDifferentHash(t *testing.T) {
	fs1 := makeTestFS(map[string]string{"main.js": "version A"})
	fs2 := makeTestFS(map[string]string{"main.js": "version B"})
	if computeBuildHash(fs1) == computeBuildHash(fs2) {
		t.Error("different content must produce different hashes")
	}
}

func TestComputeBuildHash_EmptyFS(t *testing.T) {
	if h := computeBuildHash(makeTestFS(nil)); len(h) != 8 {
		t.Errorf("empty FS: expected 8-char hash, got %q", h)
	}
}

// --- jsImportRe ---

func TestJsImportRe_PatchesRelativeImports(t *testing.T) {
	cases := []struct {
		input    string
		wantSub  string
	}{
		{`import "./status.js";`, `"./status.js?v=h"`},
		{`import {x} from "./connect.js";`, `"./connect.js?v=h"`},
		{`import {a, b} from "./app.js";`, `"./app.js?v=h"`},
	}
	for _, c := range cases {
		got := string(jsImportRe.ReplaceAll([]byte(c.input), []byte(`"$1?v=h"`)))
		if !strings.Contains(got, c.wantSub) {
			t.Errorf("input %q: want %q in output, got %q", c.input, c.wantSub, got)
		}
	}
}

func TestJsImportRe_SkipsAlreadyVersioned(t *testing.T) {
	// ? excluded from [^"?]+ so already-versioned imports are not re-matched
	got := string(jsImportRe.ReplaceAll([]byte(`import "./status.js?v=old";`), []byte(`"$1?v=new"`)))
	if strings.Contains(got, "?v=new") {
		t.Errorf("should not re-patch already-versioned import, got: %s", got)
	}
}

func TestJsImportRe_SkipsNonJsFiles(t *testing.T) {
	got := string(jsImportRe.ReplaceAll([]byte(`import "./style.css";`), []byte(`"$1?v=h"`)))
	if strings.Contains(got, "?v=h") {
		t.Errorf("should not patch non-.js import, got: %s", got)
	}
}

func TestJsImportRe_SkipsAbsoluteURLs(t *testing.T) {
	got := string(jsImportRe.ReplaceAll([]byte(`import "https://cdn.example.com/lib.js";`), []byte(`"$1?v=h"`)))
	if strings.Contains(got, "?v=h") {
		t.Errorf("should not patch absolute URL import, got: %s", got)
	}
}

// --- versionRe ---

func TestVersionRe_ReplacesAllOccurrences(t *testing.T) {
	input := `<link href="style.css?v=16"><script src="js/main.js?v=18"></script>`
	got := string(versionRe.ReplaceAll([]byte(input), []byte("?v=xyz")))
	if strings.Contains(got, "?v=16") || strings.Contains(got, "?v=18") {
		t.Errorf("old versions should be replaced, got: %s", got)
	}
	if !strings.Contains(got, "?v=xyz") {
		t.Errorf("replacement not found in: %s", got)
	}
}

// --- setCacheHeadersForJS ---

func TestSetCacheHeadersForJS_VersionedIsImmutable(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/js/main.js?v=abc12345", nil)
	setCacheHeadersForJS(w, r)
	cc := w.Header().Get("Cache-Control")
	if !strings.Contains(cc, "immutable") || !strings.Contains(cc, "max-age=31536000") {
		t.Errorf("versioned JS: want immutable cache, got %q", cc)
	}
}

func TestSetCacheHeadersForJS_UnversionedIsMustRevalidate(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/js/main.js", nil)
	setCacheHeadersForJS(w, r)
	cc := w.Header().Get("Cache-Control")
	if !strings.Contains(cc, "must-revalidate") {
		t.Errorf("unversioned JS: want must-revalidate, got %q", cc)
	}
}

func TestSetCacheHeadersForJS_NonJsNoHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/style.css?v=abc", nil)
	setCacheHeadersForJS(w, r)
	if got := w.Header().Get("Cache-Control"); got != "" {
		t.Errorf("CSS should not get Cache-Control from this function, got %q", got)
	}
}

// --- serveStaticFiles HTTP integration ---

func serveStaticTestMux(t *testing.T) http.Handler {
	t.Helper()
	mux := http.NewServeMux()
	serveStaticFiles(mux)
	return mux
}

func TestServeStaticFiles_IndexHtmlHasNoCache(t *testing.T) {
	h := serveStaticTestMux(t)
	for _, path := range []string{"/", "/index.html"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != http.StatusOK {
			t.Errorf("GET %s: want 200, got %d", path, w.Code)
		}
		cc := w.Header().Get("Cache-Control")
		if cc != "no-cache" {
			t.Errorf("GET %s: want Cache-Control=no-cache, got %q", path, cc)
		}
	}
}

func TestServeStaticFiles_IndexHtmlVersionIsHash(t *testing.T) {
	h := serveStaticTestMux(t)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	body := w.Body.String()
	// Hardcoded version numbers like ?v=18 should have been replaced by the computed hash
	if strings.Contains(body, "?v=18") || strings.Contains(body, "?v=16") {
		t.Error("index.html should not contain hardcoded version numbers after hash injection")
	}
	if !strings.Contains(body, "?v=") {
		t.Error("index.html should contain at least one versioned asset reference")
	}
}

func TestServeStaticFiles_IndexHtmlVersionConsistentAcrossRequests(t *testing.T) {
	h := serveStaticTestMux(t)
	getBody := func() string {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
		return w.Body.String()
	}
	b1, b2 := getBody(), getBody()
	if b1 != b2 {
		t.Error("index.html content should be identical across requests")
	}
}

func TestServeStaticFiles_MainJsImportsAreVersioned(t *testing.T) {
	h := serveStaticTestMux(t)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/js/main.js", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET /js/main.js: want 200, got %d", w.Code)
	}
	body := w.Body.String()
	for _, line := range strings.Split(body, "\n") {
		if strings.Contains(line, `"./`) && strings.Contains(line, ".js") && !strings.Contains(line, "?v=") {
			t.Errorf("unversioned import found in main.js: %q", strings.TrimSpace(line))
		}
	}
}

func TestServeStaticFiles_AllJsModulesHaveVersionedImports(t *testing.T) {
	h := serveStaticTestMux(t)
	for _, mod := range []string{"/js/connect.js", "/js/editor.js", "/js/app.js", "/js/run.js"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", mod, nil))
		if w.Code != http.StatusOK {
			t.Errorf("GET %s: want 200, got %d", mod, w.Code)
			continue
		}
		for _, line := range strings.Split(w.Body.String(), "\n") {
			if strings.Contains(line, `from "./`) && !strings.Contains(line, "?v=") {
				t.Errorf("%s: unversioned import: %q", mod, strings.TrimSpace(line))
			}
		}
	}
}

func TestServeStaticFiles_AllJsModuleHashesMatch(t *testing.T) {
	// All versioned imports should use the same hash value.
	h := serveStaticTestMux(t)

	extractVersion := func(body string) string {
		const marker = "?v="
		idx := strings.Index(body, marker)
		if idx < 0 {
			return ""
		}
		rest := body[idx+len(marker):]
		end := strings.IndexAny(rest, `"'& `)
		if end < 0 {
			return rest
		}
		return rest[:end]
	}

	indexW := httptest.NewRecorder()
	h.ServeHTTP(indexW, httptest.NewRequest("GET", "/", nil))
	indexVersion := extractVersion(indexW.Body.String())
	if indexVersion == "" {
		t.Fatal("no versioned asset found in index.html")
	}

	mainW := httptest.NewRecorder()
	h.ServeHTTP(mainW, httptest.NewRequest("GET", "/js/main.js", nil))
	mainVersion := extractVersion(mainW.Body.String())

	if mainVersion != "" && mainVersion != indexVersion {
		t.Errorf("index.html uses hash %q but main.js imports use hash %q", indexVersion, mainVersion)
	}
}

func TestServeStaticFiles_VersionedJsGetsImmutableCache(t *testing.T) {
	h := serveStaticTestMux(t)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/js/main.js?v=somehash", nil))
	cc := w.Header().Get("Cache-Control")
	if !strings.Contains(cc, "immutable") {
		t.Errorf("versioned JS should get immutable cache, got %q", cc)
	}
}

func TestServeStaticFiles_SpaRoutingServesIndexHtml(t *testing.T) {
	// A 22-char base62 path (valid file ID) should serve index.html for SPA routing
	h := serveStaticTestMux(t)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/AAAAAAAAAAAAAAAAAAAAAA", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("SPA route: want 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "<html") {
		t.Error("SPA route should serve index.html HTML content")
	}
}

func TestServeStaticFiles_ContentTypeIsJavascript(t *testing.T) {
	h := serveStaticTestMux(t)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/js/main.js", nil))
	ct := w.Header().Get("Content-Type")
	if !strings.Contains(ct, "javascript") {
		t.Errorf("expected javascript Content-Type for JS file, got %q", ct)
	}
}

func TestStaticFilesEmbedStructure(t *testing.T) {
	staticFS, err := fs.Sub(staticFiles, "client")
	if err != nil {
		t.Fatalf("fs.Sub: %v", err)
	}
	allowedDirs := map[string]bool{"js": true, "codemirror": true, "md": true}
	requiredFiles := []string{"index.html", "style.css", "js/main.js"}

	_ = fs.WalkDir(staticFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || p == "." || !d.IsDir() {
			return err
		}
		if !strings.Contains(p, "/") && !allowedDirs[p] {
			t.Errorf("unexpected top-level directory in embedded FS: %q", p)
		}
		return nil
	})

	for _, f := range requiredFiles {
		if _, err := fs.Stat(staticFS, f); err != nil {
			t.Errorf("required file missing from embed: %q", f)
		}
	}
}
