package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/EkaterinaGoltsova/sprint-starter/internal/jira"

	"github.com/gin-gonic/gin"
)

//GetForm is Handler for method GET
func GetForm(context *gin.Context) {
	context.HTML(http.StatusOK, "form.html", gin.H{})
}

//PostForm is Handler for method POST
func PostForm(context *gin.Context) {
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
