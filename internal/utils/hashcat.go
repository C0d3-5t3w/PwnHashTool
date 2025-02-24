package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunHashcat(hashFile, wordlist string, options []string) (string, error) {
	if _, err := exec.LookPath("hashcat"); err != nil {
		return "", fmt.Errorf("hashcat not found in PATH: %v", err)
	}

	baseName := filepath.Base(hashFile)
	potfile := strings.TrimSuffix(baseName, ".hc22000") + ".potfile"

	args := []string{"-m", "22000"}
	args = append(args, hashFile, wordlist, "--potfile-path", potfile)

	cmd := exec.Command("hashcat", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("hashcat failed: %v\nOutput: %s", err, output)
	}

	return potfile, nil
}
