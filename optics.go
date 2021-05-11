// arsenal project optics.go
package main

import (
	"fmt"
	"os"
	"path"

	"text/template"

	"github.com/gocarina/gocsv"
)

type OpticsType string

func (o OpticsType) ToString() string {
	return string(o)
}

const (
	OpticsStarting OpticsType = "Starting"
	OpticsHolo     OpticsType = "Holo"
	OpticsCombat   OpticsType = "Combat"
	OpticsSniper   OpticsType = "Sniper"
)

var allOpticsTypes = [4]OpticsType{OpticsStarting, OpticsHolo, OpticsCombat, OpticsSniper}

type Optics struct {
	Type      OpticsType `csv:"type"`
	Name      string     `csv:"name"`
	MountType string     `csv:"mount"`
	Camoflage string     `csv:"camo"`
}

type TypeToOpticsList map[string][]*Optics

func processOptics(filepath *string, writeFiles bool, tplDir string) error {
	opticsList, err := parseOpticsCsv(filepath)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	OpticsListHolo := make(TypeToOpticsList)
	OpticsListCombat := make(TypeToOpticsList)
	OpticsListSniper := make(TypeToOpticsList)
	allOptics := map[string]TypeToOpticsList{
		"holo":   OpticsListHolo,
		"combat": OpticsListCombat,
		"sniper": OpticsListSniper,
	}

	for _, item := range opticsList {
		switch item.Type {
		case OpticsHolo:
			OpticsListHolo[item.MountType] = append(OpticsListHolo[item.MountType], item)
		case OpticsCombat:
			OpticsListCombat[item.MountType] = append(OpticsListCombat[item.MountType], item)
		case OpticsSniper:
			OpticsListSniper[item.MountType] = append(OpticsListSniper[item.MountType], item)
		}
	}

	oTpl, err := template.New("optics.tpl").ParseFiles(path.Join(tplDir, "optics.tpl"))
	if err != nil {
		return fmt.Errorf("failed to parse optics template: %w", err)
	}

	err = writeTemplate("universal", oTpl, allOptics, writeFiles, "")
	if err != nil {
		return fmt.Errorf("failed to execute optics template: %w", err)
	}

	return nil
}

func parseOpticsCsv(filepath *string) ([]*Optics, error) {
	csvFile, err := os.OpenFile(*filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	oList := []*Optics{}
	if err := gocsv.UnmarshalFile(csvFile, &oList); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return oList, nil
}
