package cmd

import (
	"errors"

	"github.com/golgoth31/multiShellKonfig/internal/config"
)

var (
	debug           bool
	noID            bool
	konfGoReqID     string
	cfgFile         string
	cfgDir          string
	cfgContextsPath string
	cfgData         config.Konfigs
	homedir         string
	errNoReqID      = errors.New("request ID not set")
	cleanAll        bool
)
