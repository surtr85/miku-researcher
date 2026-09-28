package envutil

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// AutoLoadEnv attempts to load environment variables from well-known secret paths
// if they are not already set in the current process environment.
func AutoLoadEnv() {
	candidates := []string{
		"/run/secrets/web-tools-env",
	}

	if home, err := os.UserHomeDir(); err == nil && home != "" {
		candidates = append(candidates,
			filepath.Join(home, ".config/sops-nix/secrets/web-tools-env"),
			filepath.Join(home, ".env"),
		)
	}

	for _, p := range candidates {
		loadEnvFile(p)
	}
}

func loadEnvFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
			if _, exists := os.LookupEnv(k); !exists && k != "" {
				os.Setenv(k, v)
			}
		}
	}
}
