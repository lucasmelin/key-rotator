package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/lucasmelin/key-rotator/config"
	"github.com/lucasmelin/key-rotator/github"
)

func main() {
	// Check if the YAML file path is provided as an argument.
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path to YAML config file>", os.Args[0])
	}
	yamlFile := os.Args[1]

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
		secretValue, err := secretPrompt(secret.Description)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		// Iterate over each destination for the secret.
		for _, d := range secret.Destinations {
			if err = d.Destination.UpdateSecret(ctx, client, secretValue); err != nil {
				log.Fatalf("Failed to update secret: %v", err)
			}

			fmt.Println("Updated", d.Destination.GetDescription())
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
