// arsenal project common_explosives.go
package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"text/template"

	"github.com/gocarina/gocsv"
)

const (
	itemExplosives     string = "Explosives"
	itemGrenadeFrag    string = "GrenadeFrag"
	itemGrenadeSmoke   string = "GrenadeSmoke"
	itemGrenadeOther   string = "GrenadeOther"
	itemCommon         string = "Common"
	itemMedicineFAK    string = "MedicineFAK"
	itemMedicinePAK    string = "MedicinePAK"
	itemBackpack       string = "Backpack"
	itemBinoculars     string = "Binoculars"
	itemAce            string = "ACEItem"
	itemAceGrenade     string = "ACEGrenade"
	itemAceMedicine    string = "ACEMedicine"
	itemAceMedDressing string = "ACEMedDressing"
	itemAceMedInjector string = "ACEMedInjector"
	GP25HEMagazine     string = "GP25_HE"
	GP25OtherMagazine  string = "GP25_other"
	M203HEMagazine     string = "M203_HE"
	M203OtherMagazine  string = "M203_other"
)

type Item struct {
	Category string `csv:"type"`
	Name     string `csv:"name"`
}

func processCommonItems(filepath *string, writeFiles bool, tplDir string) error {
	itemsList, err := parseItemsCsv(filepath)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	commonItems := map[string][]string{
		itemCommon:        []string{},
		itemMedicineFAK:   []string{},
		itemMedicinePAK:   []string{},
		itemBackpack:      []string{},
		itemBinoculars:    []string{},
		GP25HEMagazine:    []string{},
		GP25OtherMagazine: []string{},
		M203HEMagazine:    []string{},
		M203OtherMagazine: []string{},
	}

	explosiveItems := map[string][]string{
		itemExplosives:   []string{},
		itemGrenadeFrag:  []string{},
		itemGrenadeSmoke: []string{},
		itemGrenadeOther: []string{},
	}

	ACEItems := map[string][]string{
		itemAce:            []string{},
		itemAceGrenade:     []string{},
		itemAceMedicine:    []string{},
		itemAceMedDressing: []string{},
		itemAceMedInjector: []string{},
	}

	for _, item := range itemsList {
		switch item.Category {
		case itemExplosives:
			explosiveItems[item.Category] = append(explosiveItems[item.Category], item.Name)
		case itemGrenadeFrag:
			explosiveItems[item.Category] = append(explosiveItems[item.Category], item.Name)
		case itemGrenadeSmoke:
			explosiveItems[item.Category] = append(explosiveItems[item.Category], item.Name)
		case itemGrenadeOther:
			explosiveItems[item.Category] = append(explosiveItems[item.Category], item.Name)
		case itemCommon:
			commonItems[item.Category] = append(commonItems[item.Category], item.Name)
		case itemMedicinePAK:
			commonItems[item.Category] = append(commonItems[item.Category], item.Name)
		case itemMedicineFAK:
			commonItems[item.Category] = append(commonItems[item.Category], item.Name)
		case itemBackpack:
			commonItems[item.Category] = append(commonItems[item.Category], item.Name)
		case itemBinoculars:
			commonItems[item.Category] = append(commonItems[item.Category], item.Name)
		case itemAce:
			ACEItems[item.Category] = append(ACEItems[item.Category], item.Name)
		case itemAceGrenade:
			ACEItems[item.Category] = append(ACEItems[item.Category], item.Name)
		case itemAceMedicine:
			ACEItems[item.Category] = append(ACEItems[item.Category], item.Name)
		case itemAceMedDressing:
			ACEItems[item.Category] = append(ACEItems[item.Category], item.Name)
		case itemAceMedInjector:
			ACEItems[item.Category] = append(ACEItems[item.Category], item.Name)
		case GP25HEMagazine:
			commonItems[item.Category] = append(commonItems[item.Category], item.Name)
		case GP25OtherMagazine:
			commonItems[item.Category] = append(commonItems[item.Category], item.Name)
		case M203HEMagazine:
			commonItems[item.Category] = append(commonItems[item.Category], item.Name)
		case M203OtherMagazine:
			commonItems[item.Category] = append(commonItems[item.Category], item.Name)
		}
	}

	iTpl, err := template.New("common_explosives.tpl").Funcs(template.FuncMap{
		"ToLower": func(s string) string { return strings.ToLower(s) },
	}).ParseFiles(path.Join(tplDir, "common_explosives.tpl"))
	if err != nil {
		return fmt.Errorf("failed to parse common items template: %w", err)
	}

	err = writeTemplate("universal", iTpl, explosiveItems, writeFiles, "explosives")
	if err != nil {
		return fmt.Errorf("failed to execute explosive items template: %w", err)
	}

	err = writeTemplate("", iTpl, commonItems, writeFiles, "common_items")
	if err != nil {
		return fmt.Errorf("failed to execute common items template: %w", err)
	}

	err = writeTemplate("", iTpl, ACEItems, writeFiles, "ace_items")
	if err != nil {
		return fmt.Errorf("failed to execute ACE items template: %w", err)
	}

	return nil
}

func parseItemsCsv(filepath *string) ([]*Item, error) {
	csvFile, err := os.OpenFile(*filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	items := []*Item{}
	if err := gocsv.UnmarshalFile(csvFile, &items); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return items, nil
}
