# Protocol: ps5-telemetry

A sink, not an emulation. Every request, whatever the method, Host or path, gets:

- `200`, `Content-Type: application/json`, body `{}`.
- The request body is read up to 1 MiB and discarded; nothing is stored.

Exception: `GET /healthz` is answered by ps5-common (`{"status":"ok"}`) for the container health check.

Hosts pointed here by the reference deployment: `telemetry-console.api.playstation.com`,
`smetrics.aem.playstation.com`, `feature.api.playstation.com`. What the console expects back from these hosts is
unverified (function of the names only); `{}` is a guess that costs nothing.

## Hostname census

Each request logs one info line: `msg=telemetry method=... host=... path=... content_length=... content_type=...`.
Counting distinct `host` values over a console session gives the census for PRD Gate 2. The `Host` header is what
the console sent, so point any unknown host's DNS here to learn its traffic.

## Research dump

With `PS5_TELEMETRY_DUMP_DIR` set, each request also writes `<UTC time>-<seq>.txt`: the request line, `Host`, all
headers, a blank line and the first 4 KiB of the body. Dumps can contain device identifiers; keep the directory
private and delete it when done.
