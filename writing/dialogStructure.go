package writing

import (
	"gopkg.in/yaml.v3"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

const placeholder string = "Change me!"
const dialogFileDirectory string = "dialog"
const filename string = "dialog.yaml"

var GlobalDialog *Dialog

type Dialog struct {
	Start struct {
		Greeting  string `yaml:"greeting"`
		Intro     string `yaml:"intro"`
		EnterName string `yaml:"enterName"`
		Info      string `yaml:"info"`
		Help      string `yaml:"help"`
	} `yaml:"start"`
	Game struct {
		GameHost struct {
			Update string `yaml:"update"`
			Toss   string `yaml:"toss"`
			Kill   string `yaml:"kill"`
			End    string `yaml:"end"`
			CantJoin string `yaml:"cant_join"`
			ChangeName string `yaml:"name_change"`
		} `yaml:"host"`
	} `yaml:"game"`
	Help struct {
		HelpMenu string `yaml:"menu"`
	} `yaml:"help"`
}

func BuildDialog() (error) {
	GlobalDialog = &Dialog{}
	yamlPath, err := getDialogPath(dialogFileDirectory)
	if err != nil {
		return err
	}
	if data, err := os.ReadFile(yamlPath); err == nil {
		_ = yaml.Unmarshal(data, GlobalDialog)
	}
	
	fillEmptyStrings(reflect.ValueOf(GlobalDialog))

	data, err := yaml.Marshal(*GlobalDialog)

	if err != nil {
		return err
	}

	return os.WriteFile(yamlPath, data, 0644)
}

func fillEmptyStrings(v reflect.Value) {
	if !v.IsValid(){
		return
	}
	
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		fillEmptyStrings(v.Elem())
		return
	}

	switch v.Kind(){
	case reflect.Struct:
		for i:= 0; i < v.NumField(); i++{
			field := v.Field(i)

			if field.CanSet() || field.Kind() == reflect.Struct || field.Kind() == reflect.Ptr {
				fillEmptyStrings(field)
			} else {
				if field.CanAddr() {
					fillEmptyStrings((field.Addr()))
				}
			}
		}
	case reflect.String:
		if v.String() == ""{
			v.SetString(placeholder)
		}
	case reflect.Slice, reflect.Array:
		for i:= 0 ; i < v.Len(); i++{
			fillEmptyStrings(v.Index(i))
		}
	}
}

func getBaseDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(exePath)
}

func getDialogPath(directory string) (string, error) {
	base := getBaseDir()
	targetDir := filepath.Join(base, "..", directory)

	// Check if the directory exists
	info, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		// Create the directory if missing
		if mkErr := os.MkdirAll(targetDir, 0755); mkErr != nil {
			return "", mkErr
		}
	} else if err != nil {
		return "", err
	} else if !info.IsDir() {
		return "", fmt.Errorf("path exists but is not directory")
	}

	return filepath.Join(targetDir, filename), nil
}
