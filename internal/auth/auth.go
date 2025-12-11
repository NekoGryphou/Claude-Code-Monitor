package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"claude-monitor/internal/consts"
)

// credentialsFile mirrors the expected JSON schema for the credentials file.
type credentialsFile struct {
	ClaudeOauth struct {
		AccessToken string `json:"accessToken"`
	} `json:"claudeAiOauth"`
}

// ResolveToken finds the OAuth token from the environment or a credentials file.
//
// Parameters:
//   - credPath: path to credentials JSON used when no environment token is present.
//
// Returns:
//   - token string if found.
//   - error when the token cannot be resolved.
func ResolveToken(credPath string) (string, error) {
	if token := strings.TrimSpace(os.Getenv(consts.EnvTokenName)); token != "" {
		return token, nil
	}

	path := expandHome(credPath)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf(consts.ErrReadCredentialsFmt, fmt.Errorf("%s: %w", path, err))
	}

	if info, err := os.Stat(path); err == nil {
		if mode := info.Mode(); mode&fs.ModePerm != 0 {
			if mode.Perm()&0o077 != 0 {
				return "", fmt.Errorf("credentials file %s must not be group/other readable (mode %v); set chmod 600", path, mode.Perm())
			}
		}
	}

	var creds credentialsFile
	if err := json.Unmarshal(content, &creds); err != nil {
		return "", fmt.Errorf(consts.ErrParseCredentialsFmt, fmt.Errorf("%s: %w", path, err))
	}

	token := strings.TrimSpace(creds.ClaudeOauth.AccessToken)
	if token == "" {
		return "", errors.New(consts.ErrEmptyAccessToken)
	}
	return token, nil
}

// DefaultCredPath returns the default credentials location, expanding the user's home directory when available.
//
// Returns:
//   - absolute path when the home directory is known; otherwise a relative default path.
func DefaultCredPath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, consts.DefaultCredRelPath)
	}
	return consts.DefaultCredRelPath
}

// expandHome replaces a leading "~" with the user's home directory when it
// can be determined; otherwise it returns the original path.
func expandHome(path string) string {
	if strings.HasPrefix(path, consts.TildePrefix) {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, consts.TildePrefix))
		}
	}
	return path
}
