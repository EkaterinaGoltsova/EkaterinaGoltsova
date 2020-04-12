package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/EkaterinaGoltsova/sprint-starter/jira"
)

type result struct {
	Message string
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}

func handle(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		handleGet(writer)
		return
	}

	if request.Method == http.MethodPost {
		handlePost(writer, request)
		return
	}

	handleUnknownMethod(writer, request)
}

func handleGet(writer http.ResponseWriter) {
	template.Must(template.ParseFiles("html/form.html")).Execute(writer, nil)
}

func handlePost(writer http.ResponseWriter, request *http.Request) {
	sprintNumber, err := getSprintNumber(request)

	var config jira.Config
	err = config.InitConfig(jira.ConfigPath)
	if err != nil {
		fmt.Print(err.Error())
	}

	err = jira.StartSprint(config, sprintNumber)
	result := result{}
	if err != nil {
		result.Message = fmt.Sprintf("Не удалось получить номер спринта. Ошибка - %s", err.Error())
	} else {
		result.Message = fmt.Sprintf("Спринт c номером - %d успешно начат", sprintNumber)
	}

	template.Must(template.ParseFiles("html/result.html")).Execute(writer, result)
}

func getSprintNumber(request *http.Request) (int, error) {
	sprintString := request.FormValue("sprint_number")
	return strconv.Atoi(sprintString)
}

func handleUnknownMethod(writer http.ResponseWriter, request *http.Request) {
	result := result{
		Message: fmt.Sprintf("Http метод - %s не поддерживается", request.Method),
	}

	template.Must(template.ParseFiles("html/result.html")).Execute(writer, result)
}
