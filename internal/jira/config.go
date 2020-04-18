package jira

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Host                string
	User                user
	Board               board
	Swimline            swimline
	LabelTemplate       string            `mapstructure:"label_template"`
	IssuesStatusMapping map[string]string `mapstructure:"issues_status_mapping"`
}

type user struct {
	Name     string
	Password string
}

type board struct {
	ID           int
	NameTemplate string `mapstructure:"name_template"`
}

type swimline struct {
	ID             int
	FilterTemplate string `mapstructure:"filter_template"`
}

const sprintNumberPlaceholder = "#SPRINT_NUMBER#"

func (config *Config) InitConfig(path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Невозможно прочитать файл конфига: %s", path))
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Не удалось преобразовать конфиг в структуру: %s", path))
	}

	return nil
}

func (config *Config) GetSprintLabel(sprintNumber int) string {
	return replaceSprintNumberPlaceholder(config.LabelTemplate, sprintNumber)
}

func (config *Config) GetBoardName(sprintNumber int) string {
	return replaceSprintNumberPlaceholder(config.Board.NameTemplate, sprintNumber)
}

func (config *Config) GetSwimlineFilter(sprintNumber int) string {
	return replaceSprintNumberPlaceholder(config.Swimline.FilterTemplate, sprintNumber)
}

func replaceSprintNumberPlaceholder(template string, sprintNumber int) string {
	return strings.Replace(template, sprintNumberPlaceholder, strconv.Itoa(sprintNumber), -1)
}
