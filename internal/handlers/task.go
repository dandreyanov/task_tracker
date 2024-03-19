package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"task_tracker/internal/entity"
)

type TaskRoutes struct {
	db           *sql.DB
	TaskResponse entity.TasksDTO
}

func NewTaskRoutes(database *sql.DB) *TaskRoutes {
	return &TaskRoutes{
		db:           database,
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
	_, err = t.db.Exec("INSERT INTO tasks (id, title, description, status, priority) VALUES ($1, $2, $3, $4, $5)", task.ID, task.Title, task.Description, task.Status, task.Priority)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task created"})
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	}

}

func (t *TaskRoutes) GetAllTasks(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=3600")

	rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
		return
	}
	var TaskResponse entity.TasksDTO
	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority)
		if err != nil {
			return
		}

		TaskResponse.Tasks = append(TaskResponse.Tasks, task)
	}
	TaskResponse.Total = len(TaskResponse.Tasks)
	c.JSON(http.StatusOK, TaskResponse)
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

	rows, err := t.db.Query("SELECT * FROM tasks WHERE status = $1 AND priority = $2", status, priority)
	if err != nil {
		return
	}
	var TaskResponse entity.TasksDTO
	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority)
		if err != nil {
			return
		}
		TaskResponse.Total = len(TaskResponse.Tasks)
		TaskResponse.Tasks = append(TaskResponse.Tasks, task)
	}
	c.JSON(http.StatusOK, TaskResponse)
}

func (t *TaskRoutes) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task entity.Task
	t.TaskResponse.Tasks = append(t.TaskResponse.Tasks, task)
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, errDb := t.db.Exec("UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4 WHERE id = $5", task.Title, task.Description, task.Status, task.Priority, id)
	rows, err := t.db.Query("SELECT * FROM tasks WHERE id = $1", id)
	err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority)
	c.JSON(http.StatusOK, task)
	if errDb != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": errDb.Error()})
	}
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": errDb.Error()})
	}
}

func (t *TaskRoutes) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	_, err := t.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusMultiStatus, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "task delete"})
	}
}

func (t *TaskRoutes) ListTasks(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	rows, err := t.db.Query("SELECT * FROM tasks LIMIT $1 OFFSET $2", 10, (page-1)*10)
	if err != nil {
		return
	}
	var TaskResponse entity.TasksDTO
	for rows.Next() {
		var task entity.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority)
		if err != nil {
			return
		}
		TaskResponse.Tasks = append(TaskResponse.Tasks, task)

		count := t.db.QueryRow("SELECT COUNT(*) FROM tasks")
		err := count.Scan(&TaskResponse.Total)
		if err != nil {
			return
		}
	}
	c.JSON(http.StatusOK, TaskResponse)
}
