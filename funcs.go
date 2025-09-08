package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"errors"
	"strings"
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

	var longest_text_length int
	for _, task := range tasks {
		if len(task.Text) > longest_text_length {
			longest_text_length = len(task.Text)
		}
	}

	for _, task := range tasks {
		fmt.Println(task.ID, "|", task.Text, strings.Repeat(" ", longest_text_length-len(task.Text)), "|", task.Timestamp, "|", task.Done)
	}

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
