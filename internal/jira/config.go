package jira

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config struct
type Config struct {
	Host                  string
	User                  user
	Board                 board
	Swimline              swimline
	Subfilter             subfilter
	LabelTemplate         string `mapstructure:"label_template"`
	SprintStatusID        int    `mapstructure:"sprint_status_id"`
	BacklogTransitionID   int    `mapstructure:"backlog_transition_id"`
	CurrentSprintStatuses string `mapstructure:"current_sprint_statuses"`
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

type subfilter struct {
	QueryTemplate string `mapstructure:"query_template"`
}

const sprintNumberPlaceholder = "#SPRINT_NUMBER#"

//InitConfig method
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

//GetSprintLabel method
func (config *Config) GetSprintLabel(sprintNumber int) string {
	return replaceSprintNumberPlaceholder(config.LabelTemplate, sprintNumber)
}

//GetPreviousSprintLabel method
func (config *Config) GetPreviousSprintLabel(sprintNumber int) string {
	return replaceSprintNumberPlaceholder(config.LabelTemplate, sprintNumber-1)
}

//GetBoardName method
func (config *Config) GetBoardName(sprintNumber int) string {
	return replaceSprintNumberPlaceholder(config.Board.NameTemplate, sprintNumber)
}

//GetSwimlineFilter method
func (config *Config) GetSwimlineFilter(sprintNumber int) string {
	return replaceSprintNumberPlaceholder(config.Swimline.FilterTemplate, sprintNumber)
}

//GetSubfilterQuery method
func (config *Config) GetSubfilterQuery(sprintNumber int) string {
	return replaceSprintNumberPlaceholder(config.Subfilter.QueryTemplate, sprintNumber)
}

func replaceSprintNumberPlaceholder(template string, sprintNumber int) string {
	return strings.Replace(template, sprintNumberPlaceholder, strconv.Itoa(sprintNumber), -1)
}
