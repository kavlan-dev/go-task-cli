package app

import (
	"fmt"
	"go-task-cli/internal/model"
	"os"
	"strconv"
	"strings"
)

type TaskService interface {
	AddTask(description string) (*model.Task, error)
	UpdateTask(id int, description string) error
	DeleteTask(id int) error
	MarkTask(id int, status model.TaskStatus) error
	ListTasks(status model.TaskStatus) ([]model.Task, error)
}

func Run(serv TaskService) {
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

		task, err := serv.AddTask(strings.Join(args, " "))
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
		fmt.Printf("Задача добавлена успешно (ID: %d)\n", task.Id)
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

		err = serv.UpdateTask(id, strings.Join(args[1:], " "))
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
		fmt.Printf("Задача обновлена успешно (ID: %d)\n", id)
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

		err = serv.DeleteTask(id)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
		fmt.Printf("Задача удалена (ID: %d)\n", id)
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

		err = serv.MarkTask(id, model.StatusTodo)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
		fmt.Printf("Задача пометлена как TODO (ID: %d)\n", id)
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

		err = serv.MarkTask(id, model.StatusInProgress)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
		fmt.Printf("Задача пометлена как в процессе (ID: %d)\n", id)
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

		err = serv.MarkTask(id, model.StatusDone)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}
		fmt.Printf("Задача пометлена как выполненная (ID: %d)\n", id)
	case "list":
		var status model.TaskStatus
		if len(args) != 0 {
			status = model.TaskStatus(args[0])
		}

		tasks, err := serv.ListTasks(status)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("Задачи не найдены.")
			return
		}

		fmt.Println("Задачи:")
		for _, task := range tasks {
			fmt.Println("ID:", task.Id)
			fmt.Println("Описание:", task.Description)
			fmt.Println("Статус:", task.Status)
			fmt.Println("Создано:", task.CreatedAt)
			fmt.Println("Обновлено:", task.UpdatedAt)
			fmt.Println("-------------------")
		}
	default:
		fmt.Printf("Неверная команда: %s\n", command)
		return
	}
}
