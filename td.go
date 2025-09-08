package main

import (
	"fmt"
	"os"
	"errors"
	"strconv"
	"strings"
)

func main() {
	commands := map[string]func(args []string) error {
		"add": func(args []string) error {
			if len(args) < 1 {
				return errors.New("add: missing 'text' argument")
			}

			text := strings.Join(args, " ")

			return add(text)
		},
		"delete": func(args []string) error {
			if len(args) < 1 {
				return errors.New("delete: missing 'taskId' argument")
			}
			taskId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New(fmt.Sprintf("delete: invalid 'taskId': %v", err))
			}
			return delete(taskId)
		},
		"edit": func(args []string) error {
			if len(args) < 2 {
				return errors.New("edit: missing 'taskId' or 'newText' arguments")
			}
			taskId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New(fmt.Sprintf("edit: invalid 'taskId': %v", err))
			}

			text := strings.Join(args[1:], " ")
			return edit(taskId, text)
		},
		"list": func(args []string) error {
			return list()
		},
		"done": func(args []string) error {
			if len(args) < 1 {
				return errors.New("delete: missing 'taskId' argument")
			}
			taskId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New(fmt.Sprintf("delete: invalid 'taskId': %v", err))
			}
			return done(taskId)
		},
		"undone": func(args []string) error {
			if len(args) < 1 {
				return errors.New("delete: missing 'taskId' argument")
			}
			taskId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New(fmt.Sprintf("delete: invalid 'taskId': %v", err))
			}
			return undone(taskId)
		},
	}

	if len(os.Args) < 2 {
		fmt.Println("Error: no command provided")
		fmt.Println("Available commands: add, edit, delete, list")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	if cmd, ok:= commands[command]; ok {
		err := cmd(args)
		if err != nil {
			fmt.Println("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Error: unknown command '%s'\n", command)
		fmt.Println("Available commands: add, edit, delete, list")
		os.Exit(1)
	}
}
