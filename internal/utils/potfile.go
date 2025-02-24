package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParsePotfile reads a potfile and creates a .pass file with extracted passwords
func ParsePotfile(potfilePath string) (string, error) {
	// Open the potfile
	file, err := os.Open(potfilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open potfile: %v", err)
	}
	defer file.Close()

	// Create output .pass file
	outputPath := strings.TrimSuffix(potfilePath, ".potfile") + "_password.txt"
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outFile)

	// Process each line
	for scanner.Scan() {
		line := scanner.Text()
		if lastIndex := strings.LastIndex(line, ":"); lastIndex != -1 {
			// Write the password (everything after the last colon)
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
