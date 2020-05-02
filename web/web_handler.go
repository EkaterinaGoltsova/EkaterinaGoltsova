package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/EkaterinaGoltsova/sprint-starter/internal/jira"

	"github.com/gin-gonic/gin"
)

//Handler struct
type Handler struct {
	JiraService jira.Service
}

//GetForm is Handler for method GET
func (handler Handler) GetForm(context *gin.Context) {
	context.HTML(http.StatusOK, "form.html", gin.H{})
}

//PostForm is Handler for method POST
func (handler Handler) PostForm(context *gin.Context) {
	sprintNumber, err := strconv.Atoi(context.PostForm("sprint_number"))
	if err != nil {
		context.JSON(http.StatusBadRequest, getErrorResponse(err))
		return
	}

	err = handler.JiraService.StartSprint(sprintNumber)
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
