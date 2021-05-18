// arsenal project squads.go
package main

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"text/template"

	"github.com/gocarina/gocsv"
)

type RoleName string

const (
	roleTL             RoleName = "tl"
	roleTLHvy          RoleName = "tl_heavy"
	roleRifleLight     RoleName = "rifleman_light"
	roleRifle          RoleName = "rifleman"
	roleRifleHvy       RoleName = "rifleman_heavy"
	roleGL             RoleName = "grenadier"
	roleGLHvy          RoleName = "grenadier_heavy"
	roleSharp          RoleName = "sharpshooter"
	roleMark           RoleName = "marksman"
	roleLAT            RoleName = "lat"
	roleHAT            RoleName = "hat"
	roleAA             RoleName = "aa"
	roleMedic          RoleName = "medic"
	roleLMG            RoleName = "lmg"
	roleMG             RoleName = "mg"
	roleSniper         RoleName = "sniper"
	roleCrewSMG        RoleName = "crew_smg"
	roleCrewCarbine    RoleName = "crew_carbine"
	roleSpecTL         RoleName = "spec_tl"
	roleSpecRifleLight RoleName = "spec_rifleman_light"
	roleSpecRifle      RoleName = "spec_rifleman"
	roleSpecGL         RoleName = "spec_grenadier"
	roleSpecMark       RoleName = "spec_marksman"
	roleSpecSniper     RoleName = "spec_sniper"
	roleSpecMG         RoleName = "spec_mg"
	roleSpecSharp      RoleName = "spec_sharpshooter"
	roleSpecLAT        RoleName = "spec_lat"
	roleSpecHAT        RoleName = "spec_hat"
	roleSpecAA         RoleName = "spec_aa"
	roleSpecMedic      RoleName = "spec_medic"
	roleSpecLMG        RoleName = "spec_lmg"
)

type RoleTraits struct {
	Traits string
}

func (t *RoleTraits) UnmarshalCSV(csv string) error {
	if len(csv) > 0 {
		t.Traits = fmt.Sprintf("\"%s\"", strings.Join(strings.Split(csv, ","), "\",\""))
	}
	return nil
}

type Role struct {
	Name         RoleName   `csv:"name"`
	WpnClass     string     `csv:"pri_wpn_class"`
	NMags        int        `csv:"nmags"`
	Optics       string     `csv:"optics"`
	MuzzleDevice bool       `csv:"muzzle_device"`
	RailDevice   bool       `csv:"frontrail_device"`
	NGrenades    int        `csv:"ngrenades"`
	VestType     string     `csv:"vest"`
	Backpack     bool       `csv:"has_backpack"`
	Traits       RoleTraits `csv:"traits"`

	BaseSquad    int `csv:"base_sq"`
	LightSquad   int `csv:"light_sq"`
	HeavySquad   int `csv:"heavy_sq"`
	SniperFT     int `csv:"sniper_ft"`
	ATSquad      int `csv:"at_sq"`
	AASquad      int `csv:"aa_sq"`
	SpecFT       int `csv:"spec_ft"`
	SpecSniperFT int `csv:"spec_sniper_ft"`
	PatrolSquad  int `csv:"patrol_sq"`
	AmbushSquad  int `csv:"ambush_sq"`
	CrewComp     int `csv:"crew_comp"`
	SentryComp   int `csv:"sentry_comp"`
}

func (r Role) getCsvFieldName(strField string) string {
	return getCsvFieldName(r, strField)
}

type Squad []RoleName

func processSquads(filepath *string, writeFiles bool, tplDir string) error {
	rolesList, err := parseSquadsCsv(filepath)
	allSquads := make(map[string]Squad)

	for _, r := range rolesList {
		fields := []int{
			r.BaseSquad, r.LightSquad, r.HeavySquad, r.SniperFT, r.ATSquad,
			r.AASquad, r.SpecFT, r.SpecSniperFT, r.PatrolSquad, r.AmbushSquad,
			r.CrewComp, r.SentryComp,
		}

		val := reflect.Indirect(reflect.ValueOf(r))
		for i, f := range fields {
			fName := r.getCsvFieldName(val.Type().Field(10 + i).Name)
			for j := 0; j < f; j++ {
				allSquads[fName] = append(allSquads[fName], r.Name)
			}
		}
	}

	sTpl, err := template.New("squads.tpl").Funcs(template.FuncMap{
		"IsLastItem": func(i, listSize int) bool { return i == listSize-1 },
	}).ParseFiles(path.Join(tplDir, "squads.tpl"))
	if err != nil {
		return fmt.Errorf("failed to parse roles and squads template: %w", err)
	}

	payload := struct {
		Roles  []*Role
		Squads map[string]Squad
	}{Roles: rolesList, Squads: allSquads}
	err = writeTemplate("", sTpl, payload, writeFiles, "")
	if err != nil {
		return fmt.Errorf("failed to execute roles and squads template: %w", err)
	}

	return err
}

func parseSquadsCsv(filepath *string) ([]*Role, error) {
	csvFile, err := os.OpenFile(*filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	roles := []*Role{}
	if err := gocsv.UnmarshalFile(csvFile, &roles); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return roles, nil
}
