package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestMCPServerLifecycle(t *testing.T) {
	in := bytes.NewBufferString(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}` + "\n" +
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}` + "\n")
	var out bytes.Buffer

	server := NewServer(in, &out)
	err := server.Serve(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(lines))
	}

	var resp1, resp2 Response
	if err := json.Unmarshal([]byte(lines[0]), &resp1); err != nil {
		t.Fatalf("unmarshal resp1 error: %v", err)
	}
	if err := json.Unmarshal([]byte(lines[1]), &resp2); err != nil {
		t.Fatalf("unmarshal resp2 error: %v", err)
	}

	resMap, ok := resp1.Result.(map[string]any)
	if !ok || resMap["serverInfo"] == nil {
		t.Errorf("expected serverInfo in init response")
	}

	toolsMap, ok := resp2.Result.(map[string]any)
	if !ok || toolsMap["tools"] == nil {
		t.Errorf("expected tools list in response")
	}
}
