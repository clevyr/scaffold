package main

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/clevyr/scaffold/appconfig"
	"regexp"
)

var validationRegex, _ = regexp.Compile("^[0-9]*[kmg]$")

func askQuestions(appConfig *appconfig.AppConfig) (err error) {
	// App Name
	err = survey.AskOne(&survey.Input{
		Message: "What is the application name?",
		Default: appConfig.AppName,
	}, &appConfig.AppName, survey.WithValidator(survey.Required))
	if err != nil {
		return
	}

	// Database
	err = survey.AskOne(&survey.Select{
		Message: "Choose which main database server to configure:",
		Options: []string{"PostgreSQL", "MariaDB"},
		Default: appConfig.Database,
	}, &appConfig.Database, survey.WithValidator(survey.Required))
	if err != nil {
		return
	}

	// Enabled PhpModules
	appConfig.EnableSelectedDatabase()
	err = survey.AskOne(&survey.MultiSelect{
		Message: "Choose which PHP modules to enable:",
		Options: appConfig.PhpModules.ToOptionsSlice(),
		Default: appConfig.PhpModules.ToDefaultSlice(),
	}, &appConfig.PhpModules)
	if err != nil {
		return
	}

	// Admin Gen
	err = survey.AskOne(&survey.Select{
		Message: "Choose which admin generator to include:",
		Options: []string{"None", "Nova", "Backpack"},
		Default: appConfig.AdminGen,
	}, &appConfig.AdminGen)
	if err != nil {
		return
	}

	// MailDev
	err = survey.AskOne(&survey.Confirm{
		Message: "Use MailDev as local mail backend?",
		Default: appConfig.MailDev,
		Help: "If enabled, MailDev will listen on http://localhost:1080 and Laravel will be configured accordingly.",
	}, &appConfig.MailDev)
	if err != nil {
		return
	}

	// Max Upload Size
	err = survey.AskOne(
		&survey.Input{
			Message: "What is the maximum upload size that should be allowed?",
			Default: appConfig.MaxUploadSize,
			Help: "Configures the maximum allowed upload size. " +
				"Supports the suffixes \"k\" (kilobytes), \"m\" (megabytes) and \"g\" (gigabytes).",
		},
		&appConfig.MaxUploadSize,
		survey.WithValidator(func(val interface{}) error {
			if str, ok := val.(string); !ok || !validationRegex.MatchString(str) {
				return errors.New("Make sure to enter a size followed by \"k\" (kilobytes), \"m\" (megabytes) or \"g\" (gigabytes).")
			}
			return nil
		}),
	)
	if err != nil {
		return
	}

	return
}