package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	// TODO fill this in directly or through environment variable
	// Build a DSN e.g. postgres://username:password@host:port/dbName
	DB_DSN = "postgres://fcricryh:F5a7wATfocTUNww1Dm14AfebtPaysqIn@john.db.elephantsql.com/fcricryh"
)

type Todo struct {
	Id     int
	Title  string
	IsDone bool
}

var db *sql.DB

func main() {
  connectDb()
	defer db.Close()

	todo := Todo{Title: "do something"}
	err := AddTodoTx(&todo)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("New", todo)

	err = DeleteTodo(todo.Id)
	if err != nil {
		log.Println(err)
		return
	}
	
	_, err = GetTodo(todo.Id)
	if err != nil {
		log.Println("Not found Id:", todo.Id)
		return
	}
}

func connectDb() {
	var err error
	db, err = sql.Open("postgres", DB_DSN)
  if err != nil {
		log.Fatal("Cannot open DB connection", err)
  }	

  // Test that the DSN is valid before using it
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB Connected")
}

func GetTodos(status *bool) ([]Todo, error) {
	var rows *sql.Rows
	var err error
	if status != nil {
		q := "select * from todos where is_done=$1"
		rows, err = db.Query(q, *status)
	} else {
		q := "select * from todos"
		rows, err = db.Query(q)
	}
	
	if err != nil {
		return nil, err
	}
	// ปิดเมื่อจบการทำงาน
	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		todo := Todo{}
		err = rows.Scan(&todo.Id, &todo.Title, &todo.IsDone)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func GetTodo(id int) (*Todo, error) {
	q := "select * from todos where id=$1"
	row := db.QueryRow(q, id)

	todo := Todo{}
	err := row.Scan(&todo.Id, &todo.Title, &todo.IsDone)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func AddTodo(todo *Todo) error {
	stmt, err := db.Prepare("INSERT INTO todos (title) VALUES ($1) returning id")
	if err != nil {
		return err
	}
	defer stmt.Close()

  // เปลี่ยน Exec() เป็น QueryRow()
	err = stmt.QueryRow(todo.Title).Scan(&todo.Id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTodo(todo Todo) error {
	stmt, err := db.Prepare("UPDATE todos SET is_done=$1 where id=$2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(todo.IsDone, todo.Id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot update")
	}

	return nil
}

func DeleteTodo(id int) error {
	stmt, err := db.Prepare("DELETE FROM todos where id=$1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot delete")
	}

	return nil
}

func AddTodoTx(todo *Todo) error {
	// เปิดใช้งาน transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// เปลี่ยน db เป็น tx
	stmt, err := tx.Prepare("INSERT INTO todos (title) VALUES ($1) returning id, is_done")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// เปลี่ยน Exec() เป็น QueryRow()
	err = stmt.QueryRow(todo.Title).Scan(&todo.Id, &todo.IsDone)
	if err != nil {
		// เมื่อ error ก็สั่ง rollback
		tx.Rollback()
		return err
	}
	// ทำงานสำเร็จก็สั่ง commit
	tx.Commit()
	return nil
}