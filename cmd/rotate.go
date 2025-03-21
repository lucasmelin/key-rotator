package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/lucasmelin/key-rotator/config"
	"github.com/lucasmelin/key-rotator/github"
	"github.com/spf13/cobra"
)

var dryRun bool

type rotateOptions struct {
	dryRun   bool
	yamlFile string
}

var rotateCmd = &cobra.Command{
	Use:     "rotate <path to YAML config file>",
	Short:   "Rotate secrets based on the provided configuration file",
	Args:    cobra.ExactArgs(1),
	GroupID: "core-commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun {
			fmt.Println("Running in dry-run mode, no changes will be made")
		}

		opts := &rotateOptions{
			dryRun:   dryRun,
			yamlFile: args[0],
		}

		return runRotate(opts)
	},
}

func runRotate(opts *rotateOptions) error {
	// Parse the YAML configuration file.
	cfg, err := config.ParseFile(opts.yamlFile)
	if err != nil {
		return fmt.Errorf("failed to parse file: %v", err)
	}

	// Create a new GitHub client.
	client := github.NewClient()
	ctx := context.Background()

	// Iterate over each secret in the configuration.
	for _, secret := range cfg.Secrets {
		// Prompt the user to enter the value for the secret.
		secretValue, err := secretPrompt(fmt.Sprintf("%s: %s", secret.Name, secret.Description))
		if err != nil {
			return fmt.Errorf("failed to read input: %v", err)
		}

		// Display the destinations that will be updated.
		fmt.Println("The following destinations will be updated:")
		for _, d := range secret.Destinations {
			fmt.Println("-", d.Destination.GetDescription())
		}

		// Prompt the user to accept before updating.
		confirm, err := confirmPrompt("Do you want to proceed with updating these destinations? (y/N)")
		if err != nil {
			return fmt.Errorf("failed to read input: %v", err)
		}
		if !confirm {
			fmt.Println("Operation cancelled by the user.")
			continue
		}

		// Iterate over each destination for the secret.
		for _, d := range secret.Destinations {
			if opts.dryRun {
				fmt.Printf("[Dry Run] Would update %s with provided secret value for %s\n", d.Destination.GetDescription(), secret.Name)
			} else {
				if err = d.Destination.UpdateSecret(ctx, client, secretValue); err != nil {
					return fmt.Errorf("failed to update secret: %v", err)
				}
				fmt.Println("Updated", d.Destination.GetDescription())
			}
		}
	}
	return nil
}

func init() {
	rotateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print out the changes that would be made without actually making them")
}

func secretPrompt(title string) (string, error) {
	var secretValue string
	err := huh.NewInput().
		Title(title).
		Value(&secretValue).
		EchoMode(huh.EchoModePassword).
		Run()
	return secretValue, err
}

func confirmPrompt(title string) (bool, error) {
	var response string
	fmt.Print(title + " ")
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}
	return response == "y" || response == "Y", nil
}
