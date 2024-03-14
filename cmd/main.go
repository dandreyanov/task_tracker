package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"task_tracker/config"
	"task_tracker/internal/endpoints"
	"task_tracker/internal/handlers"
	"task_tracker/internal/tools"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	tasks := handlers.NewTaskRoutes()
	err = tools.LoadTasksFromFile(tasks.TaskResponse)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()

	endpoints.InitEndpoints(r, tasks)

	err = r.Run(viper.GetString("http.port"))
	if err != nil {
		return
	}
}
