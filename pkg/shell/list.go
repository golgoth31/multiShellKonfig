package shell

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rs/zerolog/log"
)

func LoadList(itemType string, itemList []string) (string, error) {
	var (
		// options     []string
		output string
	)

	// for _, item := range itemList {
	// 	options = append(options, item)
	// }

	simpleQs := &survey.Select{
		Message:  fmt.Sprintf("Select %s:", itemType),
		Options:  itemList,
		PageSize: len(itemList),
	}

	err := survey.AskOne(simpleQs, &output)
	if err != nil {
		log.Debug().Err(err).Msg("cannot ask list")

		return "", err
	}

	return output, nil
}
