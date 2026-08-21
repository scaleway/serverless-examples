package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id,omitempty"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func main() {
	// cmp.Or eliminates verbose env var fallback checks
	port := cmp.Or(os.Getenv("PORT"), "8080")

	mux := http.NewServeMux()

	// Native method routing in Go 1.22+
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("OPTIONS /mcp", handleCORS)
	mux.HandleFunc("POST /mcp", handleMCP)

	slog.Info("Starting modern MCP server", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("Server crashed", "error", err)
		os.Exit(1)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func handleCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

func handleMCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		resp.Result = map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    "scaleway-test-mcp",
				"version": "1.1.0",
			},
		}

	case "tools/list":
		resp.Result = map[string]any{
			"tools": []map[string]any{
				{
					"name":        "get_server_time",
					"description": "Returns current UTC time from the Scaleway container",
					"inputSchema": map[string]any{
						"type":       "object",
						"properties": map[string]any{},
					},
				},
				{
					"name":        "echo",
					"description": "Echoes back the input message",
					"inputSchema": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"message": map[string]any{
								"type":        "string",
								"description": "The string to echo back",
							},
						},
						"required": []string{"message"},
					},
				},
			},
		}

	case "tools/call":
		var params struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			resp.Error = map[string]any{"code": -32602, "message": "Invalid params"}
			break
		}

		switch params.Name {
		case "get_server_time":
			resp.Result = map[string]any{
				"content": []map[string]any{
					{
						"type": "text",
						"text": fmt.Sprintf("Scaleway Container UTC Time: %s", time.Now().UTC().Format(time.RFC3339)),
					},
				},
			}
		case "echo":
			msg, _ := params.Arguments["message"].(string)
			resp.Result = map[string]any{
				"content": []map[string]any{
					{
						"type": "text",
						"text": fmt.Sprintf("Echo from Scaleway: %s", msg),
					},
				},
			}
		default:
			resp.Error = map[string]any{"code": -32601, "message": "Tool not found"}
		}

	default:
		resp.Result = map[string]any{}
	}

	_ = json.NewEncoder(w).Encode(resp)
}
