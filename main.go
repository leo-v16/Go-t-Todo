package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

var COMMANDS = [...]string{"exit", "create"}

const NO_OF_COMMANDS = len(COMMANDS)

type Task struct {
	task string // `json:"task"`
	done bool   // `json:"done"`
}

func NewTask(task string) *Task {
	return &Task{task, false}
}

func ExtractTaskFromLine(task_string string) Task {
	task_pattern := regexp.MustCompile(`\$TASK:([^$]+)`)
	done_pattern := regexp.MustCompile(`\$STATUS:([^$]+)`)
	task := task_pattern.FindStringSubmatch(task_string)[1]
	status := done_pattern.FindStringSubmatch(task_string)[1]
	done := false
	if status == "DONE" {
		done = true
	}
	return Task{task, done}
}

func ExtractTaskFromFile(file_path string) []Task {
	file, err := os.Open(file_path)
	if err != nil {
		print("Task File Could Not Be Opened")
		return []Task{}
	}
	defer file.Close()

	var task_list []Task

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		task := ExtractTaskFromLine(line)
		task_list = append(task_list, task)
	}
	return task_list
}

func main() {
	// for {
	// 	command := GetCommand("COMMAND:")
	// 	if !ValidateCommand(command) {
	// 		continue
	// 	}
	// 	head := command[0]
	// 	switch head {
	// 	case "exit":
	// 		return
	// 	}
	// }
	fmt.Print(ExtractTaskFromFile("database.db"))
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
