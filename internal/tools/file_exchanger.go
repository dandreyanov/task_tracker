package tools

import (
	"encoding/json"
	"os"
	"task_tracker/internal/entity"
)

func SaveData(tasks []entity.Task) error {
	jsonData, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile("tasks.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadTasksFromFile(tr *entity.TasksDTO) error {

	data, err := os.ReadFile("tasks.json")
	if os.IsNotExist(err) {
		_, err = os.Create("tasks.json")
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &tr.Tasks)
	if err != nil {
		return err
	}
	tr.Total = len(tr.Tasks)
	return nil

}
