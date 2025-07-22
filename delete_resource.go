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

// cleanupAWSResources performs a series of AWS cleanup operations.
// It takes the IAM role name, S3 bucket name, and KMS key ARN as input.
func cleanupAWSResources(roleName, bucketName, clusterId, mcName string) error {

	// --- IAM Operations ---

	// List all policies in the role
	rolePolicyListCmd := "aws iam list-attached-role-policies --role-name rosa-hcp-bkp-" + mcName + "-" + clusterId + " | awk '{print $2}'"

	rolePolicyListOutput, err := runCommand("bash", "-c", rolePolicyListCmd)
	if err != nil {
		return fmt.Errorf("Policies are empty or role does not exist: %s", err)
	}
	fmt.Printf("Role policies list is as follows...\n")

	var detachIAMRolePolicyCmd, deleteRoleCmd string
	for _, item := range rolePolicyListOutput {
		// 1. Detach IAM Role Policy
		detachIAMRolePolicyCmd = "aws iam detach-role-policy --policy-arn " + item + " --role-name rosa-hcp-bkp-" + mcName + "-" + clusterId

		detachIAMRolePolicyOutput, err := runCommand("bash", "-c", detachIAMRolePolicyCmd)
		if err != nil {
			return fmt.Errorf("Policies are empty or role does not exist: %s", err)
		}
		fmt.Printf("Role policy %s is detached successfully.\n", detachIAMRolePolicyOutput)
	}

	// 2. Delete IAM Role
	deleteRoleCmd = "aws iam delete-role --role-name " + roleName
	deleteRoleCmdOutput, err := runCommand("bash", "-c", deleteRoleCmd)
	if err != nil {
		return fmt.Errorf("failed to delete IAM role '%s': %w", roleName, err)
	}
	fmt.Printf("Successfully deleted IAM role '%s'.\n", deleteRoleCmdOutput)

	// 3. Delete S3 Bucket (and its contents first)
	fmt.Printf("Attempting to delete S3 bucket '%s'...\n", bucketName)
	s3Cmd := "aws s3 rb s3://" + bucketName + " --force"
	// Execute the s3 command
	s3CmdOutput, err := runCommand("bash", "-c", s3Cmd)
	if err != nil {
		return fmt.Errorf("failed to execute yq command: %s", err)
	}
	// List and delete all objects in the bucket first
	fmt.Printf("Deleting all objects in bucket '%s'...\n", s3CmdOutput)

	fmt.Printf("Successfully deleted S3 bucket '%s'.\n", bucketName)

	return nil
}

func main() {
	// os.Args[1] -> Path to backup_resources.yml file
	// os.Args[2] -> role name
	// os.Args[3] -> bucket name
	// os.Args[4] -> clusterId
	// os.Args[5] -> mc name

	fmt.Println("------Delete resources-------")
	deleteResource("bsl", "BackupStorageLocation", os.Args[1])
	deleteResource("schedule", "Schedule", os.Args[1])
	deleteResource("backup", "BackupStorageLocation", os.Args[1])
	deleteResource("secret", "Secret", os.Args[1])
	cleanupAWSResources(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
}
