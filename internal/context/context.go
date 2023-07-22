package context

import (
	"fmt"
	"sort"

	"github.com/golgoth31/multiShellKonfig/internal/config"
	"github.com/golgoth31/multiShellKonfig/pkg/konfig"
	"github.com/golgoth31/multiShellKonfig/pkg/shell"
	"github.com/rs/zerolog/log"
)

func New(mskConfig []config.Konfig, homedir string) (*Context, error) {
	curContexts := &Context{}

	// get all available contexts
	for _, unitKonfig := range mskConfig {
		log.Debug().Msgf("found config: %s", unitKonfig.Path)

		kubeConfig, err := konfig.Load(unitKonfig.Path, homedir)
		if err != nil {
			return &Context{}, err
		}

		curKonfig := konfig.Konfig{
			FileID:   unitKonfig.ID,
			FilePath: unitKonfig.Path,
			Content:  kubeConfig,
		}

		curContexts.KonfigList = append(curContexts.KonfigList, &curKonfig)
	}

	return curContexts, nil
}

func (ctx *Context) GetContextList() ([]string, error) {
	contextList := shell.ShellContextList{}
	contextListString := []string{}

	for _, unitKonfig := range ctx.KonfigList {
		for _, context := range unitKonfig.Content.Contexts {
			log.Debug().Msgf("found context '%s@%s'", context.Name, unitKonfig.FilePath)

			contextList = append(
				contextList,
				shell.ContextDef{
					Name:     context.Name,
					FileID:   unitKonfig.FileID,
					FilePath: unitKonfig.FilePath,
				},
			)
		}
	}

	// Sort contextList by context name
	sort.Stable(contextList)

	log.Debug().Msgf("context list: %v", contextList)

	// Generate list of context for select
	for _, v := range contextList {
		contextListString = append(
			contextListString,
			fmt.Sprintf(
				"%s@%s",
				v.Name,
				v.FilePath,
			),
		)
	}

	return contextListString, nil
}
