package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"task_tracker/internal/entity"
	"task_tracker/internal/service"
	"task_tracker/internal/tools"
)

type TaskRoutes struct {
	TaskResponse entity.TasksDTO
}

func NewTaskRoutes() *TaskRoutes {
	return &TaskRoutes{
		TaskResponse: entity.TasksDTO{Tasks: make([]entity.Task, 0, 10), Total: 0},
	}
}

func (t *TaskRoutes) CreateTask(c *gin.Context) {
	var task entity.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}
	task.ID = uuid.New().String()
	t.TaskResponse.Tasks = append(t.TaskResponse.Tasks, task)
	c.JSON(http.StatusOK, gin.H{"message": "task created"})
	err = tools.SaveData(t.TaskResponse.Tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}

}

func (t *TaskRoutes) GetAllTasks(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, t.TaskResponse.Tasks)
}

func (t *TaskRoutes) GetFilterTasks(c *gin.Context) {
	statusStr, okStatus := c.GetQuery("status")
	priorityStr, okPriority := c.GetQuery("priority")
	if !okStatus && !okPriority {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not valid data"})
		return
	}
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	priority, err := strconv.Atoi(priorityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filterTasks := make([]entity.Task, 0, len(t.TaskResponse.Tasks))

	for _, task := range t.TaskResponse.Tasks {
		if (task.Status == status) && (task.Priority == uint8(priority)) {
			filterTasks = append(filterTasks, task)
		}
	}
	c.JSON(http.StatusOK, filterTasks)
}

func (t *TaskRoutes) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	ok, idx, task := service.GetTaskByID(id, t.TaskResponse)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t.TaskResponse.Tasks[idx] = task

	c.JSON(http.StatusOK, task)
	err = tools.SaveData(t.TaskResponse.Tasks)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}
}

func (t *TaskRoutes) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	ok := service.DeleteTaskByID(id, t.TaskResponse)
	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "task delete"})
		err := tools.SaveData(t.TaskResponse.Tasks)
		if err != nil {
			c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task not found"})
}

func (t *TaskRoutes) ListTasks(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	low := (page - 1) * 10
	high := page * 10
	if high > len(t.TaskResponse.Tasks) {
		high = len(t.TaskResponse.Tasks)
	}
	res := t.TaskResponse.Tasks[low:high]

	total := len(t.TaskResponse.Tasks)

	c.JSON(http.StatusOK, gin.H{"tasks": res, "total": total})
}
