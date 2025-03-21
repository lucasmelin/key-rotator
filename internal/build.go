package build

import (
	"runtime/debug"
)

// Version is dynamically set by the toolchain or overridden at build time.
var Version = "DEV"

// Date is dynamically set at build time.
// The format is YYYY-MM-DD.
var Date = ""

func init() {
	if Version == "DEV" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}
