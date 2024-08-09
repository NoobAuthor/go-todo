package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Task struct {
	ID        int
	Title     string
	Completed bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'list', 'remove', or 'complete' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo-cli add <task>")
			os.Exit(1)
		}
		addTask(os.Args[2])
	case "list":
		listTasks()
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo-cli remove <task-id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			os.Exit(1)
		}
		removeTask(id)
	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo-cli complete <task-id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid task ID")
			os.Exit(1)
		}
		completeTask(id)
	default:
		fmt.Println("Expected 'add', 'list', 'remove', or 'complete' subcommands")
		os.Exit(1)
	}
}

func addTask(title string) {
	tasks := readTasks()
	id := len(tasks) + 1
	newTask := Task{ID: id, Title: title, Completed: false}
	tasks = append(tasks, newTask)
	writeTasks(tasks)
	fmt.Println("Added task:", title)
}

func listTasks() {
	tasks := readTasks()
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for _, task := range tasks {
		status := " "
		if task.Completed {
			status = "x"
		}
		fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Title)
	}
}

func removeTask(id int) {
	tasks := readTasks()
	var updatedTasks []Task
	for _, task := range tasks {
		if task.ID != id {
			updatedTasks = append(updatedTasks, task)
		}
	}
	writeTasks(updatedTasks)
	fmt.Println("Removed task with ID:", id)
}

func completeTask(id int) {
	tasks := readTasks()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = true
		}
	}
	writeTasks(tasks)
	fmt.Println("Completed task with ID:", id)
}

func readTasks() []Task {
	var tasks []Task

	file, err := os.Open("tasks.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return tasks
		}
		fmt.Println("Error reading tasks:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var task Task
		err := json.Unmarshal([]byte(scanner.Text()), &task)
		if err != nil {
			fmt.Println("Error parsing task:", err)
			os.Exit(1)
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func writeTasks(tasks []Task) {
	file, err := os.Create("tasks.txt")
	if err != nil {
		fmt.Println("Error writing tasks:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, task := range tasks {
		data, err := json.Marshal(task)
		if err != nil {
			fmt.Println("Error encoding task:", err)
			os.Exit(1)
		}
		writer.WriteString(string(data) + "\n")
	}
	writer.Flush()
}
