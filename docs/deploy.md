# Deploy: ps5-telemetry

Architecture: OpenPak's DNS sends the telemetry Host names to OpenPak; the reverse proxy forwards them to this
service on `PS5_ADDR` (default `:8080`). TLS terminates in front. The Host header is only logged.

## Run

```sh
cp example.env .env
set -a; source .env; set +a
GOWORK=off go run ./cmd/ps5-telemetry
```

```sh
TMPDIR=<scratch> GOWORK=off go test ./...
GOWORK=off go vet ./... && gofmt -l .
docker build -t ps5-telemetry .   # build context must contain vendor/
```

The module consumes `github.com/openpak/ps5-common` through `replace ... => ../ps5-common` with a committed
`vendor/` (see ps5-common's deploy doc). After changing ps5-common run `GOWORK=off go mod vendor` here.

## Configuration

| Variable | Required | Default | Meaning |
|---|---|---|---|
| `PS5_ADDR` | no | `:8080` | HTTP listen address |
| `PS5_LOG_LEVEL` | no | `info` | `debug`, `info`, `warn`, `error`; the census lines are info |
| `PS5_TELEMETRY_DUMP_DIR` | no | off | write headers + first 4 KiB of body per request to this directory (research) |

## Release

Tag `v*.*.*` publishes `ghcr.io/openpak/ps5-telemetry`; deployed hosts auto-update from it.
