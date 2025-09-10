package main

import (
	"database/sql"
	"fmt"
	"time"
	"strings"
	"strconv"
)

func add(db *sql.DB, text string) error {
	timestamp := time.Now().Unix()

	stmt, err := db.Prepare("INSERT INTO tasks (text, timestamp, done) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(text, timestamp, false)
	if err != nil {
		return err
	}

	fmt.Println("Task saved.")
	return nil
}


func delete(db *sql.DB, taskId int) error {
	stmt, err := db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(taskId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		fmt.Println(fmt.Sprintf("Task %d deleted", taskId))
	} else {
		fmt.Println("No task found with given id.")
	}

	return nil
}

func edit(db *sql.DB, taskId int, newText string) error {
	stmt, err := db.Prepare("UPDATE tasks SET text = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newText, taskId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		fmt.Println(fmt.Sprintf("Task %d edited to %v", taskId, newText))
	} else {
		fmt.Println("No task found with given id.")
	}

	return nil
}

func list(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return err
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.ID, &task.Text, &task.Timestamp, &task.Done)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	}

	longest_text_length := 4
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

	if len(tasks) > 0 {
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

			// calculate relative time
			current_timestamp := time.Now().Unix()
			passed_time := current_timestamp - task.Timestamp

			var time_text string
			if passed_time / 31536000 != 0 {
				time_text = fmt.Sprintf("%dy ago", passed_time / 31536000)
			} else if passed_time / 2592000 != 0 {
				time_text = fmt.Sprintf("%dm ago", passed_time / 2592000)
			} else if passed_time / 86400 != 0 {
				time_text = fmt.Sprintf("%dd ago", passed_time / 86400)
			} else if passed_time / 3600 != 0 {
				time_text = fmt.Sprintf("%dh ago", passed_time / 3600)
			} else if passed_time / 60 != 0 {
				time_text = fmt.Sprintf("%dm ago", passed_time / 60)
			} else if passed_time < 10 {
				time_text = fmt.Sprintf("just now")
			} else {
				time_text = fmt.Sprintf("%ds ago", passed_time)
			}



			fmt.Println(
				"|",
				task.ID,
				strings.Repeat(" ", longest_id_length-len(strconv.Itoa(task.ID))),
				"|",
				task.Text,
				strings.Repeat(arrow, longest_text_length-len(task.Text)),
				"|",
				time_text,
				strings.Repeat(" ", longest_timestamp_length-len(time_text)-1),
				"|",
				done)
			counter += 1
		}
	}else {
		all_characters_length := longest_id_length + 1 + longest_text_length + 1 + 10 + 1 + 2
		spaces := all_characters_length
		fmt.Println("|", strings.Repeat(" ", spaces/2), "No tasks", strings.Repeat(" ", spaces/2), " |")
	}


	printFooter(longest_id_length, longest_text_length, longest_timestamp_length, longest_done_length)

	return nil
}


func done(db *sql.DB, taskId int) error {
	stmt, err := db.Prepare("UPDATE tasks SET done=true WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(taskId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		fmt.Printf("Task with id: %d marked as done\n", taskId)
	} else {
		fmt.Println("No task found with given id.")
	}

	return nil
}


func undone(db *sql.DB, taskId int) error {
	stmt, err := db.Prepare("UPDATE tasks SET done=false WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(taskId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		fmt.Printf("Task with id: %d marked as not done\n", taskId)
	} else {
		fmt.Println("No task found with given id.")
	}

	return nil
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
