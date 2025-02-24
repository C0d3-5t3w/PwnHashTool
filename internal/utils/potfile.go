package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ParsePotfile(potfilePath string) (string, error) {
	file, err := os.Open(potfilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open potfile: %v", err)
	}
	defer file.Close()

	outputPath := strings.TrimSuffix(potfilePath, ".potfile") + "_password.txt"
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outFile)

	for scanner.Scan() {
		line := scanner.Text()
		if lastIndex := strings.LastIndex(line, ":"); lastIndex != -1 {
			password := line[lastIndex+1:]
			if _, err := writer.WriteString(password + "\n"); err != nil {
				return "", fmt.Errorf("failed to write to output file: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading potfile: %v", err)
	}

	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("error flushing output file: %v", err)
	}

	return outputPath, nil
}

// Author: C0d3-5t3w