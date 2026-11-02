// Package sink is the catch-all handler: it discards request bodies, answers 200 {} and
// logs enough to build the console hostname census.
package sink

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

const (
	// BodyCap is how much of a request body is read before the rest is dropped.
	BodyCap = 1 << 20
	// DumpBody is how much of the body a dump file keeps.
	DumpBody = 4 << 10
)

// Handler returns the sink. dumpDir "" disables dumping.
func Handler(log *slog.Logger, dumpDir string) http.Handler {
	var seq atomic.Uint64
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(io.LimitReader(r.Body, BodyCap))
		_, _ = io.Copy(io.Discard, io.LimitReader(r.Body, BodyCap))
		log.Info("telemetry", "method", r.Method, "host", r.Host, "path", r.URL.Path,
			"content_length", r.ContentLength, "content_type", r.Header.Get("Content-Type"))
		if dumpDir != "" {
			if err := dump(dumpDir, seq.Add(1), r, body); err != nil {
				log.Warn("dump", "err", err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "{}")
	})
}

func dump(dir string, n uint64, r *http.Request, body []byte) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	if len(body) > DumpBody {
		body = body[:DumpBody]
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s HTTP/%d.%d\nHost: %s\n", r.Method, r.URL.RequestURI(), r.ProtoMajor, r.ProtoMinor, r.Host)
	for k, vs := range r.Header {
		for _, v := range vs {
			fmt.Fprintf(&b, "%s: %s\n", k, v)
		}
	}
	b.WriteString("\n")
	b.Write(body)
	name := fmt.Sprintf("%s-%06d.txt", time.Now().UTC().Format("20060102T150405.000Z"), n)
	return os.WriteFile(filepath.Join(dir, name), []byte(b.String()), 0o640)
}
