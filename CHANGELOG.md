# Changelog: ps5-telemetry

## Unreleased

- 2026-11-02: first cut. Catch-all sink: any method on any path answers `200 {}`, bodies are read up to 1 MiB and
  discarded, one `telemetry` info log line per request (method, host, path, content length, content type) for the
  hostname census. `PS5_TELEMETRY_DUMP_DIR` (off by default) writes headers plus the first 4 KiB of the body per
  request. `GET /healthz` from ps5-common. Tests with httptest. Not deployed, no tag, not tried against a console.
