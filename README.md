# PayRam Analytics MCP Server

A Go-based Model Context Protocol (MCP) server with a modular toolbox for tools. Currently exposes one tool that returns an introduction to PayRam.

## Prerequisites
- Go 1.21+

## Run (stdio)
```sh
make run
```
Send JSON-RPC 2.0 requests over stdin. Example session:

1) Initialize
```json
{"id":1,"method":"initialize","params":{}}
```
2) List tools
```json
{"id":2,"method":"tools/list","params":{}}
```
3) Call the PayRam intro tool
```json
{"id":3,"method":"tools/call","params":{"name":"payram_intro"}}
```

## Run (HTTP)
```sh
make run-http  # listens on :8080
```
Call with a single JSON-RPC request per POST:
```sh
curl -X POST -H "Content-Type: application/json" \
	-d '{"id":1,"method":"tools/list","params":{}}' \
	http://localhost:8080/
```
Health check: `curl http://localhost:8080/health`

## Available tool
- `payram_intro`: Returns a plain-text overview of PayRam and useful links.

## Structure
- `main.go`: wires stdin/stdout loop to the MCP server.
- `internal/mcp`: server routing, toolbox, and protocol handling.
- `internal/tools`: individual tool implementations (extensible).
- `internal/protocol`: shared types for requests/responses.

## Development
- Format: `make fmt`
- Test (no tests yet, but ensures build succeeds): `make test`
