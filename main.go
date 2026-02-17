package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type taskStatus string

const (
	statusTodo       taskStatus = "todo"
	statusInProgress taskStatus = "in-progress"
	statusDone       taskStatus = "done"
)

type Task struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	Status      taskStatus `json:"status"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
}

const tasksFile = "tasks.json"

func loadTasks() ([]*Task, error) {
	var tasks []*Task

	data, err := os.ReadFile(tasksFile)
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

func saveTasks(tasks []*Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации задач: %v", err)
	}

	err = os.WriteFile(tasksFile, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла задач: %v", err)
	}

	return nil
}

func taskById(tasks []*Task, id int) (*Task, error) {
	for _, task := range tasks {
		if task.Id == id {
			return task, nil
		}
	}

	return nil, fmt.Errorf("Задача не найдена (ID: %d)", id)
}

func nextId(tasks []*Task) int {
	id := 0

	for _, task := range tasks {
		if task.Id > id {
			id = task.Id
		}
	}

	id++
	return id
}

func addTask(desc string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		return
	}

	now := time.Now().Format(time.RFC3339)
	newTask := Task{
		Id:          nextId(tasks),
		Description: desc,
		Status:      statusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, &newTask)

	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf("Ошибка записи файла задач: %v\n", err)
		return
	}

	fmt.Printf("Задача добавлена успешно (ID: %d)\n", newTask.Id)
}

func updateTask(id int, desc string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		return
	}

	task, err := taskById(tasks, id)
	if err != nil {
		fmt.Printf("Задача с ID %d не найдена\n", id)
		return
	}

	now := time.Now().Format(time.RFC3339)
	task.Description = desc
	task.UpdatedAt = now
	tasks[id-1] = task

	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf("Ошибка записи файла задач: %v\n", err)
		return
	}

	fmt.Printf("Задача %d обновлена успешно\n", id)
}

func deleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		return
	}

	_, err = taskById(tasks, id)
	if err != nil {
		fmt.Printf("Задача с ID %d не найдена\n", id)
		return
	}

	tasks = append(tasks[:id-1], tasks[id:]...)

	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf("Ошибка записи файла задач: %v\n", err)
		return
	}

	fmt.Printf("Задача %d удалена успешно\n", id)
}

func markTask(id int, status taskStatus) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		return
	}

	task, err := taskById(tasks, id)
	if err != nil {
		fmt.Printf("Задача с ID %d не найдена\n", id)
		return
	}

	now := time.Now().Format(time.RFC3339)
	task.Status = status
	task.UpdatedAt = now
	tasks[id-1] = task

	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf("Ошибка записи файла задач: %v\n", err)
		return
	}

	fmt.Printf("Задача %d помечена как %s\n", id, status)
}

func listTasks(statusFilter string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		return
	}

	var filteredTasks []*Task
	if statusFilter == "" {
		filteredTasks = tasks
	} else {
		status := taskStatus(statusFilter)
		for _, task := range tasks {
			if task.Status == status {
				filteredTasks = append(filteredTasks, task)
			}
		}
	}

	if len(filteredTasks) == 0 {
		if statusFilter == "" {
			fmt.Println("Нет задач")
		} else {
			fmt.Printf("Нет задач с статусом '%s'\n", statusFilter)
		}
		return
	}

	fmt.Println("Задачи:")
	for _, task := range filteredTasks {
		fmt.Printf("ID: %d\n", task.Id)
		fmt.Printf("Description: %s\n", task.Description)
		fmt.Printf("Status: %s\n", task.Status)
		fmt.Printf("Created: %s\n", task.CreatedAt)
		fmt.Printf("Updated: %s\n", task.UpdatedAt)
		fmt.Println("---")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: task-cli <команда> [аргументы...]")
		fmt.Println("Команды:")
		fmt.Println("  add <описание> - Добавить новую задачу")
		fmt.Println("  update <id> <описание> - Обновить задачу")
		fmt.Println("  delete <id> - Удалить задачу")
		fmt.Println("  mark-in-progress <id> - Отметить задачу как в процессе")
		fmt.Println("  mark-done <id> - Отметить задачу как выполненной")
		fmt.Println("  list [статус] - Список всех задач или задач по статусу (todo, in-progress, done)")
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		if len(args) < 1 {
			fmt.Println("Использование: task-cli add <описание>")
			return
		}

		addTask(strings.Join(args, " "))
	case "update":
		if len(args) < 2 {
			fmt.Println("Использование: task-cli update <id> <описание>")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Неверный идентификатор задачи: %v\n", err)
			return
		}

		updateTask(id, strings.Join(args[1:], " "))
	case "delete":
		if len(args) != 1 {
			fmt.Println("Использование: task-cli delete <id>")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Неверный идентификатор задачи: %v\n", err)
			return
		}

		deleteTask(id)
	case "mark-todo":
		if len(args) != 1 {
			fmt.Println("Использование: task-cli mark-todo <id>")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Неверный идентификатор задачи: %v\n", err)
			return
		}

		markTask(id, "todo")
	case "mark-in-progress":
		if len(args) != 1 {
			fmt.Println("Использование: task-cli mark-in-progress <id>")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Неверный идентификатор задачи: %v\n", err)
			return
		}

		markTask(id, "in-progress")
	case "mark-done":
		if len(args) != 1 {
			fmt.Println("Использование: task-cli mark-done <id>")
			return
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Неверный идентификатор задачи: %v\n", err)
			return
		}

		markTask(id, "done")
	case "list":
		status := ""
		if len(args) != 0 {
			status = args[0]
		}

		listTasks(status)
	default:
		fmt.Printf("Неверная команда: %s\n", command)
		return
	}
}
