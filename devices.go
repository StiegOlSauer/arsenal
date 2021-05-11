// arsenal project devices.go
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

type DeviceType string

func (v DeviceType) ToString() string {
	return string(v)
}

const (
	deviceTypeMuzzle DeviceType = "Muzzle"
	deviceTypeRail   DeviceType = "Rail"
)

type Device struct {
	Type       DeviceType `csv:"type"`
	Name       string     `csv:"name"`
	Camoflage  string     `csv:"camo"`
	MuzzleType float64    `csv:"muzzle_dev_type" weight_divisor:"8"`
	GripType   float64    `csv:"grip_type" weight_divisor:"4"`
	LightType  float64    `csv:"light_type" weight_divisor:"16"`
	Score      float64
}

func (v Device) getDivisor(fieldName string) (div float64) {
	return getDivisor(v, fieldName)
}

func processDevices(filepath *string, writeFiles bool, tplDir string) error {
	devList, err := parseDevicesCsv(filepath)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	devsBlack := make(map[DeviceType][]*Device)
	devsDesert := make(map[DeviceType][]*Device)
	devsWoodland := make(map[DeviceType][]*Device)
	devsNocamo := make(map[DeviceType][]*Device)
	devicesCamoTypeMap := map[string]map[DeviceType][]*Device{
		camoBlack:    devsBlack,
		camoDesert:   devsDesert,
		camoWoodland: devsWoodland,
		camoNocamo:   devsNocamo,
	}

	for _, dev := range devList {
		dev.Score = calculateDeviceWeight(dev)

		switch dev.Camoflage {
		case camoBlack:
			devsBlack[dev.Type] = append(devsBlack[dev.Type], dev)
		case camoDesert:
			devsDesert[dev.Type] = append(devsDesert[dev.Type], dev)
		case camoWoodland:
			devsWoodland[dev.Type] = append(devsWoodland[dev.Type], dev)
		case camoNocamo:
			devsNocamo[dev.Type] = append(devsNocamo[dev.Type], dev)
		}
	}

	devTpl, err := template.New("devices.tpl").Funcs(template.FuncMap{
		"IsLastItem": func(i, listSize int) bool { return i == listSize-1 },
		"ToLower":    func(dt DeviceType) string { return strings.ToLower(string(dt)) },
		"JoinTypes":  joinTypes,
	}).ParseFiles(path.Join(tplDir, "devices.tpl"))
	if err != nil {
		return fmt.Errorf("failed to parse devices template: %w", err)
	}

	err = writeTemplate("universal", devTpl, devicesCamoTypeMap, writeFiles, "")
	if err != nil {
		log.Fatal("Failed to execute devices template:", err)
	}

	return nil
}

func parseDevicesCsv(filepath *string) ([]*Device, error) {
	csvFile, err := os.OpenFile(*filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	devices := []*Device{}
	if err := gocsv.UnmarshalFile(csvFile, &devices); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return devices, nil
}

func calculateDeviceWeight(d *Device) float64 {
	score := startingScore

	if d.MuzzleType > 0 {
		score = score / (d.MuzzleType * d.getDivisor("MuzzleType"))
	}
	if d.GripType > 0 {
		score = score / (d.GripType * d.getDivisor("GripType"))
	}
	if d.LightType > 0 {
		score = score / (d.LightType * d.getDivisor("LightType"))
	}

	return score
}
