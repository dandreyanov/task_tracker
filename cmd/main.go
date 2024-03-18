package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"log"
	"task_tracker/config"
	"task_tracker/internal/endpoints"
	"task_tracker/internal/handlers"
	"task_tracker/internal/tools"
)

var InitDB = `
CREATE TABLE IF NOT EXISTS tasks (
	id VARCHAR(36) PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
	description VARCHAR(200),
	status integer default 0,
	priority integer default 0);`

func main() {
	db, err := sql.Open("sqlite3", "taskTracker.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(InitDB)
	if err != nil {
		log.Fatal(err)
	}

	err = config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	tasks := handlers.NewTaskRoutes(db)
	err = tools.LoadTasksFromFile(&tasks.TaskResponse)
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
