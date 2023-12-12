package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"trackCommit/config"
	"trackCommit/db"
	"trackCommit/model"

	"github.com/spf13/cobra"
)

var allowedTypes = []string{"bug", "feature", "changes", "conflict"}

func isValidType(value string) bool {
	for _, t := range allowedTypes {
		if value == t {
			return true
		}
	}
	return false
}

func isValidCommitHash(hash string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", hash)
	err := cmd.Run()
	return err == nil
}

func GetAddCommand() *cobra.Command {
	var (
		typeFlag string
		hashFlag string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add the commit to the target branch",
		Run: func(cmd *cobra.Command, args []string) {
			if !isValidType(typeFlag) {
				fmt.Printf("Error: Invalid type '%s'. Allowed values are %v\n", typeFlag, allowedTypes)
				os.Exit(1)
			}

			if hashFlag == "" {
				fmt.Println("Error: --hash flag is mandatory and cannot be an empty string.")
				os.Exit(1)
			}

			if !isValidCommitHash(hashFlag) {
				fmt.Printf("Error: Invalid commit hash '%s'.\n", hashFlag)
				os.Exit(1)
			}

			commitDetails, err := getCommitDetails(hashFlag)
			if err != nil {
				fmt.Printf("Error getting commit details: %v\n", err)
				os.Exit(1)
			}

			commitDetails.Branch = "bo_revamp_merged"

			// Check if the commit hash already exists
			existingDetails, err := db.RetrieveCommitDetailsByHash(&config.DbConnection, "commitDetails", hashFlag)
			if err == nil && existingDetails.Hash != "" {
				fmt.Printf("Error: Commit with hash '%s' already exists.\n", hashFlag)
				os.Exit(1)
			}

			err = db.SaveCommitDetails(&config.DbConnection, "commitDetails", commitDetails)
			if err != nil {
				fmt.Printf("Error saving commit details: %v\n", err)
				os.Exit(1)
			}

			// Retrieve commit details
			retrievedDetails, err := db.RetrieveCommitDetailsByHash(&config.DbConnection, "commitDetails", commitDetails.Hash)
			if err != nil {
				fmt.Printf("Error retrieving commit details: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("saved successful")
			fmt.Println("retrievedDetails :", retrievedDetails)
		},
	}

	cmd.Flags().StringVarP(&typeFlag, "type", "t", "", "Type of the commit (bug, feature, changes, conflict)")
	cmd.MarkFlagRequired("type")

	cmd.Flags().StringVarP(&hashFlag, "hash", "H", "", "Hash of the commit (mandatory)")
	cmd.MarkFlagRequired("hash")

	return cmd
}

func getCommitDetails(hash string) (model.CommitDetails, error) {
	cmd := exec.Command("git", "show", "--pretty=format:%H%n%an%n%ad%n%s%n", hash)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("inside one")
		return model.CommitDetails{}, fmt.Errorf("failed to get commit details: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 4 {
		log.Println("inside two")
		return model.CommitDetails{}, fmt.Errorf("unexpected output format")
	}

	return model.CommitDetails{
		Hash:    lines[0],
		Author:  lines[1],
		Date:    lines[2],
		Message: lines[3],
	}, nil
}
