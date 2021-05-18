// arsenal project weapons.go
package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"text/template"

	"github.com/gocarina/gocsv"
)

type WeaponType string

func (w WeaponType) ToString() string {
	return string(w)
}

type WeaponRarenessBit uint8

const (
	wpnModern WeaponRarenessBit = 1 << iota
	wpnSpecops
	wpnUnique
)

type DeviceMounts uint8

const (
	wpnAllowsDevices DeviceMounts = 1 << iota
	integralBipod
	integralFrontgrip
)

const (
	wpnTypeRifle     WeaponType = "Rifle"
	wpnTypeRifleGL   WeaponType = "RifleGL"
	wpnTypeCarbine   WeaponType = "Carbine"
	wpnTypeCarbineGL WeaponType = "CarbineGL"
	wpnTypePistol    WeaponType = "Pistol"
	wpnTypeSniper    WeaponType = "Sniper"
	wpnTypeMarksman  WeaponType = "Marksman"
	wpnTypeSmg       WeaponType = "SMG"
	wpnTypeLMG       WeaponType = "LMG"
	wpnTypeMG        WeaponType = "MG"
	wpnTypeAA        WeaponType = "AA"
	wpnTypeLAT       WeaponType = "LAT"
	wpnTypeHAT       WeaponType = "HAT"
)

var allWpnTypes = [12]WeaponType{
	wpnTypeRifle, wpnTypeRifleGL, wpnTypeCarbine, wpnTypeCarbineGL, wpnTypePistol,
	wpnTypeSniper, wpnTypeSmg, wpnTypeLMG, wpnTypeMG, //wpnTypeStarting,
	wpnTypeAA, wpnTypeLAT, wpnTypeHAT,
}

type RarenessProperties struct {
	Rareness  WeaponRarenessBit
	IsModern  bool `weight_divisor:"2"`
	IsSpecops bool `weight_divisor:"1.5"`
	IsUnique  bool `weight_divisor:"3.5"`
}

func (r *RarenessProperties) UnmarshalCSV(csv string) error {
	s, err := strconv.ParseUint(csv, 10, 8)
	if err != nil {
		return err
	}
	r.IsModern = WeaponRarenessBit(s)&wpnModern != 0
	r.IsSpecops = WeaponRarenessBit(s)&wpnSpecops != 0
	r.IsUnique = WeaponRarenessBit(s)&wpnUnique != 0

	return nil
}

func (r RarenessProperties) getDivisor(fieldName string) (div float64) {
	return getDivisor(r, fieldName)
}

type MountProperties struct {
	AllowsMuzzleDevices bool `csv:"allows_muzzle_dev" weight_divisor:"1.5"`
	AllowsOptics        bool `weight_divisor:"2"`
}

type AdditionalRails struct {
	HasSideRail        bool `csv:"has_side_rail" weight_divisor:"1.5"`
	HasIntegralToprail bool `csv:"has_int_toprail" weight_divisor:"2"`
}

type IntegratedFeatures struct {
	Features          DeviceMounts
	AllowsGripsBipods bool `weight_divisor:"1.5"`
	HasBipod          bool `weight_divisor:"2"`
	HasGL             bool `weight_divisor:"2.5"`
	HasFrontGrip      bool `weight_divisor:"2"`
	HasOptics         bool `weight_divisor:"1.5"`
}

func (f *IntegratedFeatures) UnmarshalCSV(csv string) error {
	s, err := strconv.ParseUint(csv, 10, 8)
	if err != nil {
		return err
	}
	f.AllowsGripsBipods = DeviceMounts(s)&wpnAllowsDevices != 0
	f.HasBipod = DeviceMounts(s)&integralBipod != 0
	f.HasFrontGrip = DeviceMounts(s)&integralFrontgrip != 0

	return err
}

func (f IntegratedFeatures) getDivisor(fieldName string) (div float64) {
	return getDivisor(f, fieldName)
}

type MechanicalFeatures struct {
	IsBalanced        bool `csv:"balanced" weight_divisor:"4"`
	HighRateOfFire    bool `csv:"high_rof" weight_divisor:"2"`
	IsSilenced        bool `csv:"silenced" weight_divisor:"1.5"`
	HighCapacityMag   bool `csv:"hi_cap_mag" weight_divisor:"2"`
	PowerfulCartridge bool `csv:"powerful_cartridge" weight_divisor:"2"`
}

type Weapon struct {
	Type       WeaponType `csv:"type"`
	Name       string     `csv:"name"`
	MagType    string     `csv:"mag_type"`
	RailType   string     `csv:"optics_rail" weight_divisor:"1.5"`
	Camoflage  string     `csv:"camo"`
	RawScore   float64    `csv:"raw_score"`
	IsStarting int        `csv:"is_starting"`
	Score      float64
	Rareness   RarenessProperties `csv:"rareness"`
	Features   IntegratedFeatures `csv:"frontrail_devices"`
	AdditionalRails
	MountProperties
	// IntegratedFeatures
	MechanicalFeatures
}

func (w Weapon) getDivisor(fieldName string) (div float64) {
	return getDivisor(w, fieldName)
}

func processWeapons(filepath *string, writeFiles bool, tplDir string) error {
	weaponList, err := parseWeaponCsv(filepath)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	weaponsDesert := make(map[WeaponType][]*Weapon)
	weaponsBlack := make(map[WeaponType][]*Weapon)
	weaponsWoodland := make(map[WeaponType][]*Weapon)
	weaponsNocamo := make(map[WeaponType][]*Weapon)
	magazineTypes := []string{}
	camoToWeapons := map[string]map[WeaponType][]*Weapon{
		"desert":   weaponsDesert,
		"black":    weaponsBlack,
		"woodland": weaponsWoodland,
		"nocamo":   weaponsNocamo,
	}

	// Init weapon types to ensure that all are represented (limitation of mission scripts)
	for _, wpnMap := range camoToWeapons {
		for _, wpnType := range allWpnTypes {
			wpnMap[wpnType] = []*Weapon{}
		}
	}

	for _, wpn := range weaponList {
		wpn.Score = calculateWeaponWeight(wpn)

		switch wpn.Camoflage {
		case camoNocamo:
			weaponsNocamo[wpn.Type] = append(weaponsNocamo[wpn.Type], wpn)
		case camoBlack:
			weaponsBlack[wpn.Type] = append(weaponsBlack[wpn.Type], wpn)
		case camoDesert:
			weaponsDesert[wpn.Type] = append(weaponsDesert[wpn.Type], wpn)
		case camoWoodland:
			weaponsWoodland[wpn.Type] = append(weaponsWoodland[wpn.Type], wpn)
		}

		isNew := true
		for _, m := range magazineTypes {
			if wpn.MagType == m {
				isNew = false
				break
			}
		}
		if isNew {
			magazineTypes = append(magazineTypes, wpn.MagType)
		}
	}

	wTpl, err := template.New("weapons.tpl").Funcs(template.FuncMap{
		"IsLastItem": func(i, listSize int) bool { return i == listSize-1 },
		"ToLower":    func(wt WeaponType) string { return strings.ToLower(string(wt)) },
		"JoinTypes":  joinTypes,
	}).ParseFiles(path.Join(tplDir, "weapons.tpl"))
	if err != nil {
		return fmt.Errorf("failed to parse weapons template: %w", err)
	}

	for pattern, wpnMap := range camoToWeapons {
		err = writeTemplate(pattern, wTpl, wpnMap, writeFiles, "")
		if err != nil {
			return fmt.Errorf("failed to execute weapons template: %w", err)
		}
	}

	magTpl, err := template.New("magazines.tpl").Funcs(template.FuncMap{
		"StringsJoin": strings.Join,
	}).ParseFiles(path.Join(tplDir, "magazines.tpl"))
	if err != nil {
		return fmt.Errorf("failed to parse magazines template: %w", err)
	}

	err = writeTemplate("universal", magTpl, magazineTypes, writeFiles, "")
	if err != nil {
		return fmt.Errorf("failed to execute magazines template: %w", err)
	}

	return nil
}

func calculateWeaponWeight(wpn *Weapon) float64 {
	score := startingScore
	switch wpn.RailType {
	case "RIS":
		score = score / wpn.getDivisor("RailType")
		score = score / wpn.getDivisor("AllowsOptics")
	case "Dovetail":
		score = score / wpn.getDivisor("AllowsOptics")
	case "Integral":
		score = score / wpn.Features.getDivisor("HasOptics")
	}

	// Uniqueness
	if wpn.Rareness.IsModern {
		score = score / wpn.Rareness.getDivisor("IsModern")
	}
	if wpn.Rareness.IsSpecops {
		score = score / wpn.Rareness.getDivisor("IsSpecops")
	}
	if wpn.Rareness.IsUnique {
		score = score / wpn.Rareness.getDivisor("IsUnique")
	}

	// Allowed mounts
	if wpn.AllowsMuzzleDevices {
		score = score / wpn.getDivisor("AllowsMuzzleDevices")
	}
	if wpn.Features.AllowsGripsBipods {
		score = score / wpn.Features.getDivisor("AllowsGripsBipods")
	}

	// Integrated hardware
	if wpn.Features.HasBipod {
		score = score / wpn.Features.getDivisor("HasBipod")
	}
	if wpn.Features.HasFrontGrip {
		score = score / wpn.Features.getDivisor("HasFrontGrip")
	}
	if wpn.HasSideRail {
		score = score / wpn.getDivisor("HasSideRail")
	}
	if wpn.HasIntegralToprail {
		score = score / wpn.getDivisor("HasIntegralToprail")
	}
	if wpn.Type == wpnTypeRifleGL || wpn.Type == wpnTypeCarbineGL {
		score = score / wpn.Features.getDivisor("HasGL")
	}
	// Mechanical features
	if wpn.IsBalanced {
		score = score / wpn.getDivisor("IsBalanced")
	}
	if wpn.HighRateOfFire {
		score = score / wpn.getDivisor("HighRateOfFire")
	}
	if wpn.IsSilenced {
		score = score / wpn.getDivisor("IsSilenced")
	}
	if wpn.HighCapacityMag {
		score = score / wpn.getDivisor("HighCapacityMag")
	}
	if wpn.PowerfulCartridge {
		score = score / wpn.getDivisor("PowerfulCartridge")
	}

	return score
}

func parseWeaponCsv(filepath *string) ([]*Weapon, error) {
	csvFile, err := os.OpenFile(*filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	wpnList := []*Weapon{}
	if err := gocsv.UnmarshalFile(csvFile, &wpnList); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return wpnList, nil
}
