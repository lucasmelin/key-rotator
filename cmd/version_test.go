package cmd

import "testing"

func Test_formatVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		buildDate string
		want      string
	}{
		{
			name:      "valid semver version with leading v",
			version:   "v0.1.0",
			buildDate: "2025-03-20",
			want:      "key-rotator version 0.1.0 (2025-03-20)\nhttps://github.com/lucasmelin/key-rotator/releases/tag/v0.1.0\n",
		},
		{
			name:      "valid semver version without leading v",
			version:   "0.1.0",
			buildDate: "2025-03-20",
			want:      "key-rotator version 0.1.0 (2025-03-20)\nhttps://github.com/lucasmelin/key-rotator/releases/tag/v0.1.0\n",
		},
		{
			name:      "valid semver version with beta tag",
			version:   "v0.1.0-beta",
			buildDate: "2025-03-20",
			want:      "key-rotator version 0.1.0-beta (2025-03-20)\nhttps://github.com/lucasmelin/key-rotator/releases/tag/v0.1.0-beta\n",
		},
		{
			name:      "invalid semver version",
			version:   "invalid-semver-version",
			buildDate: "2025-03-20",
			want:      "key-rotator version invalid-semver-version (2025-03-20)\nhttps://github.com/lucasmelin/key-rotator/releases/latest\n",
		},
		{
			name:      "no build date",
			version:   "v0.1.0",
			buildDate: "",
			want:      "key-rotator version 0.1.0\nhttps://github.com/lucasmelin/key-rotator/releases/tag/v0.1.0\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatVersion(tt.version, tt.buildDate)
			if got != tt.want {
				t.Errorf("formatVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_changelogURL(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "valid semver version with leading v",
			version: "v0.1.0",
			want:    "https://github.com/lucasmelin/key-rotator/releases/tag/v0.1.0",
		},
		{
			name:    "valid semver version without leading v",
			version: "0.1.0",
			want:    "https://github.com/lucasmelin/key-rotator/releases/tag/v0.1.0",
		},
		{
			name:    "valid semver version with beta tag",
			version: "v0.1.0-beta",
			want:    "https://github.com/lucasmelin/key-rotator/releases/tag/v0.1.0-beta",
		},
		{
			name:    "invalid semver version",
			version: "invalid-semver-version",
			want:    "https://github.com/lucasmelin/key-rotator/releases/latest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := changelogURL(tt.version)
			if got != tt.want {
				t.Errorf("changelogURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
