package jira

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
)

//Service struct
type Service struct {
	Client *jira.Client
	Config Config
}

type editBoardNameRequestParams struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type editBoardFilterRequestParams struct {
	ID    int    `json:"id"`
	Query string `json:"query"`
}

//StartSprint is a pipeline finction for start a sprint
func (service *Service) StartSprint(sprintNumber int) error {
	err := service.moveIssues(sprintNumber)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprint("Не удалось переместить задачи в беклог"))
	}

	err = service.editBoard(sprintNumber)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprint("Не удалось отредактировать доску"))
	}

	return nil
}

func (service *Service) moveIssues(sprintNumber int) error {
	issues, err := service.searchIssues(sprintNumber)

	if err != nil {
		return err
	}

	if len(issues) == 0 {
		return nil
	}

	for _, issue := range issues {
		err = service.moveIssue(&issue)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *Service) searchIssues(sprintNumber int) (result []jira.Issue, err error) {
	jql := fmt.Sprintf("labels=%s", service.Config.GetSprintLabel(sprintNumber))
	options := &jira.SearchOptions{
		StartAt:    0,
		MaxResults: 50,
	}

	for {
		issues, resp, err := service.Client.Issue.Search(jql, options)
		if err != nil {
			return result, errors.Wrapf(err, fmt.Sprint("Не удалось найти задачи для перемещения в беклог"))
		}

		result = append(result, issues...)
		if resp.StartAt+resp.MaxResults >= resp.Total {
			return result, nil
		}

		options.StartAt += resp.MaxResults
	}
}

func (service *Service) moveIssue(issue *jira.Issue) error {
	value, isset := service.Config.IssuesStatusMapping[issue.Fields.Status.ID]
	if !isset {
		return nil
	}

	fmt.Printf("Меняем status задачи %s с %s на %s", issue.ID, issue.Fields.Status.ID, value)
	_, err := service.Client.Issue.DoTransition(issue.ID, value)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Не удалось переместить задачу - %s", issue.ID))
	}

	issue, _, err = service.Client.Issue.Get(issue.ID, &jira.GetQueryOptions{})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Не удалось получить задачу - %s", issue.ID))
	}

	return service.moveIssue(issue)
}

func (service *Service) editBoard(sprintNumber int) error {
	err := service.editBoardTitle(sprintNumber)
	if err != nil {
		return err
	}

	return service.editBoardFilter(sprintNumber)
}

func (service *Service) editBoardTitle(sprintNumber int) error {
	params := editBoardNameRequestParams{
		ID:   service.Config.Board.ID,
		Name: service.Config.GetBoardName(sprintNumber),
	}

	req, _ := service.Client.NewRequest(
		"PUT",
		"rest/greenhopper/1.0/rapidviewconfig/name",
		params,
	)
	_, err := service.Client.Do(req, nil)
	return err
}

func (service *Service) editBoardFilter(sprintNumber int) error {
	params := editBoardFilterRequestParams{
		ID:    service.Config.Swimline.ID,
		Query: service.Config.GetSwimlineFilter(sprintNumber),
	}

	req, _ := service.Client.NewRequest(
		"PUT",
		fmt.Sprintf(
			"rest/greenhopper/1.0/swimlanes/%d/%d",
			service.Config.Board.ID,
			service.Config.Swimline.ID,
		),
		params,
	)

	_, err := service.Client.Do(req, nil)
	return err
}

//CreateJiraClient func return new Jira Client
func CreateJiraClient(config Config) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: config.User.Name,
		Password: config.User.Password,
	}

	return jira.NewClient(tp.Client(), config.Host)
}
