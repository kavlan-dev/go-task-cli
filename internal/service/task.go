package service

import (
	"fmt"
	"go-task-cli/internal/model"
	"time"
)

type taskRepository interface {
	LoadTasks() ([]model.Task, error)
	SaveTasks(tasks []model.Task) error
}

type taskService struct {
	repo taskRepository
}

func NewTaskService(repo taskRepository) *taskService {
	return &taskService{repo: repo}
}

func (s *taskService) AddTask(desc string) (*model.Task, error) {
	tasks, err := s.repo.LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки задач: %w", err)
	}

	now := time.Now().Format(time.RFC3339)
	newTask := model.Task{
		Id:          s.nextId(),
		Description: desc,
		Status:      model.StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)

	if err := s.repo.SaveTasks(tasks); err != nil {
		return nil, fmt.Errorf("ошибка записи файла задач: %w", err)
	}

	return &newTask, nil
}

func (s *taskService) UpdateTask(id int, desc string) error {
	tasks, err := s.repo.LoadTasks()
	if err != nil {
		return fmt.Errorf("ошибка загрузки задач: %w", err)
	}

	i, err := s.taskIndexById(id)
	if err != nil {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	now := time.Now().Format(time.RFC3339)
	task := tasks[i]
	task.Description = desc
	task.UpdatedAt = now
	tasks[i] = task

	err = s.repo.SaveTasks(tasks)
	if err != nil {
		return fmt.Errorf("ошибка записи файла задач: %w", err)
	}

	return nil
}

func (s *taskService) DeleteTask(id int) error {
	tasks, err := s.repo.LoadTasks()
	if err != nil {
		return fmt.Errorf("ошибка загрузки задач: %w", err)
	}

	i, err := s.taskIndexById(id)
	if err != nil {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	tasks = append(tasks[:i], tasks[i+1:]...)

	err = s.repo.SaveTasks(tasks)
	if err != nil {
		return fmt.Errorf("ошибка записи файла задач: %w", err)
	}

	return nil
}

func (s *taskService) MarkTask(id int, status model.TaskStatus) error {
	tasks, err := s.repo.LoadTasks()
	if err != nil {
		return fmt.Errorf("ошибка загрузки задач: %w", err)
	}

	i, err := s.taskIndexById(id)
	if err != nil {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	now := time.Now().Format(time.RFC3339)
	task := tasks[i]
	task.Status = status
	task.UpdatedAt = now
	tasks[i] = task

	err = s.repo.SaveTasks(tasks)
	if err != nil {
		return fmt.Errorf("ошибка записи файла задач: %w", err)
	}

	return nil
}

func (s *taskService) ListTasks(statusFilter model.TaskStatus) ([]model.Task, error) {
	tasks, err := s.repo.LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки задач: %w", err)
	}

	var filteredTasks []model.Task
	if statusFilter == "" {
		filteredTasks = tasks
	} else {
		for _, task := range tasks {
			if task.Status == statusFilter {
				filteredTasks = append(filteredTasks, task)
			}
		}
	}

	return filteredTasks, nil
}

func (s *taskService) taskIndexById(id int) (int, error) {
	tasks, err := s.repo.LoadTasks()
	if err != nil {
		return 0, err
	}

	for i, task := range tasks {
		if task.Id == id {
			return i, nil
		}
	}

	return 0, fmt.Errorf("Задача не найдена (ID: %d)", id)
}

func (s *taskService) nextId() int {
	tasks, err := s.repo.LoadTasks()
	if err != nil {
		return 1
	}

	if len(tasks) == 0 {
		return 1
	}

	id := 0
	for _, task := range tasks {
		if task.Id >= id {
			id = task.Id + 1
		}
	}

	return id
}
