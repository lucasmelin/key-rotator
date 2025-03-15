package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/google/go-github/v69/github"
	"github.com/jamesruan/sodium"
	"gopkg.in/yaml.v3"
)

// KeyConfig represents the structure of the YAML configuration file.
type KeyConfig struct {
	Secrets []Secret `yaml:"secrets"`
}

// Secret represents a secret and its destinations.
type Secret struct {
	Name         string        `yaml:"name"`
	Description  string        `yaml:"description"`
	Destinations []Destination `yaml:"destinations"`
}

// Destination represents a destination where the secret should be stored.
type Destination struct {
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Repo        string `yaml:"repo"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
}

// GitHubSecret represents an encrypted GitHub secret.
type GitHubSecret struct {
	Name           string
	KeyID          string
	EncryptedValue string
}

// GitHubClient wraps the GitHub client.
type GitHubClient struct {
	*github.Client
}

// Secret destination types
const (
	typeGitHubRepository            = "github-repository"
	typeGitHubRepositoryDependabot  = "github-repository-dependabot"
	typeGitHubRepositoryEnvironment = "github-repository-environment"
)

func main() {
	// Check if the YAML file path is provided as an argument.
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path to YAML config file>", os.Args[0])
	}
	yamlFile := os.Args[1]

	// Parse the YAML configuration file.
	config, err := parseFile(yamlFile)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	// Create a new GitHub client.
	client := NewGitHubClient()
	ctx := context.Background()

	// Iterate over each secret in the configuration.
	for _, secret := range config.Secrets {
		// Prompt the user to enter the value for the secret.
		secretValue, err := secretPrompt(secret.Description)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		// Iterate over each destination for the secret.
		for _, dest := range secret.Destinations {
			ownerRepo := strings.Split(dest.Repo, "/")
			if len(ownerRepo) != 2 {
				log.Fatalf("Invalid destination format: %s", dest.Repo)
			}
			owner, repo := ownerRepo[0], ownerRepo[1]

			// Get the public key for the repository.
			var key *github.PublicKey
			switch dest.Type {
			case typeGitHubRepository:
				key, _, err = client.Actions.GetRepoPublicKey(ctx, owner, repo)
				if err != nil {
					log.Fatalf("Failed to get public key: %v", err)
				}
			case typeGitHubRepositoryDependabot:
				key, _, err = client.Dependabot.GetRepoPublicKey(ctx, owner, repo)
				if err != nil {
					log.Fatalf("Failed to get public key: %v", err)
				}
			case typeGitHubRepositoryEnvironment:
				repository, _, err := client.Repositories.Get(ctx, owner, repo)
				if err != nil {
					log.Fatalf("Failed to get repository: %v", err)
				}
				key, _, err = client.Actions.GetEnvPublicKey(ctx, int(repository.GetID()), dest.Environment)
				if err != nil {
					log.Fatalf("Failed to get public key: %v", err)
				}
			default:
				log.Fatalf("Unsupported destination type: %s", dest.Type)
			}

			// Encrypt the secret value using the public key.
			encryptedValue, err := encryptSodiumSecret(secretValue, key.GetKey())
			if err != nil {
				log.Fatalf("Failed to encrypt secret: %v", err)
			}

			// Create a GitHubSecret object.
			ghSecret := GitHubSecret{
				Name:           dest.Name,
				KeyID:          key.GetKeyID(),
				EncryptedValue: encryptedValue,
			}

			// Update the secret in the repository.
			switch dest.Type {
			case typeGitHubRepository:
				err = client.updateRepositorySecret(ctx, owner, repo, ghSecret)
				if err != nil {
					log.Fatalf("Failed to update secret: %v", err)
				}
			case typeGitHubRepositoryDependabot:
				err = client.updateDependabotSecret(ctx, owner, repo, ghSecret)
				if err != nil {
					log.Fatalf("Failed to update secret: %v", err)
				}
			case typeGitHubRepositoryEnvironment:
				err = client.updateEnvironmentSecret(ctx, owner, repo, dest.Environment, ghSecret)
				if err != nil {
					log.Fatalf("Failed to update secret: %v", err)
				}
			default:
				log.Fatalf("Unsupported destination type: %s", dest.Type)
			}

			fmt.Println("Updated", dest.Description)
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

// NewGitHubClient creates a new GitHub client with authentication.
func NewGitHubClient() GitHubClient {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatalf("The GITHUB_TOKEN environment variable must be set")
	}
	return GitHubClient{github.NewClient(nil).WithAuthToken(token)}
}

// parseFile reads and parses the YAML configuration file.
func parseFile(yamlFile string) (KeyConfig, error) {
	file, err := os.Open(yamlFile)
	if err != nil {
		return KeyConfig{}, fmt.Errorf("failed to open %s: %v", yamlFile, err)
	}
	defer file.Close()

	var config KeyConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return KeyConfig{}, fmt.Errorf("failed to decode %s: %v", yamlFile, err)
	}
	return config, nil
}

// encryptSodiumSecret encrypts a secret value using a public key.
// It uses the sodium library for encryption and returns the base64-encoded encrypted value.
func encryptSodiumSecret(secretValue string, publicKey string) (string, error) {
	secretBytes := sodium.Bytes(secretValue)
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %w", err)
	}
	sodiumKey := sodium.BoxPublicKey{Bytes: publicKeyBytes}
	encryptedSecret := secretBytes.SealedBox(sodiumKey)
	return base64.StdEncoding.EncodeToString(encryptedSecret), nil
}

// updateRepositorySecret updates a GitHub Actions secret in the repository.
func (ghc GitHubClient) updateRepositorySecret(ctx context.Context, owner string, repo string, secret GitHubSecret) error {
	s := &github.EncryptedSecret{
		Name:           secret.Name,
		KeyID:          secret.KeyID,
		EncryptedValue: secret.EncryptedValue,
	}
	_, err := ghc.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, s)
	return err
}

// updateDependabotSecret updates a GitHub Dependabot secret in the repository.
func (ghc GitHubClient) updateDependabotSecret(ctx context.Context, owner string, repo string, secret GitHubSecret) error {
	s := &github.DependabotEncryptedSecret{
		Name:           secret.Name,
		KeyID:          secret.KeyID,
		EncryptedValue: secret.EncryptedValue,
	}
	_, err := ghc.Dependabot.CreateOrUpdateRepoSecret(ctx, owner, repo, s)
	return err
}

// updateEnvironmentSecret updates a GitHub environment secret in the repository.
func (ghc GitHubClient) updateEnvironmentSecret(ctx context.Context, owner string, repo string, environment string, secret GitHubSecret) error {
	repository, _, err := ghc.Repositories.Get(ctx, owner, repo)
	if err != nil {
		log.Fatalf("Failed to get repository: %v", err)
	}

	s := &github.EncryptedSecret{
		Name:           secret.Name,
		KeyID:          secret.KeyID,
		EncryptedValue: secret.EncryptedValue,
	}
	_, err = ghc.Actions.CreateOrUpdateEnvSecret(ctx, int(repository.GetID()), environment, s)
	return err
}
