package shell

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/rs/zerolog/log"
)

func LoadList(itemList []string) (int, error) {
	var (
		options     []string
		optionIndex int
	)

	for _, item := range itemList {
		options = append(options, item)
	}

	simpleQs := &survey.Select{
		Message: "Select context:",
		Options: options,
	}

	err := survey.AskOne(simpleQs, &optionIndex)
	if err != nil {
		log.Debug().Err(err).Msg("cannot ask list")

		return 0, err
	}

	return optionIndex, nil
}
