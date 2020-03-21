package jira

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
)

const ConfigPath = "config/config.yaml"

type editBoardNameRequestParams struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type editBoardFilterRequestParams struct {
	ID    int    `json:"id"`
	Query string `json:"query"`
}

func StartSprint(config Config, sprintNumber int) error {
	jiraClient, err := createJiraClient(config)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprint("Не удалось создать клиента для jira"))
	}

	err = moveIssues(jiraClient, config, sprintNumber)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprint("Не удалось переместить задачи в беклог"))
	}

	err = editBoard(jiraClient, config, sprintNumber)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprint("Не удалось отредактировать доску"))
	}

	return nil
}

func createJiraClient(config Config) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: config.User.Name,
		Password: config.User.Password,
	}

	return jira.NewClient(tp.Client(), config.Host)
}

func moveIssues(jiraClient *jira.Client, config Config, sprintNumber int) error {
	issues, err := searchIssues(jiraClient, config.GetSprintLabel(sprintNumber))

	if err != nil {
		return err
	}

	if len(issues) == 0 {
		return nil
	}

	for _, issue := range issues {
		err = moveIssue(jiraClient, &issue, config.IssuesStatusMapping)
		if err != nil {
			return err
		}
	}

	return nil
}

func searchIssues(jiraClient *jira.Client, label string) (result []jira.Issue, err error) {
	jql := fmt.Sprintf("labels=%s", label)
	options := &jira.SearchOptions{
		StartAt:    0,
		MaxResults: 50,
	}

	for {
		issues, resp, err := jiraClient.Issue.Search(jql, options)
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

func moveIssue(jiraClient *jira.Client, issue *jira.Issue, issuesStatusMapping map[string]string) error {
	value, isset := issuesStatusMapping[issue.Fields.Status.ID]

	if !isset {
		return nil
	}

	fmt.Printf("Меняем status задачи %s с %s на %s", issue.ID, issue.Fields.Status.ID, value)
	_, err := jiraClient.Issue.DoTransition(issue.ID, value)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Не удалось переместить задачу - %s", issue.ID))
	}

	issue, _, err = jiraClient.Issue.Get(issue.ID, &jira.GetQueryOptions{})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Не удалось получить задачу - %s", issue.ID))
	}

	return moveIssue(jiraClient, issue, issuesStatusMapping)
}

func editBoard(jiraClient *jira.Client, config Config, sprintNumber int) error {
	err := editBoardTitle(jiraClient, config, sprintNumber)
	if err != nil {
		return err
	}

	return editBoardFilter(jiraClient, config, sprintNumber)
}

func editBoardTitle(jiraClient *jira.Client, config Config, sprintNumber int) error {
	params := editBoardNameRequestParams{
		ID:   config.Board.ID,
		Name: config.GetBoardName(sprintNumber),
	}

	req, _ := jiraClient.NewRequest(
		"PUT",
		"rest/greenhopper/1.0/rapidviewconfig/name",
		params,
	)
	_, err := jiraClient.Do(req, nil)
	return err
}

func editBoardFilter(jiraClient *jira.Client, config Config, sprintNumber int) error {
	params := editBoardFilterRequestParams{
		ID:    config.Swimline.ID,
		Query: config.GetSwimlineFilter(sprintNumber),
	}

	req, _ := jiraClient.NewRequest(
		"PUT",
		fmt.Sprintf("rest/greenhopper/1.0/swimlanes/%d/%d", config.Board.ID, config.Swimline.ID),
		params,
	)

	_, err := jiraClient.Do(req, nil)
	return err
}
