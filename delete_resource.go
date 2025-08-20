package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var initialListCmd string
var initialListArgs []string

type BackupStorageLocation struct {
	Spec Spec `json:"spec"`
}

// Spec matches the nested "spec" object.
type Spec struct {
	ObjectStorage ObjectStorage `json:"objectStorage"`
}

// ObjectStorage matches the nested "objectStorage" object.
type ObjectStorage struct {
	Bucket string `json:"bucket"`
}

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

// deleteBackup removes items from the original list that are present in the itemsToRemove list.
// It deletes the list items.
func deleteResource(resourceToDelete, clusterId string) {

	var backupList []string
	//var itemsToRemove []string
	initialListCmd = "oc"
	if resourceToDelete == "schedule" || resourceToDelete == "bsl" {
		var cmdArgs []string
		var cmdCommand = "oc"
		initialListArgs = []string{"get", resourceToDelete, clusterId + "-hourly"}
		fmt.Println("--- Getting Initial List ---", initialListArgs)
		itemsToRemove, err := runCommand(initialListCmd, initialListArgs...)
		if err != nil {
			fmt.Printf("Error getting initial list: %v\n", err)
			return
		}
		fmt.Printf("The entire list of items to remove: %s", itemsToRemove[1])
		cmdArgs = []string{"delete", resourceToDelete, itemsToRemove[1]}
		cmd := exec.Command(cmdCommand, cmdArgs...)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Errorf("failed to get StdoutPipe: %w", err)
		}
		fmt.Printf("All good: %w", stdout)
		fmt.Printf("The cmd command is: %s", cmd)
	} else {
		initialListArgs = []string{"get", resourceToDelete}
		itemsToRemove, err := runCommand(initialListCmd, initialListArgs...)
		for _, item := range itemsToRemove {
			if strings.Contains(item, clusterId) {
				backupList = append(backupList, item)
			}
		}
		if err != nil {
			fmt.Printf("Error getting initial list: %v\n", err)
		}
		fmt.Printf("The entire list of backup items to remove: %s", backupList)
		var cmdArgs []string
		var cmdCommand = "oc"
		for _, item := range backupList {
			cmdArgs = []string{"delete", resourceToDelete, item}
			cmd := exec.Command(cmdCommand, cmdArgs...)
			fmt.Println(cmd)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Errorf("failed to get StdoutPipe: %w", err)
			}
			log.Println(string(stdout))
			fmt.Printf("Resource deleted: %w", item)

		}
	}
}

// cleanupAWSResources performs a series of AWS cleanup operations.
// It takes the IAM role name, S3 bucket name, and KMS key ARN as input.
func cleanupAWSResources(clusterId, mcName string) error {
	// --- Get S3 bucket name ---
	initialListCmd = "oc"
	initialListArgs = []string{"get", "bsl", clusterId + "-hourly", "-o", "json"}

	cmd := exec.Command(initialListCmd, initialListArgs...)
	fmt.Println(initialListArgs)
	bucketNameOut, err := cmd.Output()
	if err != nil {
		fmt.Errorf("failed to get StdoutPipe: %w", err)
	}

	fmt.Println("S3 bucket name JSON is ", string(bucketNameOut))
	var bsl BackupStorageLocation

	// Unmarshal (parse) the JSON string into the struct.
	s3bucketerr := json.Unmarshal([]byte(bucketNameOut), &bsl)
	if s3bucketerr != nil {
		fmt.Println("Error parsing JSON:", err)
		//return
	}

	// Access the nested bucket value and print it.
	bucketName := bsl.Spec.ObjectStorage.Bucket
	fmt.Println(bucketName)

	// --- IAM Operations ---
	var roleNamePrefix = "rosa-hcp-bkp-"
	// List all policies in the role
	awsCmd := "aws"
	rolePolicyListArgs := []string{"iam", "list-attached-role-policies", "--role-name", roleNamePrefix + mcName + "-" + clusterId, "|", "awk", "{print $2}"}

	rolePolicyListOutput, err := runCommand(awsCmd, rolePolicyListArgs...)
	if err != nil {
		fmt.Printf("Policies are empty or role does not exist: %s", err)
	}
	fmt.Printf("Role policies list Output is as follows... %s\n", rolePolicyListOutput)

	for _, item := range rolePolicyListOutput {
		// 1. Detach IAM Role Policy
		var detachIAMRolePolicyListArgs = []string{"iam", "detach-role-policy", "--policy-arn", item, "--role-name", "rosa-hcp-bkp-" + mcName + "-" + clusterId}

		detachIAMRolePolicyOutput, err := runCommand(awsCmd, detachIAMRolePolicyListArgs...)
		if err != nil {
			fmt.Printf("Policies are empty or role does not exist: %s", err)
		}
		fmt.Printf("Role policy %s is detached successfully.\n", detachIAMRolePolicyOutput)
	}

	// 2. Delete IAM Role
	var deleteRoleListArgs = []string{"iam", "delete-role", "--role-name", roleNamePrefix + mcName + "-" + clusterId}
	deleteRoleCmdOutput, err := runCommand(awsCmd, deleteRoleListArgs...)
	if err != nil {
		fmt.Printf("failed to delete IAM role '%s': %w", roleNamePrefix+mcName+"-"+clusterId, err)
	}
	fmt.Printf("Successfully deleted IAM role '%s'.\n", deleteRoleCmdOutput)

	// 3. Delete S3 Bucket (and its contents first)
	fmt.Printf("Attempting to delete S3 bucket '%s'...\n", bucketName)
	s3CmdListArgs := []string{"s3", "rb", "s3://" + bucketName, " --force"}
	// Execute the s3 command
	s3CmdOutput, err := runCommand(awsCmd, s3CmdListArgs...)
	if err != nil {
		fmt.Printf("failed to execute yq command: %s", err)
	}
	//List and delete all objects in the bucket first
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
	var clusterId = os.Args[1]
	var mcName = os.Args[2]

	fmt.Println("------Delete Openshift resources-------")
	cleanupAWSResources(clusterId, mcName)

	fmt.Println("------Delete Openshift resources-------")
	deleteResource("bsl", clusterId)
	deleteResource("schedule", clusterId)
	deleteResource("backup", clusterId)
	deleteResource("secret", clusterId)
	deleteResource("backuprepository", clusterId)
}
