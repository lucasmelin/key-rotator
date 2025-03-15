package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/jamesruan/sodium"
)

// GitHub secret destination types.
const (
	TypeGitHubRepository            = "github-repository"
	TypeGitHubRepositoryDependabot  = "github-repository-dependabot"
	TypeGitHubRepositoryEnvironment = "github-repository-environment"
)

// Client wraps the GitHub client.
type Client struct {
	*github.Client
}

// NewClient creates a new GitHub client with authentication.
func NewClient() Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatalf("The GITHUB_TOKEN environment variable must be set")
	}
	return Client{github.NewClient(nil).WithAuthToken(token)}
}

// RepositorySecret represents a GitHub repository secret destination.
type RepositorySecret struct {
	Repo string `yaml:"repo"`
	Name string `yaml:"name"`
}

// UpdateSecret updates the GitHub Actions secret in the repository.
func (d RepositorySecret) UpdateSecret(ctx context.Context, client Client, secretValue string) error {
	ownerRepo := strings.Split(d.Repo, "/")
	if len(ownerRepo) != 2 {
		return fmt.Errorf("invalid destination format: %s", d.Repo)
	}
	owner, repo := ownerRepo[0], ownerRepo[1]

	key, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get public key: %v", err)
	}

	encryptedValue, err := encryptSodiumSecret(secretValue, key.GetKey())
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %v", err)
	}

	ghSecret := secret{
		Name:           d.Name,
		KeyID:          key.GetKeyID(),
		EncryptedValue: encryptedValue,
	}

	return client.updateRepositorySecret(ctx, owner, repo, ghSecret)
}

// GetDescription returns the destination description.
func (d RepositorySecret) GetDescription() string {
	return fmt.Sprintf("%s GitHub Repository Secret in the %s repository", d.Name, d.Repo)
}

// DependabotRepositorySecret represents a GitHub Dependabot secret destination.
type DependabotRepositorySecret struct {
	Repo string `yaml:"repo"`
	Name string `yaml:"name"`
}

// GetDescription returns the destination description.
func (d DependabotRepositorySecret) GetDescription() string {
	return fmt.Sprintf("%s GitHub Dependabot Repository Secret in the %s repository", d.Name, d.Repo)
}

// UpdateSecret updates the Dependabot secret in the repository.
func (d DependabotRepositorySecret) UpdateSecret(ctx context.Context, client Client, secretValue string) error {
	ownerRepo := strings.Split(d.Repo, "/")
	if len(ownerRepo) != 2 {
		return fmt.Errorf("invalid destination format: %s", d.Repo)
	}
	owner, repo := ownerRepo[0], ownerRepo[1]

	key, _, err := client.Dependabot.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get public key: %v", err)
	}

	encryptedValue, err := encryptSodiumSecret(secretValue, key.GetKey())
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %v", err)
	}

	ghSecret := secret{
		Name:           d.Name,
		KeyID:          key.GetKeyID(),
		EncryptedValue: encryptedValue,
	}

	return client.updateDependabotSecret(ctx, owner, repo, ghSecret)
}

// RepositoryEnvironmentSecret represents a GitHub environment secret destination.
type RepositoryEnvironmentSecret struct {
	Repo        string `yaml:"repo"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
}

// GetDescription returns the destination description.
func (d RepositoryEnvironmentSecret) GetDescription() string {
	return fmt.Sprintf("%s GitHub Repository Environment Secret in the %s repository's %s environment", d.Name, d.Repo, d.Environment)
}

// UpdateSecret updates the GitHub Actions environment secret in the repository.
func (d RepositoryEnvironmentSecret) UpdateSecret(ctx context.Context, client Client, secretValue string) error {
	ownerRepo := strings.Split(d.Repo, "/")
	if len(ownerRepo) != 2 {
		return fmt.Errorf("invalid destination format: %s", d.Repo)
	}
	owner, repo := ownerRepo[0], ownerRepo[1]

	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get repository: %v", err)
	}

	key, _, err := client.Actions.GetEnvPublicKey(ctx, int(repository.GetID()), d.Environment)
	if err != nil {
		return fmt.Errorf("failed to get public key: %v", err)
	}

	encryptedValue, err := encryptSodiumSecret(secretValue, key.GetKey())
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %v", err)
	}

	ghSecret := secret{
		Name:           d.Name,
		KeyID:          key.GetKeyID(),
		EncryptedValue: encryptedValue,
	}

	return client.updateEnvironmentSecret(ctx, owner, repo, d.Environment, ghSecret)
}

// updateRepositorySecret updates a GitHub Actions secret in the repository.
func (ghc Client) updateRepositorySecret(ctx context.Context, owner string, repo string, secret secret) error {
	s := &github.EncryptedSecret{
		Name:           secret.Name,
		KeyID:          secret.KeyID,
		EncryptedValue: secret.EncryptedValue,
	}
	_, err := ghc.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, s)
	return err
}

// updateDependabotSecret updates a GitHub Dependabot secret in the repository.
func (ghc Client) updateDependabotSecret(ctx context.Context, owner string, repo string, secret secret) error {
	s := &github.DependabotEncryptedSecret{
		Name:           secret.Name,
		KeyID:          secret.KeyID,
		EncryptedValue: secret.EncryptedValue,
	}
	_, err := ghc.Dependabot.CreateOrUpdateRepoSecret(ctx, owner, repo, s)
	return err
}

// updateEnvironmentSecret updates a GitHub environment secret in the repository.
func (ghc Client) updateEnvironmentSecret(ctx context.Context, owner string, repo string, environment string, secret secret) error {
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

// secret represents an encrypted GitHub secret.
type secret struct {
	Name           string
	KeyID          string
	EncryptedValue string
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
