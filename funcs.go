package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"errors"
	"strings"
	"strconv"
)

func add(text string) error {
	tasks, err := getTasks()
	if err != nil {
		return err
	}
	latestTaskId := getLatestTaskId(tasks)
	timestamp := time.Now().Unix()

	task := Task{
		ID: latestTaskId + 1,
		Text: text,
		Timestamp: timestamp,
		Done: false,
	}

	tasks = append(tasks, task)
	err = saveTasks(tasks)
	if err != nil {
		return err
	}

	fmt.Println("Saved task")
	return nil
}


func delete(taskId int) error {
	tasks, err := getTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
			err = saveTasks(tasks)
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("Task %d deleted", taskId))
			return nil
		}
	}

	return nil
}

func edit(taskId int, newText string) error {
	tasks, err := getTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == taskId {
			tasks[i].Text = newText

			err = saveTasks(tasks)
			if err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("Task %d edited to %v", taskId, newText))
			return nil
		}
	}

	return nil
}

func list() error {
	tasks, err := getTasks()
	if err != nil {
		return err
	}

	longest_text_length := 2
	longest_id_length := 2
	longest_timestamp_length := 10
	longest_done_length := 4
	var counter int

	for _, task := range tasks {
		if len(task.Text) > longest_text_length {
			longest_text_length = len(task.Text)
		}
		id_length := len(strconv.Itoa(task.ID))
		if id_length > longest_id_length {
			longest_id_length = id_length
		}
	}

	printHeader(longest_id_length, longest_text_length, longest_timestamp_length, longest_done_length)

	for _, task := range tasks {
		var done, arrow string
		if counter % 2 == 1 {
			arrow = "."
		}else{
			arrow = " "
		}

		if task.Done {
			done = " ✅  |"
		}else{
			done = " ❌  |"
		}
		fmt.Println(
			"|",
			task.ID,
			strings.Repeat(" ", longest_id_length-len(strconv.Itoa(task.ID))),
			"|",
			task.Text,
			strings.Repeat(arrow, longest_text_length-len(task.Text)),
			"|",
			task.Timestamp,
			"|",
			done)
		counter += 1
	}

	printFooter(longest_id_length, longest_text_length, longest_timestamp_length, longest_done_length)

	return nil
}


func done(taskId int) error {
	tasks, err := getTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == taskId {
			tasks[i].Done = true
			err := saveTasks(tasks)
			if err != nil {
				return err
			}
			fmt.Printf("Task with id: %d marked as done\n", taskId)
			return nil
		}
	}
	return nil
}


func undone(taskId int) error {
	tasks, err := getTasks()
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == taskId {
			tasks[i].Done = false
			err := saveTasks(tasks)
			if err != nil {
				return err
			}
			fmt.Printf("Task with id: %d marked as undone\n", taskId)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Task with task id: %d not found", taskId))
}

func getLatestTaskId(tasks []Task) int {
	biggestId := 1
	for _, task := range tasks {
		if task.ID > biggestId {
			biggestId = task.ID
		}
	}
	return biggestId
}

func getTasks() ([]Task, error) {
	data, err := ioutil.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err){
			data = []byte("[]")
			err = ioutil.WriteFile("tasks.json", data, 0644)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}


func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("tasks.json", data, 0644)
}

func printHeader(longest_id_length, longest_text_length, longest_timestamp_length, longest_done_length int){
	fmt.Println(
		"┌",
		strings.Repeat("─", longest_id_length + 1),
		"┬",
		strings.Repeat("─", longest_text_length + 1),
		"┬",
		strings.Repeat("─", longest_timestamp_length),
		"┬",
		strings.Repeat("─", longest_done_length),
		"┐",
	)
	fmt.Println(
		"| ID",
		strings.Repeat(" ", longest_id_length - 2),
		"|",
		"Text",
		strings.Repeat(" ", longest_text_length - 4),
		"|",
		"Timestamp ",
		"|",
		"Done",
		"|",
	)
	fmt.Println(
		"├",
		strings.Repeat("─", longest_id_length + 1),
		"┼",
		strings.Repeat("─", longest_text_length + 1),
		"┼",
		strings.Repeat("─", longest_timestamp_length),
		"┼",
		strings.Repeat("─", longest_done_length),
		"┤",
	)
}


func printFooter(longest_id_length, longest_text_length, longest_timestamp_length, longest_done_length int) {
	fmt.Println(
		"└",
		strings.Repeat("─", longest_id_length + 1),
		"┴",
		strings.Repeat("─", longest_text_length + 1),
		"┴",
		strings.Repeat("─", longest_timestamp_length),
		"┴",
		strings.Repeat("─", longest_done_length),
		"┘",
	)
}
