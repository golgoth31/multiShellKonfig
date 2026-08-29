// Package config ...
package config

import "uuid"

var (
	// DefaultConfig defines the base default configuration.
	DefaultConfig = Konfigs{
		Konfigs: []Konfig{
			{
				Path: "~/.kube/config",
				ID:   uuid.New().String(),
			},
		},
	}

	// The following version variables are injected at build time via the
	// -ldflags -X flag (see the Makefile).

	// Version represents the application version.
	Version string

	// Date represents the build date.
	Date string

	// BuiltBy identifies the user that built the binary.
	BuiltBy string

	// GitCommit holds the git commit the binary was built from.
	GitCommit string
)
