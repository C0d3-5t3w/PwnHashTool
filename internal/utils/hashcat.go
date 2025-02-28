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

func ProcessHashDirectory(dirPath, wordlist string, options []string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dirPath, "*.hc22000"))
	if err != nil {
		return nil, fmt.Errorf("error finding hash files: %v", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no .hc22000 files found in directory")
	}

	var outputs []string
	for _, file := range files {
		output, err := RunHashcat(file, wordlist, options)
		if err != nil {
			return outputs, fmt.Errorf("error processing %s: %v", file, err)
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}
