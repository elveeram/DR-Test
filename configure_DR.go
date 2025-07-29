package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

// runCommand executes a shell command and returns its stdout and stderr.
// It also prints the command being executed for clarity.
func runCommand(name string, arg ...string) (string, string, error) {
	cmd := exec.Command(name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Print the command being executed
	fmt.Printf("Executing command: %s %s\n", name, strings.Join(arg, " "))

	err := cmd.Run()
	if err != nil {
		// Return a detailed error message including stderr
		return stdout.String(), stderr.String(), fmt.Errorf("command failed: %w, stderr: %s", err, stderr.String())
	}
	return stdout.String(), stderr.String(), nil
}

// setupCluster performs the sequence of steps to set up and configure a cluster.
// It takes the cluster ID, name, and environment as input.
// Assuming that there is a hive cluster and the kubeconfigs (Management and Service Clusters) are already saved on users local.
func setupCluster(clusterID, clusterName, clusterEnv string) error {
	fmt.Println("--- Cluster Setup Started ---")

	// Step 1: Simulate environment variables (for reference)
	// In Go, 'export' directly to the shell is not possible.
	// These values are used internally by the program.
	fmt.Printf("\nStep 1: Setting up cluster environment variables (for reference):\n")
	fmt.Printf("  export cluster_id=%s\n", clusterID)
	fmt.Printf("  export cluster_name=%s\n", clusterName)
	fmt.Printf("  export cluster_env=%s\n", clusterEnv)

	clusterNameerr := os.Setenv("cluster_name", clusterName)
	if clusterNameerr != nil {
		// Log the error and stderr if the command fails
		return fmt.Errorf("failed to create cluster name variable: %w, stderr: %s", clusterNameerr)
	}
	fmt.Println("Cluster Name environment variable created: %s\n", os.Getenv("cluster_name"))

	clusterIDerr := os.Setenv("cluster_id", clusterID)
	if clusterIDerr != nil {
		// Log the error and stderr if the command fails
		return fmt.Errorf("failed to create cluster name variable: %w, stderr: %s", clusterIDerr)
	}
	fmt.Println("Cluster Id environment variable created: %s\n", os.Getenv("cluster_id"))

	clusterEnverr := os.Setenv("cluster_env", clusterEnv)
	if clusterEnverr != nil {
		// Log the error and stderr if the command fails
		return fmt.Errorf("failed to create cluster Env variable: %w, stderr: %s", clusterEnverr)
	}
	fmt.Println("Cluster Env environment variable created: %s\n", os.Getenv("cluster_env"))

	// Step 2: Check cluster health using 'rosa describe cluster'
	fmt.Printf("\nStep 2: Checking cluster health for '%s'...\n", clusterName)
	stdout, stderr, err := runCommand("rosa", "describe", "cluster", "--cluster="+clusterName)
	if err != nil {
		// Log the error and stderr if the command fails
		return fmt.Errorf("failed to check cluster health: %w, stderr: %s", err, stderr)
	}
	fmt.Println("Cluster Health Output:\n", stdout)

	// Basic check for cluster state. 'ready' is ideal, 'installing' is also acceptable.
	if strings.Contains(stdout, "ready") || strings.Contains(stdout, "installing") {
		fmt.Println("Cluster appears healthy or is in the installation process.")
	} else {
		fmt.Println("Warning: Cluster state is not 'ready' or 'installing'. Please review the output above.")
	}

	// Step 3: Create an admin for the cluster using 'rosa create admin'
	// fmt.Printf("\nStep 3: Creating admin for cluster '%s'...\n", clusterName)
	// adminStdout, adminStderr, err := runCommand("rosa", "create", "admin", "--cluster="+clusterName)
	// if err != nil {
	// 	// Log the error and stderr if the command fails
	// 	return fmt.Errorf("failed to create cluster admin: %w, stderr: %s", err, adminStderr)
	// }
	// fmt.Println("Admin Creation Output:\n", adminStdout)

	//oc_login_cmd := strings.Split(adminStdout, "\n")[3]

	//stdout_cluster, stderr_cluster, err_cluster := runCommand(oc_login_cmd)

	// fmt.Println("Cluster Login output:\n", stdout_cluster)
	// if err_cluster != nil {
	// 	// Log the error and stderr if the command fails
	// 	return fmt.Errorf("failed to login to cluster: %w, stderr: %s", err_cluster, stderr_cluster)
	// }

	// Important instruction for the user regarding 'oc login'
	// fmt.Println("\n--- IMPORTANT ---")
	// fmt.Println("Please save the 'oc login' command from the output above.")
	// fmt.Println("You MUST open a separate terminal and execute that 'oc login' command to log into your cluster.")
	// fmt.Println("This step is crucial for the subsequent 'oc create ns' command to work.")
	// fmt.Println("-----------------\n")

	// Step 4: Create a test namespace called 'nginx'
	// This step assumes the user has manually logged into the cluster
	// in a separate terminal using the credentials from Step 3.
	// fmt.Printf("\nStep 4: Creating 'nginx' namespace...\n")
	// fmt.Println("Note: This step assumes you have successfully logged into your cluster")
	// fmt.Println("      in a separate terminal using the 'oc login' command provided previously.")

	// nsStdout, nsStderr, err := runCommand("oc", "create", "ns", "nginx")
	// if err != nil {
	// 	// Provide guidance if namespace creation fails, likely due to not being logged in
	// 	fmt.Printf("Warning: Failed to create 'nginx' namespace. This often happens if 'oc' is not logged in or doesn't have the correct cluster context.\n")
	// 	fmt.Printf("Error: %v\nStderr: %s\n", err, nsStderr)
	// 	fmt.Println("Please ensure you are logged into the cluster in your terminal and try running 'oc create ns nginx' manually.")
	// 	return fmt.Errorf("failed to create 'nginx' namespace (check 'oc' login status): %w", err)
	// }
	// fmt.Println("Namespace Creation Output:\n", nsStdout)
	// fmt.Println("Namespace 'nginx' created successfully (or already exists).")

	// fmt.Println("\n--- Cluster Setup Completed ---")
	// fmt.Println("\n--- Logout of Cluster ---")
	// unsetKCStdout, unsetKCStderr, err := runCommand("unset", "KUBECONFIG")
	// if err != nil {
	// 	fmt.Printf("Unset KUBECONFIG command failed: %s\n", unsetKCStderr)
	// }
	// fmt.Println("KUBECONFIG unset complete %s\n", unsetKCStdout)
	// logoutStdout, logoutStderr, err := runCommand("oc", "logout")
	// if err != nil {
	// 	fmt.Println("Unset KUBECONFIG command failed: %s\n", logoutStderr)
	// }
	// fmt.Println("KUBECONFIG unset complete %s\n", logoutStdout)
	return nil
}

// createS3Bucket performs the steps to create an AWS S3 bucket.
// It sets the AWS_PROFILE and generates a unique bucket name.
func createS3Bucket(awsProfile, region string) error {
	fmt.Println("\n--- AWS S3 Bucket Creation Started ---")

	// Step 1: Generate a unique bucket name using uuidgen
	fmt.Println("Step 1: Generating unique bucket name...")
	uuidStdout, uuidStderr, err := runCommand("uuidgen")
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w, stderr: %s", err, uuidStderr)
	}
	// Clean up the UUID: remove hyphens and convert to lowercase
	generatedUUID := strings.TrimSpace(uuidStdout)
	cleanedUUID := strings.ReplaceAll(generatedUUID, "-", "")
	bucketName := fmt.Sprintf("rosa-hcp-backup-oadp-%s", strings.ToLower(cleanedUUID))
	fmt.Printf("Generated bucket name: %s\n", bucketName)

	// Step 2: Set AWS_PROFILE and create the S3 bucket
	fmt.Printf("Step 2: Creating S3 bucket '%s' in region '%s' with AWS_PROFILE='%s'...\n", bucketName, region, awsProfile)

	cmd_awsProfile := exec.Command("export", "AWS_PROFILE=", awsProfile)
	fmt.Println(cmd_awsProfile)

	// Execute the aws s3api create-bucket command
	s3Stdout, s3Stderr, err := runCommand(
		"aws", "s3api", "create-bucket",
		"--bucket", bucketName,
		"--region", region,
		"--create-bucket-configuration", fmt.Sprintf("LocationConstraint=%s", region))

	if err != nil {
		return fmt.Errorf("failed to create S3 bucket: %w, stderr: %s", err, s3Stderr)
	}
	fmt.Println("S3 Bucket Creation Output:\n", s3Stdout)
	fmt.Printf("S3 bucket '%s' created successfully.\n", bucketName)

	fmt.Println("--- AWS S3 Bucket Creation Completed ---")
	return nil
}

// createOIDCConfig retrieves the OIDC endpoint URL and extracts the OIDC ID.
func createOIDCConfig(mcName, region, clusterId string) (string, string, string, error) {
	fmt.Println("\n--- OIDC Configuration Started ---")

	// Step 1: Get the management cluster reference href
	fmt.Printf("Step 1: Getting OIDC endpoint URL for management cluster '%s' in region '%s'...\n", mcName, region)
	// The ocm get command with jq needs to be executed via bash -c to handle pipes
	ocmGetHrefCmd := fmt.Sprintf(`ocm get /api/osd_fleet_mgmt/v1/management_clusters -p search="region='%s' and name='%s'" | jq -r '.items[].cluster_management_reference.href'`, region, mcName)
	hrefStdout, hrefStderr, err := runCommand("bash", "-c", ocmGetHrefCmd)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get management cluster href: %w, stderr: %s", err, hrefStderr)
	}
	mcHref := strings.TrimSpace(hrefStdout)
	if mcHref == "" {
		return "", "", "", fmt.Errorf("management cluster href not found for %s in %s", mcName, region)
	}
	fmt.Printf("Management Cluster Href: %s\n", mcHref)

	// Step 2: Get the OIDC endpoint URL using the href
	ocmGetOIDCUrlCmd := fmt.Sprintf(`ocm get %s | jq -r '.aws.sts.oidc_endpoint_url'`, mcHref)
	oidcUrlStdout, oidcUrlStderr, err := runCommand("bash", "-c", ocmGetOIDCUrlCmd)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get OIDC endpoint URL: %w, stderr: %s", err, oidcUrlStderr)
	}
	mcOIDCUrl := strings.TrimSpace(oidcUrlStdout)
	if mcOIDCUrl == "" {
		return "", "", "", fmt.Errorf("OIDC endpoint URL not found")
	}
	fmt.Printf("mc_oidc_url: %s\n", mcOIDCUrl)

	// Step 3: Extract the OIDC ID by removing "https://"
	mcOIDCCmd := fmt.Sprintf(`echo %s | sed -r 's/https:\/\///'`, mcOIDCUrl)
	oidcStdout, oidcStderr, err := runCommand("bash", "-c", mcOIDCCmd)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to extract OIDC ID: %w, stderr: %s", err, oidcStderr)
	}
	mcOIDC := strings.TrimSpace(oidcStdout)
	if mcOIDC == "" {
		return "", "", "", fmt.Errorf("OIDC ID could not be extracted")
	}
	fmt.Printf("mc_oidc: %s\n", mcOIDC)

	// Step 4: Get OIDC ARN
	mcOIDCArnCmd := fmt.Sprintf("aws iam list-open-id-connect-providers --output json | jq -r '.OpenIDConnectProviderList[].Arn' | grep 2juqamhdcrm6o3f7i5l309lgnu72bmb0")
	oidcArnStdout, oidcArnStderr, err := runCommand("bash", "-c", mcOIDCArnCmd)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to extract OIDC ID: %w, stderr: %s", err, oidcArnStderr)
	}
	fmt.Printf("mc_oidc_arn: %s\n", oidcArnStdout)
	fmt.Println("--- OIDC Configuration Completed ---")
	return mcOIDCUrl, mcOIDC, oidcArnStdout, nil
}

// createIAMRole creates an IAM role and attaches a policy.
func createIAMRole(awsProfile, mcName, clusterID, mcOIDCUrl, mcOIDC, mcOIDCArn string) (string, error) {

	// Step 1: Define role name
	roleName := fmt.Sprintf("rosa-hcp-bkp-%s-%s", mcName, clusterID)
	fmt.Printf("Role Name: %s\n", roleName)

	// Step 2: Create IAM Role
	fmt.Printf("Step 4: Creating IAM role '%s'...\n", roleName)

	// Construct the assume role policy document
	assumeRolePolicyDoc := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Federated":"%s"},"Action":["sts:AssumeRoleWithWebIdentity"],"Condition":{"StringEquals":{"%s:sub":"system:serviceaccount:openshift-adp:velero"}}}]}`, strings.TrimRightFunc(mcOIDCArn, unicode.IsSpace), mcOIDC)

	createRoleArgs := []string{
		"iam", "create-role",
		"--role-name", roleName,
		"--assume-role-policy-document", assumeRolePolicyDoc,
		"--description", fmt.Sprintf("\"backup-role for cluster %s\"", clusterID),
	}
	createRoleStdout, createRoleStderr, err := runCommand("aws", createRoleArgs...)
	if err != nil {
		// Check if the error is due to the role already existing
		if strings.Contains(createRoleStderr, "EntityAlreadyExists") {
			fmt.Printf("Warning: IAM role '%s' already exists. Skipping creation.\n", roleName)
		} else {
			return "", fmt.Errorf("failed to create IAM role: %w, stderr: %s", err, createRoleStderr)
		}
	} else {
		fmt.Println("IAM Role Creation Output:\n", createRoleStdout)
		fmt.Printf("IAM role '%s' created successfully.\n", roleName)
	}

	// Step 5: Get Role ARN
	fmt.Printf("Step 5: Getting ARN for role '%s'...\n", roleName)
	getRoleCmd := fmt.Sprintf(`aws iam get-role --role-name %s | awk 'NR==1 {print $2}'`, roleName)
	roleArnStdout, roleArnStderr, err := runCommand("bash", "-c", getRoleCmd)
	if err != nil {
		return "", fmt.Errorf("failed to get role ARN: %w, stderr: %s", err, roleArnStderr)
	}
	roleArn := strings.TrimSpace(roleArnStdout)
	if roleArn == "" {
		return "", fmt.Errorf("role ARN could not be retrieved for role '%s'", roleName)
	}
	fmt.Printf("role_arn: %s\n", roleArn)

	// Step 6: Attach AmazonS3FullAccess policy
	fmt.Printf("Step 6: Attaching AmazonS3FullAccess policy to role '%s'...\n", roleName)
	attachPolicyStdout, attachPolicyStderr, err := runCommand(
		"aws", "iam", "attach-role-policy",
		"--policy-arn", "arn:aws:iam::aws:policy/AmazonS3FullAccess",
		"--role-name", roleName)
	if err != nil {
		return "", fmt.Errorf("failed to attach AmazonS3FullAccess policy: %w, stderr: %s", err, attachPolicyStderr)
	}
	fmt.Println("Attach Policy Output:\n", attachPolicyStdout)
	fmt.Printf("AmazonS3FullAccess policy attached to role '%s'.\n", roleName)

	// Step 7: List attached role policies for verification
	fmt.Printf("Step 7: Listing attached policies for role '%s'...\n", roleName)
	listPoliciesStdout, listPoliciesStderr, err := runCommand(
		"aws", "iam", "list-attached-role-policies",
		"--role-name", roleName)
	if err != nil {
		return "", fmt.Errorf("failed to list attached role policies: %w, stderr: %s", err, listPoliciesStderr)
	}
	fmt.Println("List Policies Output:\n", listPoliciesStdout)

	fmt.Println("--- IAM Role Creation Completed ---")
	return roleArn, nil
}

// createKMSKeyAndPolicy creates an AWS KMS key, an associated IAM policy,
// attaches a key policy to the KMS key, and attaches the IAM policy to the role.
func createKMSKeyAndPolicy(awsProfile, clusterID, clusterEnv, awsRegion, roleArn string) (string, string, error) {
	fmt.Println("\n--- KMS Key and Policy Creation Started ---")

	// Step 1: Create KMS Key
	fmt.Printf("Step 1: Creating KMS key for cluster '%s'...\n", clusterID)
	createKeyCmd := fmt.Sprintf(`aws kms create-key --description "SSE-KMS backup key: %s" --key-usage ENCRYPT_DECRYPT --key-spec SYMMETRIC_DEFAULT --tags "TagKey=Owner,TagValue=%s" "TagKey=cluster,TagValue=%s" --region %s --query 'KeyMetadata.Arn' --output text`, clusterID, clusterEnv, clusterID, awsRegion)
	kmsArnStdout, kmsArnStderr, err := runCommand("bash", "-c", createKeyCmd)
	if err != nil {
		return "", "", fmt.Errorf("failed to create KMS key: %w, stderr: %s", err, kmsArnStderr)
	}
	kmsArn := strings.TrimSpace(kmsArnStdout)
	if kmsArn == "" {
		return "", "", fmt.Errorf("KMS ARN could not be retrieved")
	}
	fmt.Printf("kms_arn: %s\n", kmsArn)

	// Step 2: Define KMS IAM Policy Name
	kmsIAMPolicyName := fmt.Sprintf("AllowSSEKMSBackupKey-%s", clusterID)
	fmt.Printf("kms_iam_policy_name: %s\n", kmsIAMPolicyName)

	// Step 3: Create IAM Policy for KMS Key
	fmt.Printf("Step 3: Creating IAM policy '%s'...\n", kmsIAMPolicyName)
	kmsIAMPolicyDoc := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"kms:Encrypt",
					"kms:Decrypt",
					"kms:GenerateDataKey",
					"kms:DescribeKey"
				],
				"Resource": "%s"
			}
		]
	}`, kmsArn)

	createPolicyArgs := []string{
		"iam", "create-policy",
		"--policy-name", kmsIAMPolicyName,
		"--policy-document", kmsIAMPolicyDoc,
	}
	createPolicyStdout, createPolicyStderr, err := runCommand("aws", createPolicyArgs...)
	if err != nil {
		if strings.Contains(createPolicyStderr, "EntityAlreadyExists") {
			fmt.Printf("Warning: IAM policy '%s' already exists. Skipping creation.\n", kmsIAMPolicyName)
		} else {
			return "", "", fmt.Errorf("failed to create IAM policy: %w, stderr: %s", err, createPolicyStderr)
		}
	} else {
		fmt.Println("Create IAM Policy Output:\n", createPolicyStdout)
		fmt.Printf("IAM policy '%s' created successfully.\n", kmsIAMPolicyName)
	}

	// Step 4: Put Key Policy on KMS Key
	fmt.Printf("Step 4: Putting key policy on KMS key '%s'...\n", kmsArn)
	// Note: You need to replace '765374464689' with your actual AWS account ID for the 'emathias' user.
	// This is hardcoded in the original prompt, but ideally should be dynamic or a configurable parameter.
	// For this example, I'm using the hardcoded value.
	kmsKeyPolicyDoc := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Sid": "AllowClusterRoleAccess",
				"Effect": "Allow",
				"Principal": {
					"AWS": "%s"
				},
				"Action": [
					"kms:Encrypt",
					"kms:Decrypt",
					"kms:GenerateDataKey",
					"kms:DescribeKey"
				],
				"Resource": "*"
			},
			{
				"Sid": "AllowEmathiasFullAccess",
				"Effect": "Allow",
				"Principal": {
					"AWS": "*"
				},
				"Action": "kms:*",
				"Resource": "*",
				"Condition": {
					"StringEquals": {
						"aws:PrincipalArn": "arn:aws:iam::765374464689:user/emathias"
					}
				}
			}
		]
	}`, roleArn)

	putKeyPolicyArgs := []string{
		"kms", "put-key-policy",
		"--key-id", kmsArn,
		"--policy-name", "default",
		"--region", awsRegion,
		"--policy", kmsKeyPolicyDoc,
	}
	putKeyPolicyStdout, putKeyPolicyStderr, err := runCommand("aws", putKeyPolicyArgs...)
	if err != nil {
		return "", "", fmt.Errorf("failed to put key policy on KMS key: %w, stderr: %s", err, putKeyPolicyStderr)
	}
	fmt.Println("Put Key Policy Output:\n", putKeyPolicyStdout)
	fmt.Printf("Key policy attached to KMS key '%s'.\n", kmsArn)

	// Step 5: Get Policy ARN for the new IAM policy
	fmt.Printf("Step 5: Getting ARN for IAM policy '%s'...\n", kmsIAMPolicyName)
	getPolicyArnCmd := fmt.Sprintf(`aws iam list-policies --query "Policies[?PolicyName=='%s'].Arn" --output text`, kmsIAMPolicyName)
	policyArnStdout, policyArnStderr, err := runCommand("bash", "-c", getPolicyArnCmd)
	if err != nil {
		return "", "", fmt.Errorf("failed to get policy ARN for '%s': %w, stderr: %s", kmsIAMPolicyName, err, policyArnStderr)
	}
	policyArn := strings.TrimSpace(policyArnStdout)
	if policyArn == "" {
		return "", "", fmt.Errorf("policy ARN for '%s' could not be retrieved", kmsIAMPolicyName)
	}
	fmt.Printf("policy_arn: %s\n", policyArn)

	// Step 6: Attach the new IAM policy to the role
	fmt.Printf("Step 6: Attaching IAM policy '%s' to role '%s'...\n", kmsIAMPolicyName, strings.Split(roleArn, "/")[1]) // Extract role name from ARN
	attachRolePolicyStdout, attachRolePolicyStderr, err := runCommand(
		"aws", "iam", "attach-role-policy",
		"--role-name", strings.Split(roleArn, "/")[1], // Extract role name from ARN
		"--policy-arn", policyArn)
	if err != nil {
		return "", "", fmt.Errorf("failed to attach IAM policy '%s' to role '%s': %w, stderr: %s", kmsIAMPolicyName, strings.Split(roleArn, "/")[1], err, attachRolePolicyStderr)
	}
	fmt.Println("Attach Role Policy Output:\n", attachRolePolicyStdout)
	fmt.Printf("IAM policy '%s' attached to role '%s'.\n", kmsIAMPolicyName, strings.Split(roleArn, "/")[1])

	// Step 7: List attached role policies for verification
	fmt.Printf("Step 7: Listing attached policies for role '%s'...\n", strings.Split(roleArn, "/")[1])
	listPoliciesStdout, listPoliciesStderr, err := runCommand(
		"aws", "iam", "list-attached-role-policies",
		"--role-name", strings.Split(roleArn, "/")[1])
	if err != nil {
		return "", "", fmt.Errorf("failed to list attached role policies: %w, stderr: %s", err, listPoliciesStderr)
	}
	fmt.Println("List Policies Output:\n", listPoliciesStdout)

	fmt.Println("--- KMS Key and Policy Creation Completed ---")
	return kmsArn, kmsIAMPolicyName, nil
}

func main() {
	// --- IMPORTANT: Replace these placeholder values with your actual cluster details ---
	// You can get these from your ROSA cluster creation process.
	clusterID := os.Args[1]   // e.g., "abc123def456"
	clusterName := os.Args[2] // e.g., "my-rosa-cluster"
	clusterEnv := os.Args[3]  // e.g., "local", "int", "john.doe"
	mcName := os.Args[4]      // e.g., "hs-mc-n1j3kghkg", hive's mangement cluster name.
	awsProfile := os.Args[5]  // e.g., "dr-account" Your local aws config should have this name.
	awsRegion := os.Args[6]   // e.g., us-west-2 can be whichever region you are looking for.
	// ----------------------------------------------------------------------------------

	// Call the setupCluster function with your cluster details
	err := setupCluster(clusterID, clusterName, clusterEnv)
	if err != nil {
		log.Fatalf("Cluster setup failed: %v", err)
	}
	// Call the createS3Bucket function with your AWS details
	err = createS3Bucket(awsProfile, awsRegion)
	if err != nil {
		log.Printf("AWS S3 bucket creation failed: %v", err)
		// Do not exit fatally here, allow other independent operations to proceed.
	}

	// Call the createOIDCConfig function with your management cluster details
	mcOIDCUrl, mcOIDC, mcOIDCArn, err := createOIDCConfig(mcName, awsRegion, clusterID)
	if err != nil {
		log.Fatalf("OIDC configuration failed: %v", err)
	}
	fmt.Printf("\nFinal OIDC URL: %s\nFinal OIDC ID: %s\nFinal OIDC Arn: %s\n", mcOIDCUrl, mcOIDC, mcOIDCArn)

	// Call the createIAMRole function with your AWS and OIDC details
	roleArn, err := createIAMRole(awsProfile, mcName, clusterID, mcOIDCUrl, mcOIDC, mcOIDCArn)
	if err != nil {
		log.Fatalf("IAM role creation and policy attachment failed: %v", err)
	}
	fmt.Printf("\nFinal IAM Role ARN: %s\n", roleArn)

	// Call the createKMSKeyAndPolicy function with your AWS, cluster, and role details
	kmsArn, kmsIAMPolicyName, err := createKMSKeyAndPolicy(awsProfile, clusterID, clusterEnv, awsRegion, roleArn)
	if err != nil {
		log.Fatalf("KMS key and policy creation failed: %v", err)
	}
	fmt.Printf("\nFinal KMS ARN: %s\nFinal KMS IAM Policy Name: %s\n", kmsArn, kmsIAMPolicyName)
}
