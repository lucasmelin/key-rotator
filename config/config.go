package config

import (
	"context"
	"fmt"
	"os"

	"github.com/lucasmelin/key-rotator/github"
	"gopkg.in/yaml.v3"
)

// KeyConfig represents the structure of the YAML configuration file.
type KeyConfig struct {
	Secrets []Secret `yaml:"secrets"`
}

// Secret represents a secret and its destinations.
type Secret struct {
	Name         string               `yaml:"name"`
	Description  string               `yaml:"description"`
	Destinations []DestinationWrapper `yaml:"destinations"`
}

// Destination represents a destination where the secret should be stored.
type Destination interface {
	UpdateSecret(ctx context.Context, client github.Client, secretValue string) error
	GetDescription() string
}

// DestinationWrapper wraps the Destination interface for custom unmarshaling.
type DestinationWrapper struct {
	Destination
}

// UnmarshalYAML custom unmarshaler for Destination.
func (d *DestinationWrapper) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]interface{}
	if err := value.Decode(&raw); err != nil {
		return err
	}

	destType, ok := raw["type"].(string)
	if !ok {
		return fmt.Errorf("destination type is required")
	}
	switch destType {
	case github.TypeGitHubRepository:
		var dest github.RepositorySecret
		if err := value.Decode(&dest); err != nil {
			return err
		}
		d.Destination = dest
	case github.TypeGitHubRepositoryDependabot:
		var dest github.DependabotRepositorySecret
		if err := value.Decode(&dest); err != nil {
			return err
		}
		d.Destination = dest
	case github.TypeGitHubRepositoryEnvironment:
		var dest github.RepositoryEnvironmentSecret
		if err := value.Decode(&dest); err != nil {
			return err
		}
		d.Destination = dest
	default:
		return fmt.Errorf("unsupported destination type: %s", destType)
	}

	return nil
}

// ParseFile reads and parses the YAML configuration file.
func ParseFile(yamlFile string) (KeyConfig, error) {
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
