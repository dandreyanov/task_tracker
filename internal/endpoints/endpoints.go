package endpoints

import (
	"github.com/gin-gonic/gin"
	"task_tracker/internal/handlers"
)

func InitEndpoints(r *gin.Engine, tr *handlers.TaskRoutes) {
	r.GET("/all", tr.GetAllTasks)
	r.POST("/task", tr.CreateTask)
	r.PUT("/task/:id", tr.UpdateTask)
	r.DELETE("/delete/:id", tr.DeleteTask)
	r.GET("/filter", tr.GetFilterTasks)
	r.GET("/list", tr.ListTasks)
}
