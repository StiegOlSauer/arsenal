// arsenal project main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"text/template"
)

const (
	startingScore float64 = 100.0

	camoNocamo   string = "nocamo"
	camoBlack    string = "black"
	camoWoodland string = "woodland"
	camoDesert   string = "desert"
)

type ArmaEntityType interface {
	ToString() string
}

func main() {
	wpnsCsvPath := flag.String("weapons-csv", "", "path to CSV file with weapon weights")
	vestsCsvPath := flag.String("vests-csv", "", "path to CSV file with vests weights")
	opticsCsvPath := flag.String("optics-csv", "", "path to CSV file with optics list")
	uniformsCsvPath := flag.String("uniforms-csv", "", "path to CSV file with list of uniforms")
	devicesCsvPath := flag.String("devices-csv", "", "path to CSV file with list of devices")
	itemsCsvPath := flag.String("items-csv", "", "path to CSV file with list of common and explosive items")
	writeFiles := flag.Bool("w", false, "Write output to files instead of stdout (default: No)")
	templatesDir := flag.String("templates-dir", "templates", "path to directory where template files (.tpl) reside")
	flag.Parse()

	if _, err := os.Stat(*templatesDir); os.IsNotExist(err) {
		log.Fatal("Template directory %s does not exist, exiting!", templatesDir)
	}

	// do checks on provided files. If file exist: process entities. If both missing: fail
	if *wpnsCsvPath != "" {
		if _, err := os.Stat(*wpnsCsvPath); err == os.ErrNotExist {
			log.Fatal("File %s does not exist, exiting!", wpnsCsvPath)
		}
		err := processWeapons(wpnsCsvPath, *writeFiles, *templatesDir)
		if err != nil {
			log.Fatal("Failed to process weapons:", err)
		}
	}

	if *vestsCsvPath != "" {
		if _, err := os.Stat(*vestsCsvPath); err == os.ErrNotExist {
			log.Fatal("File %s does not exist, exiting!", vestsCsvPath)
		}
		err := processVests(vestsCsvPath, *writeFiles, *templatesDir)
		if err != nil {
			log.Fatal("Failed to process vests:", err)
		}
	}

	if *opticsCsvPath != "" {
		if _, err := os.Stat(*opticsCsvPath); err == os.ErrNotExist {
			log.Fatal("File %s does not exist, exiting!", opticsCsvPath)
		}
		err := processOptics(opticsCsvPath, *writeFiles, *templatesDir)
		if err != nil {
			log.Fatal("Failed to process optics:", err)
		}
	}

	if *uniformsCsvPath != "" {
		if _, err := os.Stat(*uniformsCsvPath); err == os.ErrNotExist {
			log.Fatal("File %s does not exist, exiting!", uniformsCsvPath)
		}
		err := processUniforms(uniformsCsvPath, *writeFiles, *templatesDir)
		if err != nil {
			log.Fatal("Failed to process uniforms:", err)
		}
	}

	if *devicesCsvPath != "" {
		if _, err := os.Stat(*devicesCsvPath); err == os.ErrNotExist {
			log.Fatal("File %s does not exist, exiting!", devicesCsvPath)
		}
		err := processDevices(devicesCsvPath, *writeFiles, *templatesDir)
		if err != nil {
			log.Fatal("Failed to process devices:", err)
		}
	}

	if *itemsCsvPath != "" {
		if _, err := os.Stat(*itemsCsvPath); err == os.ErrNotExist {
			log.Fatal("File %s does not exist, exiting!", itemsCsvPath)
		}
		err := processCommonItems(itemsCsvPath, *writeFiles, *templatesDir)
		if err != nil {
			log.Fatal("Failed to process common items:", err)
		}
	}
}

func createDir(dirname string) error {
	if err := os.Mkdir(dirname, os.ModePerm); !os.IsExist(err) {
		return err
	}
	return nil
}

func createSQFFile(dir string, fname string) *os.File {
	fname = strings.TrimSuffix(fname, filepath.Ext(fname)) + ".sqf"
	fh, err := os.Create(filepath.Join(dir, fname))
	if err != nil {
		log.Fatal("Cannot create file for writing:", fname, err)
	}
	return fh
}

func writeTemplate(subdir string, tpl *template.Template, tData interface{}, writeToFile bool, fileNameOverride string) error {
	if len(subdir) > 0 {
		if err := createDir(subdir); err != nil {
			return fmt.Errorf("failed to create directory for mission files: %w", err)
		}
	}

	fName := tpl.Name()
	if len(fileNameOverride) > 0 {
		fName = fileNameOverride
	}

	fh := os.Stdout
	if writeToFile {
		fh = createSQFFile(subdir, fName)
		defer fh.Close()
	}
	return tpl.Execute(fh, tData)
}

func getDivisor(r interface{}, fieldName string) (div float64) {
	t := reflect.TypeOf(r)
	field, _ := t.FieldByName(fieldName)

	div, err := strconv.ParseFloat(field.Tag.Get("weight_divisor"), 32)
	if err != nil {
		panic(err)
	}

	return
}

func joinTypes(base string, str ArmaEntityType, suffix string) string {
	spacer := ", "
	if len(base) == 0 {
		spacer = ""
	}

	return base + spacer + fmt.Sprintf("_%s%ss", strings.ToLower(str.ToString()), suffix)
}
