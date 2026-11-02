# PlayStation 5 Telemetry Sink

OpenPak's PlayStation 5 Telemetry Sink.

| | |
| --- | --- |
| Network | PlayStation 5 |
| Role | Answers the console's telemetry, analytics and feature-flag hosts with `200 {}` and logs a hostname census |
| Used by | PlayStation 5 system software |
| Image | `ghcr.io/openpak/ps5-telemetry` |
| Licence | AGPL-3.0-only |

## What it serves

| Endpoint | What |
| --- | --- |
| any method, any path (`telemetry-console.api.playstation.com`, `smetrics.aem.playstation.com`, `feature.api.playstation.com`, any host pointed here) | `200`, `{}`; body discarded |
| `GET /healthz` | health check |

## Run

```sh
GOWORK=off go run ./cmd/ps5-telemetry      # listens on :8080
curl -i -X POST http://127.0.0.1:8080/any/path -d x=1
```

## Docs

- [Protocol](docs/protocol.md)
- [Deploy](docs/deploy.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
