package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunHcxPcapngTool(inputFile string, options []string) (string, error) {
	if _, err := exec.LookPath("hcxpcapngtool"); err != nil {
		return "", fmt.Errorf("hcxpcapngtool not found in PATH: %v", err)
	}

	// Create output filename based on input file
	baseName := filepath.Base(inputFile)
	outputFile := strings.TrimSuffix(baseName, filepath.Ext(baseName)) + ".hc22000"

	args := []string{"-o", outputFile}
	args = append(args, options...)
	args = append(args, inputFile)

	cmd := exec.Command("hcxpcapngtool", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("hcxpcapngtool failed: %v\nOutput: %s", err, output)
	}

	return outputFile, nil
}
