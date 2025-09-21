package shell

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
)

func LoadList(itemType string, currentItem string, itemList []string) (string, error) {
	var (
		output string
	)

	items := []string{}
	contextStyle := pterm.NewStyle(pterm.FgLightBlue, pterm.Bold)
	fileStyle := pterm.NewStyle(pterm.FgGray)

	for _, item := range itemList {
		contextOut := contextStyle.Sprint(item)
		if itemType == "context" {
			tabItem := strings.Split(item, "@")

			contextOut = contextStyle.Sprint(tabItem[0]) + fileStyle.Sprint("@"+tabItem[1])
		}

		items = append(items, pterm.DefaultBasicText.Sprint(contextOut))
	}

	prompt := pterm.DefaultInteractiveSelect.
		WithOptions(items).
		WithDefaultText(fmt.Sprintf("Select %s:", itemType)).
		WithMaxHeight(len(items))

	if currentItem != "" {
		prompt = prompt.WithDefaultOption(pterm.DefaultBasicText.Sprint(contextStyle.Sprint(currentItem)))
	}

	output, err := prompt.Show()
	if err != nil {
		log.Debug().Err(err).Msg("cannot ask list")

		return "", err
	}

	return re.ReplaceAllString(output, ""), nil
}

// make ShellContextList sortable
func (a ShellContextList) Len() int           { return len(a) }
func (a ShellContextList) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ShellContextList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
