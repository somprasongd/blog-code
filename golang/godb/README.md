# Connect to SQL Database with Go

ในภาษา Go ถ้าต้องการเขียนโปรแกรมต่อกับระบบฐานข้อมูล สามารถใช้ standard library ที่ Go เตรียมไว้ ให้ได้เลย ซึ่งก็คือ [database/sql](https://pkg.go.dev/database/sql)

## ความรู้พื้นฐาน

1. ความรู้พื้นฐานภาษา Go
2. ความรู้เรื่องภาษา SQL

## มาเริ่มกันเลย

สร้างโปรเจคใหม่ และเปิดใช้งาน Go Module

```bash
mkdir -p godb/sql
cd godb/sql
go mod init godb/sql
```

สร้างไฟล์ `main.go`

```go
// godb/sql/main.go
package main

func main() {

}
```

ในบทความนี้จะใช้ระบบฐานข้อมูลเป็น [PostgreSQL](https://www.postgresql.org/) ต้องติดตั้ง database driver ดังนี้

```bash
go get github.com/lib/pq
```

<aside>
💡 Database driver อื่นๆ ดูที่ [https://go.dev/s/sqldrivers](https://go.dev/s/sqldrivers)

</aside>

## เชื่อมต่อ Database

ในบทความนี้จะใช้ระบบฐานข้อมูลเป็น [PostgreSQL](https://www.postgresql.org/) และใช้ standard library ของ Go มาเชื่อมต่อผ่านฟังก์ชัน `sql.Open()`

```go
package main

import "database/sql"

func main() {
  // required
  var driverName, dataSourceName string
  // return มา 2 ตัว
	db, err := sql.Open(driverName, dataSourceName)

  if err != nil {
      panic(err)
  }
}
```

ซึ่งสิ่งที่ `func Open()` ต้องการ คือ driverName และ dataSourceName ดังนั้นเราจะต้องรู้ว่าฐานข้อมูลที่เราจะใช้งานนั้นใช้ driver name อะไร เช่น PostgreSQL จะใช้ https://github.com/lib/pq เป็น database driver และมี driver name เป็น `postgres`

และเชื่อมต่อฐานข้อมูล PostgreSQL จะมีรูปแบบของ dataSourceName คือ `postgres://username:password@host:port/dbName`

```go
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	// TODO fill this in directly or through environment variable
	// Build a DSN e.g. postgres://username:password@host:port/dbName
	DB_DSN = "postgres://fcricryh:F5a7wATfocTUNww1Dm14AfebtPaysqIn@john.db.elephantsql.com/fcricryh"
)

func main() {
  db, err := sql.Open("postgres", DB_DSN)
  if err != nil {
		log.Fatal("Cannot open DB connection", err)
  }
	defer db.Close()
}
```

ทดสอบการเชื่อมต่อใช้ `db.Ping()`

```go
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	// TODO fill this in directly or through environment variable
	// Build a DSN e.g. postgres://username:password@host:port/dbName
	DB_DSN = "postgres://fcricryh:F5a7wATfocTUNww1Dm14AfebtPaysqIn@john.db.elephantsql.com/fcricryh"
)

func main() {
  db, err := sql.Open("postgres", DB_DSN)
  if err != nil {
		log.Fatal("Cannot open DB connection", err)
  }
	defer db.Close()

  // force a connection and test that it worked
	err = db.Ping()
	if err != nil {
			panic(err)
	}

	log.Println("DB Connected")
}
```

ขอย้ายโค้ดการเชื่อมต่อไปไว้ที่ฟังก์ชัน `connectDb()`

```go
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	// TODO fill this in directly or through environment variable
	// Build a DSN e.g. postgres://username:password@host:port/dbName
	DB_DSN = "postgres://fcricryh:F5a7wATfocTUNww1Dm14AfebtPaysqIn@john.db.elephantsql.com/fcricryh"
)

var db *sql.DB

func main() {
  connectDb()
  // defer ต้องอยู่ที่ func main()
	defer db.Close()
}

func connectDb() {
	// เมื่อประกาศ db เป็น global ต้องประกาศ err แยกออกมา
	var err error
  // ใช้ = แทน :=
	db, err = sql.Open("postgres", DB_DSN)
  if err != nil {
		log.Fatal("Cannot open DB connection", err)
  }

  // ทดสอบการเชื่อมต่อ
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB Connected")
}
```

## การค้นหาข้อมูลจาก Database

การค้นหาข้อมูลจะใช้ฟังก์ชัน `db.Query()`

### **ตัวอย่างการใช้งาน**

```go
func main() {
  connectDb()
	defer db.Close()

	query := "select * from todos"
  rows, err := db.Query(query)
  if (err != nil) {
		// handle error
    panic(err)
  }
	// ปิดเมื่อจบการทำงาน
	defer rows.Close()

	// อ่านได้ครั้งละ 1 row return เป็น boolean ว่ามีหรือไม่
	for rows.Next() {
		// การอ่าน row ต้องเอาตัวแปรไปรับ
		id := 0
		title := ""
		is_done := false
		err = rows.Scan(&id, &title, &is_done)
		if err != nil {
			panic(err)
		}

		fmt.Println(id, title, is_done)
	}
}
```

<aside>
💡 ข้อควรระวัง ถ้าใช้ `select *` จะได้ลำดับของ column ตามในฐานข้อมูล ดังนั้นตอนเอาตัวแปรไปรับตำแหน่งต้องตรงกับในฐานข้อมูล

แนะนำให้ระบุชื่อ column ไปเลย `select id, title, is_done` ดีกว่า

</aside>

แต่เวลาใช้งานจริง เราจะสร้าง `struct{}` ขึ้นมา เพื่อเอาไปรับค่า และ return ค่ากลับไปให้ฟังก์ชันที่เรียกใช้งาน

```go
type Todo struct {
	Id     int
	Title  string
	IsDone bool
}
```

แก้การค้นหาเป็นฟังก์ชัน `func GetTodos() ([]Todo, error)` ซึ่งจะ `return []Todo` กลับออกไป และถ้ามี `error` จะไม่ได้จัดการเองฟังก์ชันนี้ โดยจะส่ง error ออกไปให้ฟังก์ชันที่เรียกมาเอาไปจัดการต่อเอง

```go
func main() {
  connectDb()
	defer db.Close()

	todos, err := GetTodos()
	if err != nil {
		log.Println(err)
    return
	}

	for _, todo := range todos {
		fmt.Println(todo)
	}
}

func GetTodos() ([]Todo, error) {
	q := "select * from todos"
	rows, err := db.Query(q)
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
```

### **การค้นหาแบบมีเงื่อนไข**

ถ้าเราต้องการเพิ่มเงื่อนไขในการค้นหา เช่น จะค้นหาข้อมูลตามสถานะที่ทำเสร็จแล้ว ก็จะเขียนภาษา sql ได้แบบนี้ `select * from todos where is_done = true` หรือเฉพาะที่ยังไม่เสร็จ `select * from todos where is_done = false` จะเห็นว่าค่า `true/false` เป็น parameter ที่สามารถเปลี่ยนแปลงได้ตามที่ต้องการ

ซึ่งการแทนค่า parameter ในของ database driver นี้ จะใช้ $n ซึ่ง n คือ ตัวเลขลำดับ เริ่มต้นที่ 1

```go
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
```

<aside>
💡 ในแต่ละ database driver จะใช้การแทนค่าตัวแปรแตกต่างกัน เช่นใน PostgreSQL ใช้ `$n` ส่วน MySQL ใช้ `?` หรือ MSSQL สามารถใช้ชื่อได้

</aside>

```go
// PostgreSQL
q := "select * from todos where is_done=$1"
rows, err = db.Query(q, *status)

// MySQL ต้องใส่ค่าให้ถูกลำดับด้วย
q := "select * from todos where is_done=?"
rows, err = db.Query(q, *status)

// MSSQL
q := "select * from todos where is_done=@status"
rows, err = db.Query(q, sql.Named("status", *status)
```

### **การค้นหาข้อมูลจาก Id**

การค้นหาจาก Id ถ้ามีข้อมูลจากได้ออกมาแค่ 1 row ดังนั้นจะใช้ `db.QueryRow()` แทน

```go
func main() {
	connectDB()
  defer db.Close()

	id := 1
	todo, err := GetTodo(id)

	if err != nil {
		fmt.Println("Not found id", id)
		return
	}

	fmt.Println(todo)
}

func GetTodo(id int) (*Todo, error) {
	q := "select * from todos where id=$1"
  // เปลี่ยนจาก Query มาเป็น QueryRow
	row := db.QueryRow(q, id)

	todo := Todo{}
	err := row.Scan(&todo.Id, &todo.Title, &todo.IsDone)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}
```

<aside>
💡 กรณีที่ `return Todo` ออกไป ถ้าต้องการให้ `return nil` ได้ จะต้องใช้ `*Todo`

</aside>

## การเพิ่มข้อมูลลง Database

การเพิ่มข้อมูลใช้ฟังก์ชัน `db.Exec()` และถ้าต้องการตรวจสอบว่าเพิ่มสำเร็จหรือไม่ให้ตรวจสอบจาก `RowsAffected()` ถ้าสำเร็จจะได้ค่ามากกว่า `0`

```go
func main() {
  connectDb()
	defer db.Close()

	todo := Todo{Title: "do something", IsDone: false}
	err := AddTodo(todo)
	if err != nil {
		log.Println(err)
		return
	}

	todos, err := GetTodos(nil)
	if err != nil {
		log.Println(err)
		return
	}

	for _, todo := range todos {
		fmt.Println(todo)
	}
}

func AddTodo(todo Todo) error {
	q := "INSERT INTO todos (title) VALUES ($1)"
	result, err := db.Exec(q, todo.Title, todo.IsDone)
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
```

หรือจะใช้วิธีการสร้าง statment ด้วย `Prepare()` แล้วส่งค่า parameter ใน `Exec()` แทน

```go
func AddTodo(todo Todo) error {
	stmt, err := db.Prepare("INSERT INTO todos (title) VALUES ($1)")
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
```

ถ้าต้องการดึงค่า id ของ row ที่เพิ่งสร้างไปตัว library ได้เตรียม `result.LastInsertId()` ไว้ให้แล้ว แต่ใน PostgreSQL ไม่รองรับ ซึ่งมีวิธีทำดังนี้

```go
func main() {
	connectDB()
	defer db.Close()

	todo := Todo{Title: "do something"}
	err := AddTodo(&todo)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(todo)
}

// รับเป็น pointer เพราะจะใส่ค่า id กลับไปให้
func AddTodo(todo *Todo) error {
  // เพิ่ม returning id เข้าไป
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
```

<aside>
💡 ต้องการมากว่า 1 ค่า ก็ได้ `stmt, err := db.Prepare("INSERT INTO todos (title) VALUES ($1) returning id, is_done")`

</aside>

## การแก้ไขข้อมูล

โค้ดก็จะเหมือนกันการเพิ่มข้อมูล

```go
func main() {
  connectDb()
	defer db.Close()

	todo := Todo{Title: "do something"}
	err := AddTodo(&todo)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("New", todo)

	todo.IsDone = true

	err = UpdateTodo(todo)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Updated", todo)
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
```

## การลบข้อมูล

```go
func main() {
  connectDb()
	defer db.Close()

	todo := Todo{Title: "do something"}
	err := AddTodo(&todo)
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
```

## การใช้งาน Transaction

ในบางครั้ง การบันทึกข้อมูลลงฐานข้อมูลไม่ได้จบใน statement เดียว อาจจะต้อง insert หลายๆ statement หรือทั้ง insert, update และ delete ตารางอื่นด้วย ในการทำ 1 business logic ถ้างานใดงานหนึ่งมีปัญหา และต้องการ rollback ทั้งหมด จะต้องใช้ transaction

```go
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
```

---

ก็จบแล้วสำหรับการทำ Select, Insert, Update และ Delete ลงฐานข้อมูลโดยใช้ standard library สิ่งสำคัญคือต้องใช้ database driver ให้ถูกต้องการระบบฐานข้อมูลที่เราใช้ และต้องไปศึกษาเพิ่มเติมว่าแต่ละ database driver มีวิธีการใช้งานที่แตกต่างกันอย่างไร
