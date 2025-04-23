package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

var COMMANDS = [...]string{"exit", "create", "list", "help"}

const NO_OF_COMMANDS = len(COMMANDS)

type Task struct {
	task string // `json:"task"`
	done bool   // `json:"done"`
}

func NewTask(task string) *Task {
	return &Task{task, false}
}

func (T Task) String() string {
	done := "NOT DONE"
	if T.done {
		done = "DONE"
	}
	return "[TASK: " + T.task + " | STATUS: " + done + "]"
}

func ExtractTaskFromLine(task_string string) Task {
	task_pattern := regexp.MustCompile(`\$TASK:([^$]+)`)
	done_pattern := regexp.MustCompile(`\$STATUS:([^$]+)`)
	task := task_pattern.FindStringSubmatch(task_string)[1]   // captures (.*?)
	status := done_pattern.FindStringSubmatch(task_string)[1] // captures (.*?)
	done := false
	if status == "DONE" {
		done = true
	}
	return Task{task, done}
}

func ExtractTaskFromFile(file_path string) []Task {
	file, err := os.Open(file_path)
	if err != nil {
		print("Task File Could Not Be Opened To Extract")
		return []Task{}
	}
	defer file.Close()

	var task_list []Task

	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		line := scanner.Text()
		task := ExtractTaskFromLine(line)
		task_list = append(task_list, task)
	}
	return task_list
}

func ConvertTaskToString(T []Task) []string {
	var task_string_list []string
	for _, task := range T {
		task_desp := task.task
		task_status := "NOT"
		if task.done {
			task_status = "DONE"
		}
		task_string := "\n$TASK:" + task_desp + "$STATUS:" + task_status
		task_string_list = append(task_string_list, task_string)
	}
	return task_string_list
}

func AddTasksToFile(file_path string, T []Task) {
	file, err := os.OpenFile(file_path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		println("Task File Could Not Be Opened To Add")
	}
	defer file.Close()

	task_string_list := ConvertTaskToString(T)
	for _, task_string := range task_string_list {
		file.WriteString(task_string)
	}
}

func main() {
	file_path := "./database/database.db"
	task_list := ExtractTaskFromFile(file_path)
	for {
		command := GetCommand(">>:")
		if !ValidateCommand(command) {
			return
		}
		ClearScreen()
		fmt.Println(command)
		head := command[0]
		switch head {
		case "exit":
			AddTasksToFile(file_path, task_list)
			return
		case "create":
			title := command[1]
			task_list = append(task_list, Task{title, false})
			return
		case "list":
			ListTaskList(task_list)
		case "help":
			return
		}
	}
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func Input[T any](ouput string) T {
	var input T
	print(ouput)
	fmt.Scan(&input)
	return input
}

func GetCommand(output string) []string {
	print(output)
	var command_str string
	fmt.Scan(&command_str)
	command := strings.Fields(command_str)
	return command
}

func ValidateCommand(command []string) bool {
	if len(command) < 1 {
		return false
	}
	if !slices.Contains(COMMANDS[:], command[0]) {
		return false
	}
	return true
}

func ListTaskList(task_list []Task) {
	for idx, task := range task_list {
		fmt.Println(idx, task)
	}
}
