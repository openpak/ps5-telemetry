package sink

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func quiet() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func TestCatchAll(t *testing.T) {
	h := Handler(quiet(), "")
	for _, m := range []string{"GET", "POST", "PUT", "DELETE"} {
		req := httptest.NewRequest(m, "http://telemetry-console.api.playstation.com/v1/x/y", strings.NewReader("data"))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "{}" {
			t.Fatalf("%s: %d %q", m, rec.Code, rec.Body.String())
		}
	}
}

func TestDump(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "d")
	h := Handler(quiet(), dir)
	req := httptest.NewRequest("POST", "http://smetrics.aem.playstation.com/b/ss", strings.NewReader(strings.Repeat("a", 10000)))
	req.Header.Set("X-Test", "yes")
	h.ServeHTTP(httptest.NewRecorder(), req)
	files, _ := os.ReadDir(dir)
	if len(files) != 1 {
		t.Fatalf("want 1 dump, got %d", len(files))
	}
	b, _ := os.ReadFile(filepath.Join(dir, files[0].Name()))
	s := string(b)
	if !strings.Contains(s, "X-Test: yes") || !strings.Contains(s, "Host: smetrics.aem.playstation.com") {
		t.Fatalf("headers missing: %q", s[:200])
	}
	if i := strings.Index(s, "\n\n"); len(s)-i-2 != DumpBody {
		t.Fatalf("body len %d", len(s)-i-2)
	}
}

func TestNoDumpByDefault(t *testing.T) {
	dir := t.TempDir()
	Handler(quiet(), "").ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if f, _ := os.ReadDir(dir); len(f) != 0 {
		t.Fatal("unexpected dump")
	}
}
