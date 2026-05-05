# tesserato

A Go CLI on top of Docker and Docker Compose, built to centralize the day-to-day operations of multiple containerized projects running on a single host (typically a VPS).

If you manage several Docker-based projects with distinct `docker-compose.yml` files spread across directories — and you find yourself typing long `docker compose -f ... up -d --build` incantations every day — `tesserato` collapses that workflow into short, project-aware commands driven by a single YAML config.

## Features

| Command | What it does |
| --- | --- |
| `tesserato status` | Lists running containers grouped by the projects defined in your config. Falls back to a flat list when no config is present. |
| `tesserato up <project> [--build]` | Runs `docker compose up -d` (optionally `--build`) inside the project's directory. |
| `tesserato down <project> [--volumes]` | Stops and removes the project's containers via `docker compose down`. |
| `tesserato logs <project> <service> [-f] [--tail N]` | Tails the logs of a single service, resolved through the config's logical-name → container-name map. |
| `tesserato health` | For every declared service across all projects, reports `Up`, `Exited (code) ...` or `missing`. |
| `tesserato clean [--all] [-y]` | Prunes stopped containers and dangling images. With `--all`, runs the more aggressive `docker system prune -a --volumes`. Asks for confirmation unless `-y`. |
| `tesserato deploy <project> [--no-pull] [--branch X]` | Production deploy flow: optional `git checkout <branch>`, then `git pull --ff-only`, then `docker compose up -d --build`. Fails fast on git errors. |
| `tesserato --version` | Prints the binary version (default `dev`; can be injected at build time). |

A global `--config <path>` flag overrides config discovery for any command.

## Installation

### Prerequisites

- Go 1.26+ (`go version` to verify)
- Docker and Docker Compose (the CLI shells out to them)
- Git (only needed for `tesserato deploy`)

### Build and install

```powershell
go install github.com/RafhaelH/cli_go/cmd/tesserato@latest
```

This drops the binary in `$(go env GOPATH)\bin\tesserato.exe` (Windows) or `$(go env GOPATH)/bin/tesserato` (macOS/Linux), which the Go installer adds to your `PATH` automatically.

To build from a clone:

```powershell
git clone https://github.com/RafhaelH/cli_go.git
cd cli_go
go install ./cmd/tesserato
```

### Build with embedded version

The version reported by `tesserato --version` is set via `-ldflags`:

```powershell
go install -ldflags "-X 'github.com/RafhaelH/cli_go/internal/cli.Version=v0.1.0'" ./cmd/tesserato
```

Pair it with `git describe` in a release pipeline to tag binaries automatically.

## Configuration

`tesserato` is driven by a YAML file. By default it looks for, in order:

1. The path passed via `--config <path>`.
2. `./tesserato.yaml` or `./tesserato.yml` in the current directory.
3. `~/.config/tesserato/config.yaml`.

If none are found and no command needs the config (e.g. `clean`, or a no-arg `status` falling back to flat mode), the CLI runs anyway. Otherwise it prints an actionable error explaining where it looked.

### Example config

```yaml
projects:
  portal_astech:
    path: /home/user/projects/portal_astech
    compose_file: docker-compose.prod.yml
    services:
      backend: portal_astech_backend
      postgres: portal_astech_postgres
      redis: portal_astech_redis

  knowledgeai:
    path: /home/user/projects/knowledgeai
    compose_file: docker-compose.yml
    services:
      api: knowledgeai_api
      frontend: knowledgeai_frontend
      db: knowledgeai_db
```

Schema:

| Field | Required | Description |
| --- | --- | --- |
| `projects` | yes | Top-level map keyed by the logical project name you'll pass on the CLI. |
| `projects.<name>.path` | yes | Absolute path to the project directory (where the compose file lives). |
| `projects.<name>.compose_file` | no | Compose file name. Defaults to whatever `docker compose` would pick if omitted. |
| `projects.<name>.services` | yes | Map of logical service name → real container name. The container name is what `docker ps` prints under `NAMES`. |

The `services` mapping is what powers the project-grouped views: `tesserato status` and `tesserato health` correlate `docker ps` output against this map to know which container belongs to which project.

## Usage

Once installed and configured:

```powershell
# What's running, by project?
tesserato status

# Start a project
tesserato up portal_astech --build

# Tail one service
tesserato logs portal_astech backend -f --tail 100

# Stop a project
tesserato down portal_astech

# Production deploy
tesserato deploy portal_astech

# Health snapshot across all projects
tesserato health

# Free disk space (with confirmation)
tesserato clean

# Aggressive cleanup
tesserato clean --all -y
```

Pass `--config /etc/tesserato.yaml` (or any path) to any command to override config discovery.

## Project structure

```
cli_go/
├── cmd/
│   └── tesserato/
│       └── main.go              Thin entrypoint
├── internal/
│   ├── cli/
│   │   ├── root.go              Cobra root, global flags, version
│   │   ├── status.go            tesserato status
│   │   ├── up.go                tesserato up
│   │   ├── down.go              tesserato down
│   │   ├── logs.go              tesserato logs
│   │   ├── clean.go             tesserato clean (+ confirm helper)
│   │   ├── health.go            tesserato health
│   │   ├── deploy.go            tesserato deploy
│   │   ├── helpers.go           lookupProject, indexByName, etc.
│   │   └── print.go             printPerProject (shared table renderer)
│   ├── config/
│   │   ├── config.go            YAML schema, Load, Resolve, ErrNotFound
│   │   └── config_test.go
│   ├── docker/
│   │   ├── docker.go            ListContainers, ListAll, RunCompose, RunDocker
│   │   └── docker_test.go
│   └── proc/
│       └── proc.go              Generic streaming process runner
├── docker-compose.yml           Sample compose file for local testing
├── tesserato.yml                Sample tesserato config
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

The split between `cmd/` and `internal/` follows the standard Go layout: `cmd/<name>/main.go` per binary, `internal/` for code that should not be importable from outside this module.

## Development

### Run from source

```powershell
go run ./cmd/tesserato status
```

### Tests

```powershell
go test ./...
go test -v ./...        # verbose, showing each subtest
```

The current suite covers `internal/config` (Load + Resolve sentinel cases) and `internal/docker` (the `parseContainers` JSON-line parser, exercised through an injected `io.Reader` so it does not require Docker to be running). Both packages use table-driven tests with `t.Run` subtests — the idiomatic Go pattern.

### Static checks

```powershell
gofmt -l ./internal ./cmd     # lists files needing formatting
go vet ./...                  # built-in linter
```

### Reinstall after changes

```powershell
go install ./cmd/tesserato
```

Re-run any time after editing — `go install` rebuilds and replaces the binary on `PATH` in seconds.

## Roadmap

Done:

- Phase 1 MVP: `status`, `up`, `down`, `logs`.
- Phase 2: `clean`, `health`.
- Phase 3: `deploy` (git pull + rebuild).
- Polish: real binary via `go install`, version injection via ldflags, friendly multi-line error messages with available-projects suggestions, first tests, refactor to `internal/proc` and shared `printPerProject`.

Planned:

- Concurrent inspections in `health` using goroutines + `errgroup`.
- End-to-end tests for the `cli` package using `cobra.Command.SetArgs` and a mockable docker layer.
- GitHub Actions CI: tests, `go vet`, `golangci-lint`, cross-platform release binaries with version embedded.
- Colored / styled output for `status` and `health` (likely `fatih/color` or `lipgloss`).
- `--json` flag on `status` and `health` for shell-script consumption.

## Why "tesserato"?

A *tesseract* is a four-dimensional analogue of a cube — a structure that contains many cubes in a single coherent object. The same idea applies here: many Docker projects, one place to run them.
