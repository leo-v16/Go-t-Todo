package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
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
		task_string := "$TASK:" + task_desp + "$STATUS:" + task_status + "\n"
		task_string_list = append(task_string_list, task_string)
	}
	return task_string_list
}

func AddTasksToFile(file_path string, T []Task) {
	file, err := os.OpenFile(file_path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		println("Task File Could Not Be Opened To Add")
	}
	defer file.Close()

	task_string_list := ConvertTaskToString(T)
	for _, task_string := range task_string_list {
		file.WriteString(task_string)
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
	reader := bufio.NewReader(os.Stdin)
	command_line, _ := reader.ReadString('\n')
	command_line = strings.TrimSpace(command_line)

	pull_pattern := regexp.MustCompile(`"([^"]+)"|(\S+)`)
	matches := pull_pattern.FindAllStringSubmatch(command_line, -1)
	var command []string
	for _, match := range matches {
		if match[1] != "" {
			command = append(command, match[1])
		} else {
			command = append(command, match[2])
		}
	}
	return command
}

func DisplayCommand(command []string) {
	print("[")
	for idx, str := range command {
		if strings.Contains(str, " ") {
			fmt.Printf(`"%v"`, str)
		} else {
			print(str)
		}
		if idx != len(command)-1 {
			print(" ")
		}
	}
	print("]\n")
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

func DeleteTask(task_list []Task, task_no_list ...int) []Task {
	check := true
	task_no_idx := 0
	for idx := range task_list {
		if check && idx == task_no_list[task_no_idx] {
			task_no_idx++
			if task_no_idx == len(task_no_list) {
				check = false
			}
		}
		if idx+task_no_idx >= len(task_list) {
			break
		}
		task_list[idx] = task_list[idx+task_no_idx]
	}
	task_list = task_list[:len(task_list)-len(task_no_list)]
	return task_list
}

func CreateTask(task_list []Task, tasks ...Task) []Task {
	task_list = append(task_list, tasks...)
	return task_list
}

func main() {
	file_path := "./database/database.db"
	task_list := ExtractTaskFromFile(file_path)
	for {
		command := GetCommand(">>:")
		// if !ValidateCommand(command) {
		// 	return
		// }
		ClearScreen()
		DisplayCommand(command)
		head := command[0]
		switch head {
		case "exit":
			AddTasksToFile(file_path, task_list)
			return
		case "create":
			task_no_len := len(command) - 1
			task_new_list := make([]Task, task_no_len)
			for idx := range task_new_list {
				task_new_list[idx] = Task{command[idx+1], false}
			}
			task_list = CreateTask(task_list, task_new_list...)
		case "delete":
			task_no_len := len(command) - 1
			task_no_list := make([]int, task_no_len)
			for idx := range task_no_list {
				task_no_list[idx], _ = strconv.Atoi(command[idx+1])
			}
			task_list = DeleteTask(task_list, task_no_list...)
		case "list":
			ListTaskList(task_list)
		case "help":
			return
		}
	}
}
