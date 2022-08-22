package main

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// TODO fill this in directly or through environment variable
	// Build a DSN e.g. postgres://username:password@host:port/dbName
	// or "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	DB_DSN = "postgres://fcricryh:F5a7wATfocTUNww1Dm14AfebtPaysqIn@john.db.elephantsql.com/fcricryh"
)

type Todo struct {
	ID        uint
	Title     string
	Completed bool `gorm:"column:is_done"`
}

var db *gorm.DB

func main() {
	connectDb()

	TestTx()
}

func connectDb() {
	var err error
	db, err = gorm.Open(postgres.Open(DB_DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DryRun: false,
	})

	if err != nil {
		log.Fatal("Cannot open DB connection", err)
	}

	log.Println("DB Connected")
}

func GetTodos() ([]Todo, error) {
	todos := []Todo{}
	// return มาเป็น tx *gorm.DB
	tx := db.Find(&todos) // "SELECT * FROM todos"
	// ดึง error จาก tx
	if tx.Error != nil {
		return nil, tx.Error
	}
	// tx.RowsAffected returns found records count, equals `len(todos)`
	if tx.RowsAffected == 0 {
		return nil, errors.New("no todos")
	}
	return todos, nil
}

func GetTodosWithStatus(wheres map[string]interface{}) ([]Todo, error) {
	todos := []Todo{}

	// return มาเป็น tx *gorm.DB
	tx := db.Where(wheres).Find(&todos) // "SELECT * FROM todos where is_done = ?"
	// ดึง error จาก tx
	if tx.Error != nil {
		return nil, tx.Error
	}
	// tx.RowsAffected returns found records count, equals `len(todos)`
	if tx.RowsAffected == 0 {
		return nil, errors.New("no todos")
	}
	return todos, nil
}

func GetTodo(id uint) (*Todo, error) {
	todo := Todo{}
	// return มาเป็น tx *gorm.DB
	tx := db.First(&todo, id) // "SELECT * FROM todos where id = ?"
	// ดึง error จาก tx
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &todo, nil
}

func AddTodo(todo *Todo) error {
	result := db.Create(todo) // pass pointer of data to Create

	// result.Error -> returns error
	if result.Error != nil {
		return result.Error
	}

	// result.RowsAffected -> returns inserted records count
	if result.RowsAffected <= 0 {
		return errors.New("cannot insert")
	}

	return nil
}

func UpdateTodoStatus(id uint, completed bool) (*Todo, error) {
	todo := Todo{ID: id}
	tx := db.Model(&todo).Select("*").Omit("ID").Updates(Todo{Title: "change task", Completed: false})
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected <= 0 {
		return nil, errors.New("cannot update")
	}

	db.First(&todo)

	return &todo, nil
}

func DeleteTodo(id uint) error {
	tx := db.Delete(&Todo{}, id)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected <= 0 {
		return errors.New("cannot delete")
	}

	return nil
}

func TestTx() {
	todo := Todo{Title: "do something"}

	tx := db.Begin()

	tx.Create(&todo)

	id := todo.ID

	todo.Title = "Update in transaction"
	todo.Completed = true

	tx.Save(&todo)

	tx.Rollback()

	result := db.First(&todo, id)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	if result.RowsAffected <= 0 {
		fmt.Println("record not found")
		return
	}

	fmt.Println(todo)

}
