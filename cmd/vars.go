package cmd

import (
	"errors"

	"github.com/golgoth31/multiShellKonfig/internal/config"
)

var (
	debug           bool
	noID            bool
	cfgFile         string
	cfgDir          string
	cfgContextsPath string
	cfgData         config.Konfigs
	homedir         string
	errNoReqID      = errors.New("request ID not set")
	cleanAll        bool
)

const (
	filePerm = 0600
	dirPerm  = 0700
)
