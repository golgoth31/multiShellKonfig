package context

import "github.com/golgoth31/multiShellKonfig/pkg/konfig"

type Context struct {
	MskReqID   string
	KonfigList []*konfig.Konfig
}
