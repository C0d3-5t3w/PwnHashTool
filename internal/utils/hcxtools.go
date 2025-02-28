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

func ProcessPcapDirectory(dirPath string, options []string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dirPath, "*.pcap"))
	if err != nil {
		return nil, fmt.Errorf("error finding PCAP files: %v", err)
	}

	pcapngFiles, err := filepath.Glob(filepath.Join(dirPath, "*.pcapng"))
	if err != nil {
		return nil, fmt.Errorf("error finding PCAPNG files: %v", err)
	}

	files = append(files, pcapngFiles...)
	if len(files) == 0 {
		return nil, fmt.Errorf("no PCAP/PCAPNG files found in directory")
	}

	var outputs []string
	for _, file := range files {
		output, err := RunHcxPcapngTool(file, options)
		if err != nil {
			return outputs, fmt.Errorf("error processing %s: %v", file, err)
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}
