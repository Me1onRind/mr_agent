# mr_agent

## Directory Structure

```
.
├── cmd                         Application entry
│   ├── api                     API service entry
│   └── cli                     CLI service entry
└── internal
    ├── app                     Application assembly
    │   ├── api                  API wiring & bootstrap
    │   └── cli                  CLI wiring & flags
    ├── config                  Config loading & structs
    ├── domain                  Domain layer
    │   └── dialog               Dialog domain
    ├── errcode                 Error codes & wrappers
    ├── infrastructure          Infrastructure & middleware
    │   ├── agent                Agent client & tools
    │   │   ├── client            Agent client
    │   │   └── tools             Agent tool implementations
    │   ├── cache                Caching
    │   │   └── llm              LLM cache
    │   ├── db                   Database clients & registry
    │   ├── goroutine            Goroutine helpers
    │   ├── logger               Logger implementation
    │   ├── middleware           Middleware
    │   └── session              Session store (in-memory)
    ├── initialize              Initializers
    ├── pkg                     Shared utilities
    │   └── strutil              String utilities
    ├── protocol                Protocol & routing
    │   ├── agent                Agent tool protocol
    │   └── http                 HTTP protocol
    │       ├── chat             Chat HTTP routes
    │       ├── ping             Ping HTTP routes
    │       └── tools            Tools HTTP routes
    └── usecase                 Use case layer
        ├── chat                Chat use case
        ├── ping                Ping use case
        └── tools               Tools use case
```

## Features

- `Ping` endpoint: health check & availability
- `Chat` endpoint: chat entrypoint
- `Tools` endpoint: tool execution entrypoint
- `CLI` entry: dialog with LLM in terminal (readline-based, supports multi-byte deletion)
- Agent tools: tool registry and built-in implementations
- Config: local config loader and typed configs
- Database: MySQL client init and registry
- Tracing: OpenTracing initialization hook
- Middleware: access logs, JSON gateway, recovery, tracing
- Cache: LLM cache and simple cache implementation
- Session: in-memory session store for dialog context
- Core utilities: logging, error codes, goroutine helpers, string utils

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
