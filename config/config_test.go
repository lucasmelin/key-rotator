package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lucasmelin/key-rotator/github"
	"gopkg.in/yaml.v3"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		expectError bool
		expected    KeyConfig
	}{
		{
			name: "GitHub repository secret",
			yamlContent: `
secrets:
  - name: test-secret
    description: A test secret
    destinations:
      - type: github-repository
        description: Test repo secret
        repo: owner/repo
        name: TEST_SECRET
`,
			expectError: false,
			expected: KeyConfig{
				Secrets: []Secret{
					{
						Name:        "test-secret",
						Description: "A test secret",
						Destinations: []DestinationWrapper{
							{
								Destination: github.RepositorySecret{
									Description: "Test repo secret",
									Repo:        "owner/repo",
									Name:        "TEST_SECRET",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "GitHub repository dependabot secret",
			yamlContent: `
secrets:
  - name: test-secret
    description: A test secret
    destinations:
      - type: github-repository-dependabot
        description: Test dependabot secret
        repo: owner/repo
        name: TEST_SECRET
`,
			expectError: false,
			expected: KeyConfig{
				Secrets: []Secret{
					{
						Name:        "test-secret",
						Description: "A test secret",
						Destinations: []DestinationWrapper{
							{
								Destination: github.DependabotRepositorySecret{
									Description: "Test dependabot secret",
									Repo:        "owner/repo",
									Name:        "TEST_SECRET",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "GitHub repository environment secret",
			yamlContent: `
secrets:
  - name: test-secret
    description: A test secret
    destinations:
      - type: github-repository-environment
        description: Test environment secret
        repo: owner/repo
        environment: prod
        name: TEST_SECRET
`,
			expectError: false,
			expected: KeyConfig{
				Secrets: []Secret{
					{
						Name:        "test-secret",
						Description: "A test secret",
						Destinations: []DestinationWrapper{
							{
								Destination: github.RepositoryEnvironmentSecret{
									Description: "Test environment secret",
									Repo:        "owner/repo",
									Environment: "prod",
									Name:        "TEST_SECRET",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Invalid secret type",
			yamlContent: `
secrets:
  - name: test-secret
    description: A test secret
    destinations:
      - type: invalid-type
        description: Invalid type secret
`,
			expectError: true,
		},
		{
			name: "Missing secret type",
			yamlContent: `
secrets:
  - name: test-secret
    description: A test secret
    destinations:
      - description: Missing type secret
`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tmpFile := filepath.Join(tempDir, "test-file.yaml")
			f, err := os.Create(tmpFile)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			if _, err := f.Write([]byte(tt.yamlContent)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			if err := f.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			config, err := ParseFile(tmpFile)
			if (err != nil) != tt.expectError {
				t.Fatalf("ParseFile error = %v, expectError %v", err, tt.expectError)
			}
			if !tt.expectError && !cmp.Equal(config, tt.expected) {
				t.Errorf("Expected config %+v, got %+v", tt.expected, config)
			}
		})
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := ParseFile("nonexistent.yaml")
	if err == nil {
		t.Fatal("Expected error for nonexistent file, got nil")
	}
}

func TestUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		expectError bool
		expected    interface{}
	}{
		{
			name: "Valid GitHub repository secret",
			yamlContent: `type: github-repository
description: Test repo secret
repo: owner/repo
name: TEST_SECRET`,
			expectError: false,
			expected: github.RepositorySecret{
				Description: "Test repo secret",
				Repo:        "owner/repo",
				Name:        "TEST_SECRET",
			},
		},
		{
			name:        "Unsupported type",
			yamlContent: "type: unsupported-type\ndescription: Unsupported type secret",
			expectError: true,
		},
		{
			name:        "Missing type",
			yamlContent: `description: Missing type secret`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wrapper DestinationWrapper
			err := yaml.Unmarshal([]byte(tt.yamlContent), &wrapper)
			if (err != nil) != tt.expectError {
				t.Fatalf("UnmarshalYAML error = %v, expectError %v", err, tt.expectError)
			}
			if !tt.expectError {
				if !cmp.Equal(wrapper.Destination, tt.expected) {
					t.Errorf("Expected destination %+v, got %+v", tt.expected, wrapper.Destination)
				}
			}
		})
	}
}
