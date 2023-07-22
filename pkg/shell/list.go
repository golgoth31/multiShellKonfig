package shell

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func LoadList(itemType string, itemList []string) (string, error) {
	var (
		output string
	)

	items := []string{}
	contextStyle := pterm.NewStyle(pterm.FgLightBlue, pterm.Bold)
	fileStyle := pterm.NewStyle(pterm.FgGray)

	for _, item := range itemList {
		tabItem := strings.Split(item, "@")

		items = append(items, pterm.DefaultBasicText.Sprint(contextStyle.Sprint(tabItem[0])+fileStyle.Sprint("@"+tabItem[1])))

	}

	prompt := pterm.DefaultInteractiveSelect.
		WithOptions(items).
		WithDefaultText(fmt.Sprintf("Select %s:", itemType))

	output, err := prompt.Show()
	if err != nil {
		log.Debug().Err(err).Msg("cannot ask list")

		return "", err
	}

	return re.ReplaceAllString(output, ""), nil
}
