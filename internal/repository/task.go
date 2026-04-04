package repository

import (
	"encoding/json"
	"fmt"
	"go-task-cli/internal/model"
	"os"
)

type taskRepository struct {
	tasksFile string
}

func NewTaskRepository(tasksFile string) *taskRepository {
	return &taskRepository{tasksFile: tasksFile}
}

func (r *taskRepository) LoadTasks() ([]model.Task, error) {
	var tasks []model.Task

	data, err := os.ReadFile(r.tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("ошибка загрузки задач: %v", err)
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга файла задач: %v", err)
	}

	return tasks, nil
}

func (r *taskRepository) SaveTasks(tasks []model.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации задач: %v", err)
	}

	err = os.WriteFile(r.tasksFile, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла задач: %v", err)
	}

	return nil
}
