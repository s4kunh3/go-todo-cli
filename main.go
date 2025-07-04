package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      int
	Name    string
	Done    bool
	DueTime time.Time
}

var tasks []Task

/*
	Method I used for secure user input handling (begginer way)
	var, _ := reader.ReadString('\n') read until enter is pressed
	vat = strings.TrimSpace(var) gets rid of useless space in the input

	then depending on the type i handle the cases so no other input gets
	accepted and generates an unexpected error.

	If you (the reader) knows a better way pls tell me, i wanna know.
*/

// TODO: Add input validation
func createTask(reader *bufio.Reader) {
	fmt.Print("Enter task name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		return
	}
	name = strings.TrimSpace(name)
	if name == "" {
		fmt.Println("Task name cannot be empty.")
		return
	}
	fmt.Print("Enter due time (YYYY-MM-DD HH:MM) or press Enter to skip: ")
	dueDate, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading due time: ", err)
		return
	}
	dueDate = strings.TrimSpace(dueDate)

	var parsedDate time.Time
	if dueDate != "" {
		layout := "2006-01-02 15:04"
		parsedDate, err = time.Parse(layout, dueDate)
		if err != nil {
			fmt.Println("Invalid date format. Please use YYYY-MM-DD HH:MM.")
		}
	}
	tasks = append(tasks, Task{ID: getID(), Name: name, DueTime: parsedDate})
	fmt.Println("Task added!")
}

func showTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	fmt.Printf("%-4s | %-20s | %-10s | %-16s\n", "ID", "Task Name", "Status", "Due Time")
	fmt.Println("----------------------------------------------------------------")
	defer fmt.Println("----------------------------------------------------------------")
	for _, task := range tasks {
		status := "Incomplete"
		if task.Done {
			status = "Complete"
		}
		wrappedName := wrapText(task.Name, 20)
		for i, line := range wrappedName {
			if i == 0 {
				due := "-"
				if !task.DueTime.IsZero() {
					due = task.DueTime.Format("2006-01-02 15:04")
				}
				fmt.Printf("[%-1d]  | %-20s | %-10s | %-16s\n", task.ID, line, status, due)
			} else {
				fmt.Printf("     | %-20s |\n", line)
			}
		}
	}
}

func completeTask(reader *bufio.Reader) {
	fmt.Println("Enter task ID to mark as complete: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		return
	}
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("Task ID cannot be empty.")
		return
	}
	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			fmt.Println("Task marked as complete!")
			return
		}
	}
	fmt.Println("Task not found")
}

func deleteTask(reader *bufio.Reader) {
	fmt.Print("Enter task ID to delete: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		return
	}
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("Task ID cannot be empty.")
		return
	}
	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = slices.Delete(tasks, i, i+1)
			fmt.Println("Task deleted successfully.")
			return
		}
	}
	fmt.Println("Task not found.")
}

func saveTasks() error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile("tasks.json", data, 0644) // File, data yo be saved, permissions in octal
}

func loadTasks() error {
	file, err := os.ReadFile("tasks.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Add a function to edit tasks
func main() {
	err := loadTasks()
	if err != nil {
		fmt.Println("Could not load tasks:", err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Choose an option")
		fmt.Println("1: Create a task")
		fmt.Println("2: Show tasks")
		fmt.Println("3: Complete task")
		fmt.Println("4: Delete task")
		fmt.Println("5: Exit")
		fmt.Print("Option: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		option, err := strconv.Atoi(input)
		if err != nil || option < 1 || option > 5 {
			fmt.Println("Invalid input. Please choose an available option")
			continue
		}
		switch option {
		case 1:
			createTask(reader)
		case 2:
			showTasks()
		case 3:
			completeTask(reader)
		case 4:
			deleteTask(reader)
		case 5:
			defer fmt.Println("Exiting...")
			return
		}

		err = saveTasks()
		if err != nil {
			fmt.Println("Error saving tasks:", err)
		}
	}
}
