# Connect to SQL Database with **sqlx**

จากบทความที่เรื่องการเชื่อมต่อฐานข้อมูลโดยใช้ standard library ที่ Go เตรียมไว้นั้น จะเห็นว่ามีหลายๆ จุด ยังใช้งานไม่ค่อยสะดวกเท่าไหร่ เช่น เรื่องการคิวรี่ข้อมูล แล้วต้องเอาตัวแปรไปรับค่าที่ละตัว และต้องใส่ลำดับให้ถูกต้องด้วยนั้น สามารถใช้ library [sqlx](https://github.com/jmoiron/sqlx) มาช่วยแก้ปัญหาเหล่านั้นได้ ซึ่งมี concept หลักๆ คือ

- Marshal rows into structs (with embedded struct support), maps, and slices
- Named parameter support including prepared statements
- `Get` and `Select` to go quickly from query to struct/slice

## ความรู้พื้นฐาน

1. ความรู้พื้นฐานภาษา Go
2. ความรู้เรื่องภาษา SQL

## มาเริ่มกันเลย

สร้างโปรเจคใหม่ และเปิดใช้งาน Go Module

```bash
mkdir -p godb/sqlx
cd godb/sqlx
go mod init godb/sqlx
```

สร้างไฟล์ `main.go`

```go
// godb/sqlx/main.go
package main

func main() {

}
```

ในบทความนี้จะใช้ระบบฐานข้อมูลเป็น [PostgreSQL](https://www.postgresql.org/) เริ่มจากการติดตั้ง database driver https://github.com/lib/pq และ library [sqlx](https://github.com/jmoiron/sqlx)

```bash
go get github.com/lib/pq
go get github.com/jmoiron/sqlx
```

<aside>
💡 Database driver อื่นๆ ดูที่ [https://go.dev/s/sqldrivers](https://go.dev/s/sqldrivers)

</aside>

การใช้งาน [sqlx](https://github.com/jmoiron/sqlx) เบื้องต้นสามารถใช้งานได้เหมือน `database/sql` เลย แค่เปลี่ยนจาก `sql` เป็น `sqlx` ก็สามารถใช้งานได้เลย โดยไม่ต้องแก้ไขโค้ดส่วนอื่นๆ

<aside>
💡 ใช้ `sqlx.DB` แทน `sql.DB`

ใช้ `sqlx.Tx` แทน `sql.Tx`

ใช้ `sqlx.Stmt` แทน `sql.Stmt` หรือจะใช้ `sqlx.NamedStmt` เพื่อใช้ named parameters

ใช้ `sqlx.Rows` แทน `sql.Row` ซึ่งจะได้จากการ return ของ `Qeuryx`

ใช้ `sqlx.Row` แทน `sql.Row` ซึ่งจะได้จากการ return ของ `QeuryRowx`

</aside>

## เชื่อมต่อ Database

การใช้งานจะเหมือนกัน standard library `database/sql` แค่เปลี่ยนจาก `sql.Open()` เป็น `sqlx.Open()`

```go
package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	// TODO fill this in directly or through environment variable
	// Build a DSN e.g. postgres://username:password@host:port/dbName
	DB_DSN = "postgres://fcricryh:F5a7wATfocTUNww1Dm14AfebtPaysqIn@john.db.elephantsql.com/fcricryh"
)

var db *sqlx.DB

func main() {
  connectDb()
	defer db.Close()

}

func connectDb() {
	var err error
	db, err = sqlx.Open("postgres", DB_DSN)
  if err != nil {
		log.Fatal("Cannot open DB connection", err)
  }

 // force a connection and test that it worked
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB Connected")
}
```

ตัว [sqlx](https://github.com/jmoiron/sqlx) ได้เตรียม `sqlx.Connect()` ไว้ให้โดยทำงานทั้งเชื่อมต่อ และตรวจสอบโดยใช้ ping ให้เรียบร้อยแล้ว

```go
func connectDb() {
	var err error
	db, err = sqlx.Connect("postgres", DB_DSN)
  if err != nil {
		log.Fatal("Cannot open DB connection", err)
  }

	log.Println("DB Connected")
}
```

## การค้นหาข้อมูลจาก Database

การค้นหาข้อมูลสามารถใช้ `db.Query()` ได้เหมือน `database/sql` หรือจะเปลี่ยนเป็น `db.Queryx()` ของ [sqlx](https://github.com/jmoiron/sqlx) ก็ได้ แต่ใช้ `db.Select()` ดีกว่า เพราะจะช่วยเอาข้อมูลใส่ Slice ให้เลย

สร้าง `struct{}` ขึ้นมา เพื่อเอาไปรับค่า ถ้าชื่อไม่ตรงกับในฐานข้อมูลให้ใช้ `db:"column_name"`

```go
type Todo struct {
	Id     int
	Title  string
	IsDone bool `db:"is_done"`
}
```

ค้นหาโดยใช้ `db.Select()`

```go
func main() {
  connectDb()
	defer db.Close()

	todos, err := GetTodos()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(todos)
}

func GetTodos() ([]Todo, error) {
	// Query the database, storing results in a []Todo (wrapped in []interface{})
	todos := []Todo{}
	err := db.Select(&todos, "SELECT * FROM todos")
	if err != nil {
		return nil, err
	}
	return todos, nil
}
```

### **การค้นหาแบบมีเงื่อนไข**

ถ้าเราต้องการเพิ่มเงื่อนไขในการค้นหา เช่น จะค้นหาข้อมูลตามสถานะที่ทำเสร็จแล้ว ก็จะเขียนภาษา sql ได้แบบนี้ `select * from todos where is_done = true` หรือเฉพาะที่ยังไม่เสร็จ `select * from todos where is_done = false` จะเห็นว่าค่า `true/false` เป็น parameter ที่สามารถเปลี่ยนแปลงได้ตามที่ต้องการ

ซึ่งการแทนค่า parameter ในของ database driver https://github.com/lib/pq นั้นจะใช้ `$n` ซึ่ง `n` คือ ตัวเลขลำดับ เริ่มต้นที่ 1

```go
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
```

### **การค้นหาแบบมีเงื่อนไขโดยใช้ named parameters**

ตัว database driver https://github.com/lib/pq นั้น ไม่สามารถใช้ named parameters ได้ แต่ตัว [sqlx](https://github.com/jmoiron/sqlx) มีความสามารถนี้เพิ่มมาให้ โดยใช้ฟังก์ชัน `db.NamedQuery()`

<aside>
💡 การแทนค่า named parameters สามารถใช้ได้ 2 แบบ
1. `map[string]interface{}` สามารถตั้งชื่อ parameters ยังไงก็ได้
2. `structs` ชื่อ parameters ต้องเหมือนในฐานข้อมูล

</aside>

```go
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
```

### **การค้นหาข้อมูลจาก Id**

การค้นหาจาก Id ถ้ามีข้อมูลจากได้ออกมาแค่ 1 row จะใช้ `db.QueryRow()` ได้เหมือนเดิม หรือจะเปลียนมาเป็น `db.QueryRowx()` ของ [sqlx](https://github.com/jmoiron/sqlx) ก็ได้ แต่ใช้ `db.Get()` ดีกว่าเพราะจะได้ค่าใส่ struct กลับมาให้เลย

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
  todo := Todo{}
	q := "select * from todos where id=$1"

	err := db.Get(&todo, q, id)

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

การเพิ่มข้อมูลสามารถใช้ `db.Exec()` และตรวจสอบว่าเพิ่มสำเร็จหรือไม่ด้วย `RowsAffected()` ถ้าสำเร็จจะได้ค่ามากกว่า `0` ได้เหมือน `database/sql`

แต่ตัว [sqlx](https://github.com/jmoiron/sqlx) มีฟังก์ชัน `db.NamedExec()` มาให้ ซึ่งจะเป็นแบบ named parameters

### **ใช้การแทนค่าด้วย `map[string]interface{}`**

```go
func main() {
	connectDb()
	defer db.Close()

	AddTodoWithMap("do something")

	todos, _ := GetTodos(nil)

	for _, todo := range todos {
		fmt.Println(todo)
	}
}

func AddTodoWithMap(title string) error {
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
```

### **ใช้การแทนค่าด้วย `struct`**

```go
func main() {
	connectDb()
	defer db.Close()

	todo := Todo{Title: "do something struct"}
	AddTodoWithStruct(todo)

	todos, _ := GetTodos(nil)

	for _, todo := range todos {
		fmt.Println(todo)
	}
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
```

หรือจะใช้วิธีการสร้าง statment ด้วย `db.Prepare()` หรือ `db.Preparex()` แล้วส่งค่า parameter ใน `Exec()` แทน

```go
func AddTodo(todo Todo) error {
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
	stmt, err := db.Preparex("INSERT INTO todos (title) VALUES ($1) returning id")
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

ถ้าต้องการใช้ named parameters ก็ใช้ `db.PrepareNamed()`

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
```

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
	stmt, err := db.Preparex("UPDATE todos SET is_done=$1 where id=$2")
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
	stmt, err := db.Preparex("DELETE FROM todos where id=$1")
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

ซึ่งสามารถทำได้เหมือน `database/sql` เลย แต่แนะนำใช้ `db.Beginx()` แทน ดีกว่า เพราะจะได้ใช้ `Get()` หรือ `Select()` ที่อยู่ใน `tx.Prepraex()`

```go
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
```

---

ก็จบแล้วสำหรับการทำ Select, Insert, Update และ Delete ลงฐานข้อมูลโดยใช้ [sqlx](https://github.com/jmoiron/sqlx) แทน standard library สามารถศึกษาเพิ่มเติมได้จาก [http://jmoiron.github.io/sqlx/](http://jmoiron.github.io/sqlx/)
