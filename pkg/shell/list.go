package shell

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rs/zerolog/log"
)

func LoadPterm(itemList []ContextDef) (ContextDef, error) {
	var (
		options     []string
		optionIndex int
	)

	for _, item := range itemList {
		options = append(options, fmt.Sprintf("%s (file: %s)", item.Name, item.FilePath))
	}

	simpleQs := &survey.Select{
		Message: "Select context:",
		Options: options,
	}

	err := survey.AskOne(simpleQs, &optionIndex)
	if err != nil {
		log.Debug().Err(err).Msg("cannot ask list")

		return ContextDef{}, err
	}

	return itemList[optionIndex], nil
}
