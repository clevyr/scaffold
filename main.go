//go:generate pkger
package main

import (
	"flag"
	"fmt"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/clevyr/scaffold/appconfig"
	"os"
	"path"
)

func main() {
	var err error

	var context string
	flag.StringVar(&context, "C", ".", "Run as if the application was started in the given path.")

	var versionFlag bool
	flag.BoolVar(&versionFlag, "v", false, "Prints the current versionFlag.")

	flag.Parse()

	if versionFlag {
		printVersion()
		return
	}

	context = path.Clean(context)
	if context != "." {
		err = os.MkdirAll(context, os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.Chdir(context)
		if err != nil {
			panic(err)
		}
	}

	appConfig := appconfig.Defaults
	_ = appConfig.ImportFromFile()
	err = askQuestions(&appConfig)
	if err == terminal.InterruptErr {
		fmt.Println("Interrupted")
		return
	} else if err != nil {
		panic(err)
	}

	err = os.Setenv("COMPOSER_MEMORY_LIMIT", "-1")
	if err != nil {
		panic(err)
	}

	err = appConfig.GenerateAppKey()
	if err != nil {
		panic(err)
	}

	err = initLaravel(appConfig)
	if err != nil {
		panic(err)
	}

	err = appConfig.ExportToFile()
	if err != nil {
		panic(err)
	}
	
	err = generateTemplate(appConfig)
	if err != nil {
		panic(err)
	}

	err = composerRequire(appConfig)
	if err != nil {
		panic(err)
	}

	err = npmInstall()
	if err != nil {
		panic(err)
	}
}
