package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/payram/payram-analytics-mcp-server/internal/mcp"
	"github.com/payram/payram-analytics-mcp-server/internal/protocol"
	"github.com/payram/payram-analytics-mcp-server/internal/tools"
)

func main() {
	httpAddr := flag.String("http", "", "HTTP listen address (e.g., :8080). If set, server runs over HTTP instead of stdio.")
	flag.Parse()

	if err := run(*httpAddr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(httpAddr string) error {
	ctx := context.Background()
	tb := mcp.NewToolbox(tools.PayramIntro())
	server := mcp.NewServer(tb)

	if httpAddr != "" {
		log.Printf("starting HTTP MCP server on %s", httpAddr)
		return mcp.RunHTTP(server, httpAddr)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		raw, err := readFrame(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("read error: %w", err)
		}
		if raw == "" {
			continue
		}

		var req protocol.Request
		if err := json.Unmarshal([]byte(raw), &req); err != nil {
			emitResponse(mcp.WriteError(defaultID(), -32700, "invalid JSON", err))
			continue
		}

		resp, err := server.Handle(ctx, req)
		if err != nil {
			emitResponse(mcp.WriteError(req.ID, -32603, "internal error", err))
			continue
		}
		emitResponse(resp)
	}
}

func readFrame(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return "", nil
	}

	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "content-length:") {
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("malformed Content-Length header")
		}
		n, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return "", fmt.Errorf("invalid Content-Length: %w", err)
		}
		// consume until blank line
		for {
			h, err := r.ReadString('\n')
			if err != nil {
				return "", err
			}
			if strings.TrimSpace(h) == "" {
				break
			}
		}
		body := make([]byte, n)
		nread, err := io.ReadFull(r, body)
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) && nread > 0 {
				return strings.TrimSpace(string(body[:nread])), nil
			}
			return "", err
		}
		return strings.TrimSpace(string(body)), nil
	}

	return trimmed, nil
}

func defaultID() any {
	return "0"
}

func emitResponse(resp protocol.Response) {
	body, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode response: %v\n", err)
		return
	}
	// Emit newline-delimited JSON; inspector stdio expects raw JSON per line.
	if _, err := os.Stdout.Write(append(body, '\n')); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write response: %v\n", err)
	}
}
