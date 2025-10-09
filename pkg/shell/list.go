package shell

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog/log"
)

func LoadList(itemType string, currentItem string, itemList []string) (string, error) {
	// Build plain, unstyled options; display equals value for simplicity
	items := make([]string, 0, len(itemList))
	for _, item := range itemList {
		// Keep the raw item (e.g., "context@file") to return an unambiguous value
		items = append(items, item)
	}

	// Pre-set selected value to enable default selection behavior
	selected := strings.TrimSpace(currentItem)

	selectField := huh.NewSelect[string]().
		Title(fmt.Sprintf("Select %s:", itemType)).
		Options(huh.NewOptions(items...)...).
		Value(&selected)

	if err := huh.NewForm(huh.NewGroup(selectField)).Run(); err != nil {
		log.Debug().Err(err).Msg("cannot ask list")
		return "", err
	}

	return strings.TrimSpace(selected), nil
}

// make ShellContextList sortable
func (a ShellContextList) Len() int           { return len(a) }
func (a ShellContextList) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ShellContextList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
