// arsenal project uniforms.go
package main

import (
	"fmt"
	"os"
	"path"

	"text/template"

	"github.com/gocarina/gocsv"
)

type Uniform struct {
	Name      string `csv:"name"`
	Camoflage string `csv:"camo"`
}

func processUniforms(filepath *string, writeFiles bool, tplDir string) error {
	uniformsList, err := parseUniformsCsv(filepath)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	uniformsNocamo := []string{}
	uniformsBlack := []string{}
	uniformsDesert := []string{}
	uniformsWoodland := []string{}

	camoToUniform := map[string][]string{
		"nocamo":   uniformsNocamo,
		"black":    uniformsBlack,
		"desert":   uniformsDesert,
		"woodland": uniformsWoodland,
	}

	for _, item := range uniformsList {
		switch item.Camoflage {
		case camoNocamo:
			camoToUniform[item.Camoflage] = append(camoToUniform[item.Camoflage], item.Name)
		case camoBlack:
			camoToUniform[item.Camoflage] = append(camoToUniform[item.Camoflage], item.Name)
		case camoDesert:
			camoToUniform[item.Camoflage] = append(camoToUniform[item.Camoflage], item.Name)
		case camoWoodland:
			camoToUniform[item.Camoflage] = append(camoToUniform[item.Camoflage], item.Name)
		}
	}

	uTpl, err := template.New("uniforms.tpl").ParseFiles(path.Join(tplDir, "uniforms.tpl"))
	if err != nil {
		return fmt.Errorf("failed to parse uniforms template: %w", err)
	}

	err = writeTemplate("universal", uTpl, camoToUniform, writeFiles, "")
	if err != nil {
		return fmt.Errorf("failed to execute uniforms template: %w", err)
	}

	return nil
}

func parseUniformsCsv(filepath *string) ([]*Uniform, error) {
	csvFile, err := os.OpenFile(*filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	oList := []*Uniform{}
	if err := gocsv.UnmarshalFile(csvFile, &oList); err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	return oList, nil
}
