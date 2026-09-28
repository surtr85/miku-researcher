package reader

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type PageResult struct {
	URL     string `json:"url"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

type Reader struct {
	defuddlePath string
	timeout      time.Duration
}

func NewReader(timeout time.Duration) *Reader {
	p, err := exec.LookPath("defuddle")
	if err != nil {
		p = "defuddle"
	}
	if timeout <= 0 {
		timeout = 25 * time.Second
	}
	return &Reader{
		defuddlePath: p,
		timeout:      timeout,
	}
}

// ReadPage fetches and extracts clean markdown using defuddle
func (r *Reader) ReadPage(ctx context.Context, targetURL string, maxChars int) (*PageResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// Run defuddle parse --markdown <url>
	cmd := exec.CommandContext(ctx, r.defuddlePath, "parse", "--markdown", targetURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return &PageResult{
			URL:   targetURL,
			Error: fmt.Sprintf("defuddle extraction error: %s", errMsg),
		}, nil
	}

	content := strings.TrimSpace(stdout.String())
	if maxChars > 0 && len(content) > maxChars {
		content = content[:maxChars] + "\n\n... [Content truncated for token efficiency] ..."
	}

	return &PageResult{
		URL:     targetURL,
		Content: content,
	}, nil
}
