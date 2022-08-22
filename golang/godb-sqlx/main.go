package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
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
	IsDone bool `db:"is_done"`
}

var db *sqlx.DB

func main() {
	connectDb()
	defer db.Close()

	todo := Todo{Title: "do something"}
	err := AddTodoTx(&todo)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(todo)
}

func connectDb() {
	var err error
	db, err = sqlx.Connect("postgres", DB_DSN)
	if err != nil {
		log.Fatal("Cannot open DB connection", err)
	}

	log.Println("DB Connected")
}

func GetTodos(status *bool) ([]Todo, error) {
	// Query the database, storing results in a []Todo (wrapped in []interface{})
	todos := []Todo{}
	var err error
	if status != nil {
		err = db.Select(&todos, "SELECT * FROM todos where is_done=$1", *status)
	} else {
		err = db.Select(&todos, "SELECT * FROM todos")
	}

	if err != nil {
		return nil, err
	}
	return todos, nil
}

func GetTodosNQ(status bool) ([]Todo, error) {
	// ใช้ map[string]interface{}
	// rows, err := db.NamedQuery(`SELECT * FROM todos WHERE is_done=:status`, map[string]interface{}{"status": status})
	// หรือใช้แบบ struct
	rows, err := db.NamedQuery(`SELECT * FROM todos WHERE is_done=:is_done`, Todo{IsDone: status})

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		todo := Todo{}
		err = rows.StructScan(&todo)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func GetTodo(id int) (*Todo, error) {
	todo := Todo{}
	q := "select * from todos where id=$1"

	err := db.Get(&todo, q, id)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func AddTodoWithMap(title string) error {
	// use map
	personMap := map[string]interface{}{
		"title": title,
	}

	q := "INSERT INTO todos (title) VALUES (:title)"
	result, err := db.NamedExec(q, personMap)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot insert")
	}

	return nil
}

func AddTodoWithStruct(todo Todo) error {
	q := "INSERT INTO todos (title) VALUES (:title)"
	result, err := db.NamedExec(q, &todo)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot insert")
	}

	return nil
}

func AddTodoPreparex(todo Todo) error {
	stmt, err := db.Preparex("INSERT INTO todos (title) VALUES ($1)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(todo.Title)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		return errors.New("cannot insert")
	}

	return nil
}

func AddTodo(todo *Todo) error {
	// เพิ่ม returning id เข้าไป
	nstmt, err := db.PrepareNamed("INSERT INTO todos (title) VALUES (:title) returning id")
	if err != nil {
		return err
	}
	defer nstmt.Close()

	// ใช้ Get หรือ Select กรณีคิวรี่ได้หลาย row ได้เลย
	err = nstmt.Get(todo, *todo)
	if err != nil {
		return err
	}

	return nil
}

func AddTodoTx(todo *Todo) error {
	// เปิดใช้งาน transaction ใช้ Beginx เพื่อใช้ Preparex
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	// เปลี่ยน db เป็น tx
	stmt, err := tx.Preparex("INSERT INTO todos (title) VALUES ($1) returning id, is_done")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// ถ้าใช้ preparex จะใช้ Get กับ Select ได้
	err = stmt.Get(todo, todo.Title)
	if err != nil {
		// เมื่อ error ก็สั่ง rollback
		tx.Rollback()
		return err
	}
	// ทำงานสำเร็จก็สั่ง commit
	tx.Commit()
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
