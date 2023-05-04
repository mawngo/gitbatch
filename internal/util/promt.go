package util

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AskConfig(key, message string) string {
	var value = viper.GetString(key)
	if value == "" {
		err := survey.AskOne(&survey.Input{Message: message}, &value, survey.WithValidator(survey.Required))
		cobra.CheckErr(err)
		viper.Set(key, value)
	}
	return value
}

func AskPassword(message string) string {
	var value = ""
	err := survey.AskOne(&survey.Password{Message: message}, &value, survey.WithValidator(survey.Required))
	cobra.CheckErr(err)
	return value
}
