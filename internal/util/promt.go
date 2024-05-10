package util

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func AskToken() string {
	mode := GetMode()
	return AskConfig(mode+".token", fmt.Sprintf("Enter %s token", cases.Title(language.English, cases.NoLower).String(GetMode())))
}

func askPassword() string {
	var value = ""
	err := survey.AskOne(&survey.Password{Message: fmt.Sprintf("Enter %s password", cases.Title(language.English, cases.NoLower).String(GetMode())) + " (this won't be saved)"}, &value, survey.WithValidator(survey.Required))
	cobra.CheckErr(err)
	return value
}

func AskAuth() transport.AuthMethod {
	if IsSSH() {
		return nil
	}
	user := GetUser()
	token := getToken()
	if token != "" {
		return &http.BasicAuth{
			Username: user,
			Password: token,
		}
	}
	return &http.BasicAuth{
		Username: user,
		Password: askPassword(),
	}
}

func IsSSH() bool {
	return viper.GetString("user") == "@ssh"
}

func GetMode() string {
	mode := viper.GetString("mode")
	return lo.Ternary(mode == "", "gitlab", mode)
}

func GetUser() string {
	user := viper.GetString("user")
	return lo.Ternary(user == "", "@token", user)
}

func getToken() string {
	mode := GetMode()
	return viper.GetString(mode + ".token")
}
