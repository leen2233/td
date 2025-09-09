package main

import (
	"database/sql"
	"fmt"
	"os"
	"errors"
	"strconv"
	"strings"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	commands := map[string]func(db *sql.DB, args []string) error {
		"add": func(db *sql.DB, args []string) error {
			if len(args) < 1 {
				return errors.New("add: missing 'text' argument")
			}

			text := strings.Join(args, " ")

			return add(db, text)
		},
		"delete": func(db *sql.DB, args []string) error {
			if len(args) < 1 {
				return errors.New("delete: missing 'taskId' argument")
			}
			taskId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New(fmt.Sprintf("delete: invalid 'taskId': %v", err))
			}
			return delete(db, taskId)
		},
		"edit": func(db *sql.DB, args []string) error {
			if len(args) < 2 {
				return errors.New("edit: missing 'taskId' or 'newText' arguments")
			}
			taskId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New(fmt.Sprintf("edit: invalid 'taskId': %v", err))
			}

			text := strings.Join(args[1:], " ")
			return edit(db, taskId, text)
		},
		"list": func(db *sql.DB, args []string) error {
			return list(db)
		},
		"done": func(db *sql.DB, args []string) error {
			if len(args) < 1 {
				return errors.New("delete: missing 'taskId' argument")
			}
			taskId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New(fmt.Sprintf("delete: invalid 'taskId': %v", err))
			}
			return done(db, taskId)
		},
		"undone": func(db *sql.DB, args []string) error {
			if len(args) < 1 {
				return errors.New("delete: missing 'taskId' argument")
			}
			taskId, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New(fmt.Sprintf("delete: invalid 'taskId': %v", err))
			}
			return undone(db, taskId)
		},
	}

	if len(os.Args) < 2 {
		fmt.Println("Error: no command provided")
		fmt.Println("Available commands: add, edit, delete, list")
		os.Exit(1)
	}

	// initialize sqlite3 database
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT,
			timestamp INTEGER,
			done BOOLEAN
		)
	`)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	if cmd, ok:= commands[command]; ok {
		err := cmd(db, args)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Error: unknown command '%s'\n", command)
		fmt.Println("Available commands: add, edit, delete, list")
		os.Exit(1)
	}
}
