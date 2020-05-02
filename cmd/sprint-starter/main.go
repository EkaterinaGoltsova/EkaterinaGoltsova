package main

import (
	"log"

	"github.com/EkaterinaGoltsova/sprint-starter/internal/jira"
	"github.com/EkaterinaGoltsova/sprint-starter/web"
	docopt "github.com/docopt/docopt-go"
	"github.com/gin-gonic/gin"
)

var usage = `
Usage:
    sprint-starter --config <path> --templates <path> --port <port>

Options:
	-h, --help      	Show this help.
	--config <path>   	Path to config file.
	--templates <path>  Path to templates files
	--port <port> 		Port for run
`

func main() {
	arguments, _ := docopt.ParseDoc(usage)

	var config jira.Config
	err := config.InitConfig(arguments["--config"].(string))
	if err != nil {
		log.Fatalln(err)
		return
	}

	jiraClient, err := jira.CreateJiraClient(config)
	if err != nil {
		log.Fatalln(err)
		return
	}

	router := gin.Default()
	router.LoadHTMLGlob(arguments["--templates"].(string) + "*")

	webHandler := web.Handler{
		JiraService: jira.Service{
			Client: jiraClient,
			Config: config,
		},
	}

	router.GET("/", webHandler.GetForm)
	router.POST("/", webHandler.PostForm)

	router.Run(":" + arguments["--port"].(string))
}
