package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runCommand captures the stdout of a given command and returns it as a slice of strings,
// where each string is a line from the command's output.
func runCommand(command string, args ...string) ([]string, error) {
	fmt.Printf("Executing command: %s %s\n", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...) // Create the command
	stdout, err := cmd.StdoutPipe()       // Get a pipe to read the standard output
	if err != nil {
		return nil, fmt.Errorf("failed to get StdoutPipe: %w", err)
	}

	if err := cmd.Start(); err != nil { // Start the command
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	//var outputLines []string
	var firstColumnList []string
	scanner := bufio.NewScanner(stdout) // Create a scanner to read line by line
	for scanner.Scan() {
		line := scanner.Text() // Read a line and trim whitespace
		// if line != "" {                           // Only add non-empty lines
		// 	outputLines = append(outputLines, line)
		// }
		fields := strings.Fields(line) // Split the line by whitespace
		if len(fields) > 0 {
			firstColumnList = append(firstColumnList, fields[0]) // Add the first field to the list
		}
	}

	if err := scanner.Err(); err != nil { // Check for scanner errors
		return nil, fmt.Errorf("error reading command output: %w", err)
	}

	if err := cmd.Wait(); err != nil { // Wait for the command to finish and check for errors
		return nil, fmt.Errorf("command exited with error: %w", err)
	}

	return firstColumnList, nil
}

// getMetadataNameFromYAML extracts the metadata.name for a given kind from a YAML file.
func getMetadataNameFromYAML(kind, filePath string) (string, error) {
	fmt.Printf("\n--- Extracting metadata.name for kind '%s' from '%s' ---\n", kind, filePath)

	// Construct the yq command
	yqCmd := fmt.Sprintf("yq '. | select(.kind == \"%s\") | .metadata.name' %s", kind, filePath)

	// Execute the yq command
	yqCmdOutput, err := runCommand("bash", "-c", yqCmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute yq command: %s", err)
	}

	extractedName := strings.TrimSpace(yqCmdOutput[0])
	if extractedName == "" {
		return "", fmt.Errorf("no metadata.name found for kind '%s' in file '%s'", kind, filePath)
	}

	fmt.Printf("Extracted name: %s\n", extractedName)
	fmt.Println("--- Metadata Name Extraction Completed ---")
	return extractedName, nil
}

// deleteBackup removes items from the original list that are present in the itemsToRemove list.
// It deletes the list items.
func deleteResource(resourceToDelete, kind, filePath string) {
	var initialListCmd string
	var initialListArgs []string

	extractName, err := getMetadataNameFromYAML(kind, filePath)
	if resourceToDelete != "backup" {
		initialListCmd = "oc"
		initialListArgs = []string{"get", resourceToDelete, extractName}
	} else if strings.Contains(extractName, extractName) {
		initialListCmd = "oc"
		initialListArgs = []string{"get", resourceToDelete}
	} else {
		fmt.Println("Continuing deleting openshift resources.")
	}

	fmt.Println("--- Getting Initial List ---")
	itemsToRemove, err := runCommand(initialListCmd, initialListArgs...)
	if err != nil {
		//fmt.Printf("Error getting initial list: %v\n", err)
		return
	}
	// Create a map for quick lookup of items to remove
	//removeMap := make(map[string]struct{})
	var cmdCommand string
	var cmdArgs []string

	cmdCommand = "oc"

	for _, item := range itemsToRemove {
		cmdArgs = []string{"delete", resourceToDelete, item}
		cmd := exec.Command(cmdCommand, cmdArgs...)
		fmt.Println(cmd)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Errorf("failed to get StdoutPipe: %w", err)
		}
		fmt.Println("All good: %w", stdout)

	}
}

func main() {

	fmt.Println("------Delete resources-------")
	deleteResource("bsl", "BackupStorageLocation", os.Args[1])
	deleteResource("schedule", "Schedule", os.Args[1])
	deleteResource("backup", "BackupStorageLocation", os.Args[1])
	deleteResource("secret", "Secret", os.Args[1])
}
