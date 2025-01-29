package pkg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func IsGitRepo() bool {
	_, err := os.Stat(".git")
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return true
}

func GetGithubRepositoryUrl() (string, error) {
	// prevent dubious ownership in containers
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Add the current directory as a safe directory
	cmdSafeDir := exec.Command("git", "config", "--global", "--add", "safe.directory", currentDir)
	var stderrSafeDir bytes.Buffer
	cmdSafeDir.Stderr = &stderrSafeDir
	err = cmdSafeDir.Run()
	if err != nil {
		return "", fmt.Errorf("failed to add safe directory: %s", strings.TrimSpace(stderrSafeDir.String()))
	}

	// get the remote url
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get Git origin: %s", strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(out.String()), nil
}
