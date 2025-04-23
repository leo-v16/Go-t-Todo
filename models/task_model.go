package models

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
