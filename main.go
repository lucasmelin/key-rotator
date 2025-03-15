package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/lucasmelin/key-rotator/config"
	"github.com/lucasmelin/key-rotator/github"
)

func main() {
	// Define the dry-run flag.
	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", false, "Print out the API calls without making them")
	flag.Parse()

	if dryRun {
		fmt.Println("Running in dry-run mode, no changes will be made")
	}

	// Check if the YAML file path is provided as an argument.
	if len(flag.Args()) < 1 {
		log.Fatalf("Usage: %s [--dry-run] <path to YAML config file>", os.Args[0])
	}
	yamlFile := flag.Args()[0]

	// Parse the YAML configuration file.
	cfg, err := config.ParseFile(yamlFile)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	// Create a new GitHub client.
	client := github.NewClient()
	ctx := context.Background()

	// Iterate over each secret in the configuration.
	for _, secret := range cfg.Secrets {
		// Prompt the user to enter the value for the secret.
		secretValue, err := secretPrompt(fmt.Sprintf("%s: %s", secret.Name, secret.Description))
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		// Display the destinations that will be updated.
		fmt.Println("The following destinations will be updated:")
		for _, d := range secret.Destinations {
			fmt.Println("-", d.Destination.GetDescription())
		}

		// Prompt the user to accept before updating.
		confirm, err := confirmPrompt("Do you want to proceed with updating these destinations? (y/N)")
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}
		if !confirm {
			fmt.Println("Operation cancelled by the user.")
			continue
		}

		// Iterate over each destination for the secret.
		for _, d := range secret.Destinations {
			if dryRun {
				fmt.Printf("[Dry Run] Would update %s with provided secret value for %s\n", d.Destination.GetDescription(), secret.Name)
			} else {
				if err = d.Destination.UpdateSecret(ctx, client, secretValue); err != nil {
					log.Fatalf("Failed to update secret: %v", err)
				}
				fmt.Println("Updated", d.Destination.GetDescription())
			}
		}
	}
}

// secretPrompt prompts the user to enter a secret value.
func secretPrompt(title string) (string, error) {
	var secretValue string
	err := huh.NewInput().
		Title(title).
		Value(&secretValue).
		EchoMode(huh.EchoModePassword).
		Run()
	return secretValue, err
}

// confirmPrompt prompts the user to confirm an action.
func confirmPrompt(title string) (bool, error) {
	var response string
	fmt.Print(title + " ")
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}
	return response == "y" || response == "Y", nil
}
