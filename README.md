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
make run-http  # listens on :3333
```
Call with a single JSON-RPC request per POST:
```sh
curl -X POST -H "Content-Type: application/json" \
	-d '{"id":1,"method":"tools/list","params":{}}' \
	http://localhost:3333/
```
Health check: `curl http://localhost:3333/health`

## Available tool
- `payram_intro`: Returns a plain-text overview of PayRam and useful links.
- `payram_analytics`: Calls PayRam analytics APIs. Actions:
	- `list_groups`: GET analytics groups (requires `PAYRAM_ANALYTICS_TOKEN`; `PAYRAM_ANALYTICS_BASE_URL` or `base_url` argument must be set).
	- `graph_data`: POST group/graph data. Args: `group_id` (int), `graph_id` (int), `payload` (object, optional; defaults to `{ "analytics_date_filter": "last_30_days" }`).
		Example payloads: filters like `group_by_network_currency_filter`, `in_query_currency_filter`, etc., as provided by the API.

## Chat orchestrator (UI)
Launch a minimal chat UI that routes tool calls through the MCP server (HTTP mode required):

```sh
# terminal 1
make run-http

# terminal 2
OPENAI_API_KEY=sk-... make run-chat CHAT_PORT=3000 MCP_SERVER_URL=http://localhost:3333/
```

Then open http://localhost:3000/ and ask for a PayRam intro to see the tool call in action.

Environment:
- `OPENAI_API_KEY` (required)
- `OPENAI_MODEL` (default: `gpt-4o-mini`)
- `OPENAI_BASE_URL` (default: `https://api.openai.com/v1`)

## OpenAI-compatible chat API
Expose `/v1/chat/completions` that routes tool calls to the MCP server, with a PayRam system prompt.

```sh
# terminal 1: MCP HTTP
make run-http

# terminal 2: Chat API
CHAT_API_KEY=secret OPENAI_API_KEY=sk-... make run-chat-api CHAT_API_PORT=2358 MCP_SERVER_URL=http://localhost:3333/
```

Call example:
```sh
curl -X POST http://localhost:2358/v1/chat/completions \
	-H "Authorization: Bearer secret" \
	-H "Content-Type: application/json" \
	-d '{
		"model": "gpt-4o-mini",
		"messages": [{"role":"user","content":"List analytics groups"}]
	}'
```

Configuration (.env or env vars):
- `CHAT_API_KEY` (required for auth)
- `OPENAI_API_KEY` (required), `OPENAI_MODEL` (default `gpt-4o-mini`), `OPENAI_BASE_URL` (default `https://api.openai.com/v1`)
- `MCP_SERVER_URL` (HTTP endpoint for MCP server; default `http://localhost:3333/`)

## Structure
- `main.go`: wires stdin/stdout loop to the MCP server.
- `internal/mcp`: server routing, toolbox, and protocol handling.
- `internal/tools`: individual tool implementations (extensible).
- `internal/protocol`: shared types for requests/responses.

## Development
- Format: `make fmt`
- Test (no tests yet, but ensures build succeeds): `make test`

## Updates and releases
- Secrets: set repository secret `PAYRAM_UPDATE_ED25519_PRIVKEY_B64` to the base64-encoded 64-byte Ed25519 private key used to sign manifests (public key is logged during the workflow run).
- Tagging: push tags `vX.Y.Z` for stable and `vX.Y.Z-beta.N` for beta. The release workflow builds linux/amd64 binaries `chat.bin` (from `cmd/chat-api`) and `mcp.bin` (from `cmd/mcp-server`), uploads them as GitHub Release assets, generates `updates/<channel>/manifest.json` + `.sig`, and commits them back to `updates/` on the default branch.
- Agent config for zero-infra updates:
	- `PAYRAM_AGENT_UPDATE_BASE_URL=https://raw.githubusercontent.com/PayRam/analytics-mcp-server/main/updates`
	- `PAYRAM_AGENT_UPDATE_CHANNEL=stable` (or `beta`)
	- `PAYRAM_AGENT_UPDATE_PUBKEY_B64=<matching Ed25519 public key>`
