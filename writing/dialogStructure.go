package writing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

const placeholder string = "Change me!"
const dialogFileDirectory string = "dialog"

var GlobalDialog *Dialog

type Dialog struct {
	Start struct {
		Greeting  string `json:"greeting"`
		Intro     string `json:"intro"`
		EnterName string `json:"enterName"`
		Info      string `json:"info"`
		Help      string `json:"help"`
	} `json:"start"`
	Game struct {
		GameHost struct {
			Update string `json:"update"`
			Toss   string `json:"toss"`
			Kill   string `json:"kill"`
			End    string `json:"end"`
			CantJoin string `json:"cant_join"`
			ChangeName string `json:"name_change"`
		} `json:"host"`
	} `json:"game"`
	Help struct {
		HelpMenu string `json:"menu"`
	} `json:"help"`
}

func BuildDialog() (error) {
	GlobalDialog = &Dialog{}
	jsonPath, err := getDialogPath(dialogFileDirectory)
	if err != nil {
		return err
	}
	if data, err := os.ReadFile(jsonPath); err == nil {
		_ = json.Unmarshal(data, GlobalDialog)
	}

	fillEmptyStrings(GlobalDialog)

	data, err := json.MarshalIndent(*GlobalDialog, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(jsonPath, data, 0644)
}

func fillEmptyStrings(ptr interface{}) {
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				field.Set(reflect.New(fieldType.Type.Elem()))
			}
			fillEmptyStrings(field.Interface())
		} else if field.Kind() == reflect.String {
			if field.String() == "" {
				field.SetString(placeholder)
			}
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
		return "", fmt.Errorf("Path exists but is not directory")
	}

	return filepath.Join(targetDir, "dialog.json"), nil
}
