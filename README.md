# mr_agent

## Directory Structure

```
.
├── cmd                         Application entry
│   └── api                     API service entry
│   └── cli                     CLI service entry
└── internal
    ├── app                     Application assembly
    │   └── api                  API wiring & bootstrap
    │   └── cli                  CLI wiring & bootstrap
    ├── domain                  Domain layer
    │   └── dialog               Dialog domain
    ├── errcode                 Error codes & wrappers
    ├── infrastructure          Infrastructure & middleware
    │   ├── cache                Caching
    │   │   └── llm              LLM cache
    │   ├── client               Client wrappers
    │   │   └── agent            Agent client
    │   ├── goroutine            Goroutine helpers
    │   ├── logger               Logger implementation
    │   └── middleware           Middleware
    ├── initialize              Initializers
    ├── pkg                     Shared utilities
    │   └── strutil              String utilities
    ├── protocol                Protocol & routing
    │   └── http                 HTTP protocol
    │       ├── chat             Chat HTTP routes
    │       └── ping             Ping HTTP routes
    └── usecase                 Use case layer
        ├── chat                Chat use case
        └── ping                Ping use case
```

## Features

- `Ping` endpoint: health check & availability
- `Chat` endpoint: chat entrypoint
- `CLI` entry: dialog with LLM in terminal
- Middleware: access logs, JSON gateway, recovery, tracing
- Cache: LLM cache and simple cache implementation
- Core utilities: logging, error codes, goroutine helpers, client wrappers

## Run
### API Service

```
make http_run
```

Or run directly:

```
go run ./cmd/api/main.go
```

### CLI Service

```
go run ./cmd/cli/main.go
```
