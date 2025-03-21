package cmd

import (
	"fmt"
	"regexp"
	"strings"

	build "github.com/lucasmelin/key-rotator/internal"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(formatVersion(build.Version, build.Date))
	},
	Hidden: true,
}

func formatVersion(version string, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	var date string
	if buildDate != "" {
		date = fmt.Sprintf(" (%s)", buildDate)
	}
	return fmt.Sprintf("key-rotator version %s%s\n%s\n", version, date, changelogURL(version))
}

func changelogURL(version string) string {
	path := "https://github.com/lucasmelin/key-rotator"
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	// If the version is not a valid semver, return the latest release URL.
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", path)
	}
	return fmt.Sprintf("%s/releases/tag/v%s", path, strings.TrimPrefix(version, "v"))
}
