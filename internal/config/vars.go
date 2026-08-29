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

	// Dynamic version retrieve with ldflags.

	// Version represent version of application.
	Version string

	// Date represent date of build.
	Date string

	// BuiltBy represent date of build.
	BuiltBy string

	// BuiltBy represent date of build.
	GitCommit string
)
