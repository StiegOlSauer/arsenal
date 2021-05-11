// arsenal project vests.go
package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/gocarina/gocsv"
)

type VestType string

func (v VestType) ToString() string {
	return string(v)
}

const (
	vestTypeNormal   VestType = "Normal"
	vestTypeHeavy    VestType = "Heavy"
	vestTypeHeadgear VestType = "Headgear"
)

type Vest struct {
	Type        VestType `csv:"type"`
	Name        string   `csv:"name"`
	Weight      float64  `csv:"weight" weight_divisor:"0.5"`
	BProtection float64  `csv:"ballistic_protection" weight_divisor:"2.5"`
	EProtection float64  `csv:"explosive_protection" weight_divisor:"0.5"`
	Capacity    float64  `csv:"capacity" weight_divisor:"1"`
	Camoflage   string   `csv:"camo"`
	Score       float64
}

func (v Vest) getDivisor(fieldName string) (div float64) {
	return getDivisor(v, fieldName)
}

func processVests(filepath *string, writeFiles bool, tplDir string) error {
	vestList, err := parseVestsCsv(filepath)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	vestsH := make(map[string][]*Vest)
	vestsN := make(map[string][]*Vest)
	vestsHdg := make(map[string][]*Vest)
	typeToCamo := map[VestType]map[string][]*Vest{
		vestTypeHeavy:    vestsH,
		vestTypeNormal:   vestsN,
		vestTypeHeadgear: vestsHdg,
	}

	for _, vest := range vestList {
		vest.Score = calculateVestWeight(vest)

		switch vest.Type {
		case vestTypeHeavy:
			vestsH[vest.Camoflage] = append(vestsH[vest.Camoflage], vest)
		case vestTypeNormal:
			vestsN[vest.Camoflage] = append(vestsN[vest.Camoflage], vest)
		case vestTypeHeadgear:
			vestsHdg[vest.Camoflage] = append(vestsHdg[vest.Camoflage], vest)
		}
	}

	vTpl, err := template.New("vests.tpl").Funcs(template.FuncMap{
		"IsLastItem": func(i, listSize int) bool { return i == listSize-1 },
		"ToLower":    func(vt VestType) string { return strings.ToLower(string(vt)) },
		"JoinTypes":  joinTypes,
	}).ParseFiles(path.Join(tplDir, "vests.tpl"))
	if err != nil {
		return fmt.Errorf("failed to parse vests template: %w", err)
	}

	err = writeTemplate("universal", vTpl, typeToCamo, writeFiles, "")
	if err != nil {
		log.Fatal("Failed to execute devices template:", err)
	}

	return nil
}

func calculateVestWeight(v *Vest) float64 {
	score := startingScore

	if v.BProtection > 0 {
		score = score / (v.BProtection * v.getDivisor("BProtection"))
	}
	if v.EProtection > 0 {
		score = score / (v.EProtection * v.getDivisor("EProtection"))
	}
	if v.Capacity > 0 {
		score = score / (v.Capacity * v.getDivisor("Capacity"))
	}
	if v.Weight > 0 {
		// score = score / (v.Weight * v.getDivisor("Weight"))
		score = score * (v.Weight * v.getDivisor("Weight"))
	}

	return score
}

func parseVestsCsv(filepath *string) ([]*Vest, error) {
	csvFile, err := os.OpenFile(*filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	vList := []*Vest{}
	if err := gocsv.UnmarshalFile(csvFile, &vList); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return vList, nil
}
