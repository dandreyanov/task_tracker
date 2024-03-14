package service

import "task_tracker/internal/entity"

func GetTaskByID(id string, t entity.TasksDTO) (bool, int, entity.Task) {
	for idx, task := range t.Tasks {
		if task.ID == id {
			return true, idx, task
		}
	}
	return false, -1, entity.Task{}
}

func DeleteTaskByID(id string, t entity.TasksDTO) bool {
	ok, idx, _ := GetTaskByID(id, t)
	if ok {
		t.Tasks = append(t.Tasks[:idx], t.Tasks[idx+1:]...)
		return true
	}
	return false // task not found
}
