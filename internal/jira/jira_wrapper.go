package jira

import (
	"fmt"
	"strconv"

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

type editBoardSubFilterRequestParams struct {
	Query string `json:"query"`
}

//StartSprint is a pipeline finction for start a sprint
func (service *Service) StartSprint(sprintNumber int) error {
	err := service.processIssuesWithSprintStatus(sprintNumber)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprint("Не удалось переместить задачи в беклог"))
	}

	err = service.processCurrentIssues(sprintNumber)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprint("Не удалось добавить лейбл нового спринта задачам из предыдущего спринта"))
	}

	err = service.editBoard(sprintNumber)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprint("Не удалось отредактировать доску"))
	}

	return nil
}

func (service *Service) processIssuesWithSprintStatus(sprintNumber int) error {
	jql := fmt.Sprintf("status=%d", service.Config.SprintStatusID)
	issues, err := service.searchIssues(jql)

	if err != nil {
		return err
	}

	if len(issues) == 0 {
		return nil
	}

	for _, issue := range issues {
		fmt.Printf("Обработка задачи: %s", issue.ID)
		err = service.addSprintLabel(&issue, sprintNumber)
		if err != nil {
			return err
		}
		err = service.moveIssueFromSprintStatusToBacklog(&issue)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *Service) searchIssues(jql string) (result []jira.Issue, err error) {
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

func (service *Service) addSprintLabel(issue *jira.Issue, sprintNumber int) error {
	data := make(map[string]interface{})
	data["labels"] = append(issue.Fields.Labels, service.Config.GetSprintLabel(sprintNumber))

	update := make(map[string]interface{})
	update["fields"] = data

	_, err := service.Client.Issue.UpdateIssue(issue.ID, update)

	return err
}

func (service *Service) moveIssueFromSprintStatusToBacklog(issue *jira.Issue) error {
	_, err := service.Client.Issue.DoTransition(issue.ID, strconv.Itoa(service.Config.BacklogTransitionID))
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Не удалось переместить задачу - %s", issue.ID))
	}

	issue, _, err = service.Client.Issue.Get(issue.ID, &jira.GetQueryOptions{})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Не удалось получить задачу - %s", issue.ID))
	}

	return nil
}

func (service *Service) processCurrentIssues(sprintNumber int) error {
	jql := fmt.Sprintf(
		"labels=%s&status IN (%s)",
		service.Config.GetPreviousSprintLabel(sprintNumber),
		service.Config.CurrentSprintStatuses,
	)
	issues, err := service.searchIssues(jql)

	if err != nil {
		return err
	}

	if len(issues) == 0 {
		return nil
	}

	fmt.Print(len(issues))
	for _, issue := range issues {
		fmt.Printf("Обработка задачи: %s", issue.ID)
		err = service.addSprintLabel(&issue, sprintNumber)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *Service) editBoard(sprintNumber int) error {
	err := service.editBoardTitle(sprintNumber)
	if err != nil {
		return err
	}

	err = service.editBoardFilter(sprintNumber)
	if err != nil {
		return err
	}

	return service.editBoardSubFilter(sprintNumber)
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

func (service *Service) editBoardSubFilter(sprintNumber int) error {
	params := editBoardSubFilterRequestParams{
		Query: service.Config.GetSubfilterQuery(sprintNumber),
	}

	req, _ := service.Client.NewRequest(
		"PUT",
		fmt.Sprintf(
			"rest/greenhopper/1.0/subqueries/%d/",
			service.Config.Board.ID,
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
