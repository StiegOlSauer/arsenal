// arsenal project vehicles.go
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
	// Vehicle roles (for battlegroup compositions)
	roleRespawn   = "respawn_veh"
	roleRepair    = "repair_veh"
	roleFOB       = "fob_veh"
	roleTransport = "transport"
	roleSupport   = "support"
	roleATGM      = "atgm"
	roleAPC       = "apc"
	roleIFV       = "ifv"
	vehRoleAA     = "aa"
	roleArtillery = "artillery"
	roleTank      = "tank"
	roleAttack    = "attack"
)

type VehicleCamo struct {
	CamoList string
}

func (vc *VehicleCamo) UnmarshalCSV(csv string) error {
	vc.CamoList = fmt.Sprintf("[\"%s\"]", strings.Join(strings.Split(csv, ","), "\",\""))
	return nil
}

type Vehicle struct {
	Type          string      `csv:"type"`
	Role          string      `csv:"role"`
	Name          string      `csv:"name"`
	Class         string      `csv:"classname"`
	Family        string      `csv:"family"`
	Cost          int         `csv:"cost"`
	FuelCost      int         `csv:"fuel_cost"`
	UpgradeCost   int         `csv:"upgrade_cost"`
	Rareness      float64     `csv:"rareness"`
	CamosWoodland VehicleCamo `csv:"woodland"`
	CamosDesert   VehicleCamo `csv:"desert"`
	CamosNocamo   VehicleCamo `csv:"nocamo"`
	IsUnlocked    bool        `csv:"is_starting"`
	RawProperties string      `csv:"special"`
}

func processVehicles(filepath *string, writeFiles bool, tplDir string) error {
	vehList, err := parseVehiclesCsv(filepath)
	vehiclesByRole := make(map[string][]*Vehicle)

	for _, veh := range vehList {
		switch veh.Role {
		case roleRespawn:
			vehiclesByRole[roleRespawn] = append(vehiclesByRole[roleRespawn], veh)
		case roleRepair:
			vehiclesByRole[roleRepair] = append(vehiclesByRole[roleRepair], veh)
		case roleFOB:
			vehiclesByRole[roleFOB] = append(vehiclesByRole[roleFOB], veh)
		case roleTransport:
			vehiclesByRole[roleTransport] = append(vehiclesByRole[roleTransport], veh)
		case roleSupport:
			vehiclesByRole[roleSupport] = append(vehiclesByRole[roleSupport], veh)
		case roleATGM:
			vehiclesByRole[roleATGM] = append(vehiclesByRole[roleATGM], veh)
		case roleAPC:
			vehiclesByRole[roleAPC] = append(vehiclesByRole[roleAPC], veh)
		case roleIFV:
			vehiclesByRole[roleIFV] = append(vehiclesByRole[roleIFV], veh)
		case vehRoleAA:
			vehiclesByRole[vehRoleAA] = append(vehiclesByRole[vehRoleAA], veh)
		case roleArtillery:
			vehiclesByRole[roleArtillery] = append(vehiclesByRole[roleArtillery], veh)
		case roleTank:
			vehiclesByRole[roleTank] = append(vehiclesByRole[roleTank], veh)
		case roleAttack:
			vehiclesByRole[roleAttack] = append(vehiclesByRole[roleAttack], veh)
		}
	}

	vehTpl, err := template.New("vehicles.tpl").ParseFiles(path.Join(tplDir, "vehicles.tpl"))
	if err != nil {
		return fmt.Errorf("failed to parse vehicles template: %w", err)
	}

	err = writeTemplate("", vehTpl, vehiclesByRole, writeFiles, "")
	if err != nil {
		return fmt.Errorf("failed to execute vehicles template: %w", err)
	}

	return err
}

func parseVehiclesCsv(filepath *string) ([]*Vehicle, error) {
	csvFile, err := os.OpenFile(*filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	vehs := []*Vehicle{}
	if err := gocsv.UnmarshalFile(csvFile, &vehs); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return vehs, nil
}
