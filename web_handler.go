package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/EkaterinaGoltsova/sprint-starter/jira"

	"github.com/gin-gonic/gin"
)

func getForm(context *gin.Context) {
	context.HTML(http.StatusOK, "form.html", gin.H{})
}

func postForm(context *gin.Context) {
	sprintNumber, err := strconv.Atoi(context.PostForm("sprint_number"))
	if err != nil {
		context.JSON(http.StatusBadRequest, getErrorResponse(err))
		return
	}

	var config jira.Config
	err = config.InitConfig(jira.ConfigPath)
	if err != nil {
		context.JSON(http.StatusInternalServerError, getErrorResponse(err))
		return
	}

	err = jira.StartSprint(config, sprintNumber)
	if err != nil {
		context.JSON(http.StatusInternalServerError, getErrorResponse(err))
		return
	}

	var meta interface{}
	context.JSON(http.StatusOK, getSuccessResponse(
		fmt.Sprintf("Спринт c номером - %d успешно начат", sprintNumber),
		meta,
		err,
	))
}
