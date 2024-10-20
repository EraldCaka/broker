package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
)

const (
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[0;34m"
	colorNone   = "\033[0m"
)

func main() {

	fileName := os.Getenv("GO_MOD_PATH")
	disablePublic := os.Getenv("DISABLE_PUBLIC")

	if fileName == "" {
		fileName = "../../go.mod"
	}

	body, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("unable to read file: %+v\n", err)
	}

	modFile, err := modfile.Parse("go.mod", body, nil)
	if err != nil {
		log.Fatalf("unable to read parse mod: %+v\n", err)
	}

	if len(modFile.Require) > 0 && disablePublic != "true" {
		log.Printf("%supdating public libraries . . .%s\n", colorBlue, colorNone)

		for _, v := range modFile.Require {

			if !strings.HasPrefix(v.Mod.Path, "github.com/EraldCaka/watermill-example") {

				var commandParams []string
				commandParams = append(commandParams, "get")
				commandParams = append(commandParams, "-u")
				commandParams = append(commandParams, v.Mod.Path)

				cmd := exec.Command("go", commandParams...) // #nosec

				cmdResult, err := cmd.Output()

				if len(cmdResult) == 0 && err != nil {
					log.Printf("error on update library %s: %+v\n", v.Mod.Path, err)
				} else {
					log.Printf("library successful updated %s\n", v.Mod.Path)
				}
			}
		}
	}

	time.Sleep(time.Second * 1)

	var commandParams []string
	commandParams = append(commandParams, "mod")
	commandParams = append(commandParams, "tidy")

	cmd := exec.Command("go", commandParams...) // #nosec

	cmdResult, err := cmd.Output()

	if len(cmdResult) == 0 && err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			log.Printf("%s%s%s", colorRed, exitError.Stderr, colorNone)
		}
		log.Printf("%serror on go mod tidy: %+v%s\n", colorRed, err, colorRed)
	} else {
		log.Printf("%sgo mod tidy executed successfully%s\n", colorGreen, colorNone)
	}
}
