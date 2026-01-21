package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	Done        bool   `json:"done"`
	CreatedAt   string `json:"created_at"`
	CompletedAt string `json:"completed_at,omitempty"`
}

type TodoList struct {
	Tasks  []Task `json:"tasks"`
	NextId int    `json:"next_id"`
}

const maxTaskLength = 200
const tasksPath = "tasks.json"

func loadTasks() (*TodoList, error) {
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &TodoList{NextId: 1}, nil
		}

		return nil, err
	}

	var tl TodoList
	if err := json.Unmarshal(data, &tl); err != nil {
		return nil, err
	}

	return &tl, nil
}

func saveTasks(tl *TodoList) error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tasksPath, data, 0644)
}

func parseTaskId(strId string) (int, bool) {
	id, err := strconv.Atoi(strId)
	if err != nil {
		fmt.Println("Ошибка: не верный id")
		return 0, false
	}

	return id, true
}

func findTaskIndex(tl *TodoList, id int) int {
	for i := range tl.Tasks {
		if tl.Tasks[i].Id == id {
			return i
		}
	}

	return -1
}

func validateTask(tl *TodoList, task Task) error {
	if len(task.Content) > maxTaskLength {
		return fmt.Errorf("Ошибка: текст задачи не должен превышать %d символов\n", maxTaskLength)
	}

	if strings.TrimSpace(task.Content) == "" {
		return fmt.Errorf("Ошибка: новый текст задачи не может быть пустым")
	}

	for _, t := range tl.Tasks {
		if strings.EqualFold(t.Content, task.Content) {
			return fmt.Errorf("Ошибка: задача с таким заголовком уже существует")
		}
	}

	return nil
}

func listTasks(tl *TodoList) {
	if len(tl.Tasks) == 0 {
		fmt.Println("Список задач пуст")
		return
	}

	fmt.Println("Список задач:")
	for _, task := range tl.Tasks {
		status := " "
		if task.Done {
			status = "x"
		}

		fmt.Printf("%d [%s], %s (создана: %s)", task.Id, status, task.Content, task.CreatedAt)
		if task.Done && task.CompletedAt != "" {
			fmt.Printf(", выполнена: %s", task.CompletedAt)
		}

		fmt.Println()
	}
}

func addTask(tl *TodoList, content string) {
	task := Task{
		Id:        tl.NextId,
		Content:   content,
		Done:      false,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	err := validateTask(tl, task)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tl.Tasks = append(tl.Tasks, task)
	tl.NextId++
	fmt.Printf("Добавлена задача %d: %s\n", task.Id, content)
}

func toggleTask(tl *TodoList, strId string) {
	id, ok := parseTaskId(strId)
	if !ok {
		return
	}

	index := findTaskIndex(tl, id)
	if index == -1 {
		fmt.Println("Задача не найдена")
		return
	}

	tl.Tasks[index].Done = !tl.Tasks[index].Done
	var status string
	if tl.Tasks[index].Done {
		status = "выполнено"
		tl.Tasks[index].CompletedAt = time.Now().Format("2006-01-02 15:04:05")
	} else {
		status = "не выполнено"
		tl.Tasks[index].CompletedAt = ""
	}

	fmt.Printf("Задача #%d отмечена как %s\n", id, status)
}

func deleteTask(tl *TodoList, strId string) {
	id, ok := parseTaskId(strId)
	if !ok {
		return
	}

	index := findTaskIndex(tl, id)
	if index == -1 {
		fmt.Println("Задача не найдена")
		return
	}

	tl.Tasks = append(tl.Tasks[:index], tl.Tasks[index+1:]...)
	fmt.Printf("Задача #%d была удалена\n", id)
}

func clearAllTasks(tl *TodoList) {
	tl.Tasks = []Task{}
	tl.NextId = 1
	fmt.Println("Все задачи очищены")
}

func completeAllTasks(tl *TodoList) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	for i := range tl.Tasks {
		if !tl.Tasks[i].Done {
			tl.Tasks[i].Done = true
			tl.Tasks[i].CompletedAt = currentTime
		}
	}

	fmt.Println("Все задачи отмечены как выполненные")
}

func main() {
	listFlag := flag.Bool("list", false, "List all tasks")
	addFlag := flag.String("add", "", "Add a new task")
	toggleFlag := flag.String("toggle", "", "Toggle task status (provide task ID)")
	deleteFlag := flag.String("delete", "", "Delete a task (provide task ID)")
	clearFlag := flag.Bool("clear", false, "Clear all tasks")
	completeAllFlag := flag.Bool("complete-all", false, "Mark all tasks as complete")

	flag.Parse()

	tl, err := loadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		return
	}

	if *listFlag {
		listTasks(tl)
		return
	}

	if *addFlag != "" {
		content := *addFlag
		addTask(tl, content)
		if err := saveTasks(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err.Error())
			return
		}
		return
	}

	if *toggleFlag != "" {
		id := *toggleFlag
		toggleTask(tl, id)
		if err := saveTasks(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err.Error())
			return
		}
		return
	}

	if *deleteFlag != "" {
		id := *deleteFlag
		deleteTask(tl, id)
		if err := saveTasks(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err)
			return
		}
		return
	}

	if *clearFlag {
		clearAllTasks(tl)
		if err := saveTasks(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err)
			return
		}
		return
	}

	if *completeAllFlag {
		completeAllTasks(tl)
		if err := saveTasks(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err)
			return
		}
		return
	}

	flag.Usage()
}
