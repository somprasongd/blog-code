# Connect to SQL Database with GORM

ในการพัฒนาโปรแกรมติดต่อฐานข้อมูล หลายๆ คนอาจไม่ถนัดการใช้ภาษา SQL หรืออยากหาอะไรมาช่วยให้เขียนโค้ดสั้นลง หรือมาช่วยให้ทำงานง่ายขึ้น เร็วขึ้น ซึ่งในภาษา Go นั้นมี [GORM](https://github.com/go-gorm/gorm) ซึ่งเป็น ORM library มี feature ให้ใช้งานครบ ใช้งานง่าย และมีระบบ Auto Migrations มาให้ด้วย

## ความรู้พื้นฐาน

1. ความรู้พื้นฐานภาษา Go
2. ความรู้พื้นฐานภาษา SQL

## มาเริ่มกันเลย

สร้างโปรเจคใหม่ และเปิดใช้งาน Go Module

```bash
mkdir -p godb/gorm
cd godb/gorm
go mod init godb/gorm
```

สร้างไฟล์ `main.go`

```go
// godb/gorm/main.go
package main

func main() {

}
```

ในบทความนี้จะใช้ระบบฐานข้อมูลเป็น [PostgreSQL](https://www.postgresql.org/) ดังนั้นต้องติดตั้ง [GORM](https://github.com/go-gorm/gorm) และ database driver

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

<aside>
💡 GORM officially supports databases MySQL, PostgreSQL, SQLite, SQL Server

</aside>

## เชื่อมต่อ Database

การเชื่อมต่อ Database ต้องใช้ `gorm.Open()` ซึ่งต้องการ 2 อย่าง คือ `dialector` และ `config` ดังนั้นต้องสร้าง dialector ของ database ที่ใช้งานขึ้นมาก่อน

```go
package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// TODO fill this in directly or through environment variable
	// Build a DSN e.g. postgres://username:password@host:port/dbName
  // or "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	DB_DSN = "postgres://fcricryh:F5a7wATfocTUNww1Dm14AfebtPaysqIn@john.db.elephantsql.com/fcricryh"
)

var db *sqlx.DB

func main() {
  connectDb()
}

func connectDb() {
	var err error
	db, err = gorm.Open(postgres.Open(DB_DSN), &gorm.Config{})

	if err != nil {
		log.Fatal("Cannot open DB connection", err)
	}

	log.Println("DB Connected")
}
```

<aside>
💡 รูปแบบของ Data Source Name ใช้ได้ 2 แบบ
1. *`dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Bangkok"`
2. `dsn := "postgres://username:password@host:port/dbName"`*

</aside>

## สร้าง Model

เนื่องจาก GORM เป็น ORM ดังนั้นเราจะต้องสร้าง Model ซึ่งเป็น `struct` ที่มีโครงสร้างเหมือนกับตารางในฐานข้อมูลขึ้นมาก่อน

```go
type Test struct {
	ID     			 uint
	Name  			 string
	CreatedAt    time.Time
  UpdatedAt    time.Time
}
```

### **Migration**

GORM มีระบบ migrate มาให้ ซึ่งจะเอา Model ไปสร้างเป็น database schema ให้ และคอยอัพเดทให้อยู่ตลอดเวลา โดยใช้ `db.AutoMigrate`

```go
func main() {
	connectDb()

	db.AutoMigrate(&Test{})
}
```

เมื่อทดลองรันโปรแกรม GORM จะไปสร้างตารางชื่อ `tests` ขึ้นมาให้เลย

### **แสดงคำสั่ง SQL**

ถ้าต้องการให้แสดง sql ที่ GORM สร้างขึ้นมาให้ ให้เปลี่ยน logger level เป็น `Info`

```go
db, err = gorm.Open(postgres.Open(DB_DSN), &gorm.Config{
  Logger:logger.Default.LogMode(logger.Info),
})
```

เมื่อรันโปรแกรมใหม่อีกครั้งจะมีคำสั่ง sql แสดงขึ้นมาแล้ว

<aside>
💡 ถ้าต้องการรันดู sql อย่างเดียว แต่ไม่ต้องการให้ทำงานจริง ให้เพิ่ม option  `DryRun: true`

</aside>

### **Conventions**

จะตัวอย่างข้างบนจะเห็นว่า GORM จะมีการสร้างมีตารางชื่อ `tests` และมี `id` เป็น `primary key` ทั้งๆ ที่เราไม่ได้ตั้งค่าอะไรเลย เหตุที่เป็นอย่างนั้น เพราะว่า GORM จะใช้วิธี convention over configuration ทำให้ไม่ต้องไปตั้งค่าอะไรเลย แค่ทำให้ตรงตาม convention ของ GORM เท่านั้นก็พอ แต่ถ้าปรับแก้ไขค่ายังสามารถแก้ไขได้ ซึ่งมี Conventions ตามนี้

1. `ID` จะเป็น Primary Key
2. ชื่อตารางจะถูกสร้างเป็น pluralizes \*\*\*\*จาก struct เป็น `snake_cases` เช่น Test จะได้ตารางชื่อ tests หรือ Person จะได้ชื่อตารางเป็น people ถ้าหากต้องการได้ชื่อเป็น persons ต้อง implement `Tabler` inferface ดังนี้

   ```go
   type Person struct {
     ID uint
     name string
   }
   // TableName overrides the table name used by User to `profiles`
   func (Person) TableName() string {
     return "person"
   }
   ```

3. ชื่อคอลัมน์จะถูกตั้งเป็น `snake_case`
4. ถ้ามี Model มี field ชื่อ `CreatedAt` และถ้าตอน insert แล้วไม่กำหนดมาให้ GORM จะใส่ค่าให้เอง
5. ถ้ามี Model มี field ชื่อ `UpdatedAt` และถ้าตอน insert หรือ update แล้วไม่กำหนดมาให้ GORM จะใส่ค่าให้เอง
6. `gorm.Model` ถ้าใน Model ของเรามี `ID`, `CreatedAt`, `UpdatedAt` และ `DeletedAt` สามารถใส่เป็น `gorm.Model` แทนได้เลย

```go
type Test struct {
  gorm.Model
  Name string
}
// equals
type Test struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
  Name string
}
```

### **Fields Tag**

จากตัวอย่างด้านบน ตารางที่ถูกสร้างมา จะเห็นว่า `name` เป็น `text` ถ้าต้องการเปลี่ยนเป็น `varchar` และกำหนดเป็น `not null`สามารถทำได้โดยการใส่ tag

```go
type Test struct {
  ID        uint
  Name      string `gorm:"type:varchar(50);not null"`
}
// หรือจะใช้ size:50 ก็ได้
type Test struct {
  ID        uint
  Name      string `gorm:"size:50;not null"`
}
// หรือต้องการจะเปลี่ยนชื่อคอลัมน์ก็ได้
type Test struct {
  ID        uint
  Name      string `gorm:"column:myname;size:50;not null"`
}
```

<aside>
💡 ดูเพิ่มเติมได้ที่ [https://gorm.io/docs/models.html#Fields-Tags](https://gorm.io/docs/models.html#Fields-Tags)

</aside>

## การค้นหาข้อมูลจาก Database

การค้นหาข้อมูลสามารถใช้ `db.Find(&Model)`

โดยจะใช้ตารางเดิมที่มีอยู่แล้วซึ่งก็คือ todos เพื่อนำมาสร้าง model

```go
type Todo struct {
	ID     uint
	Title  string
	Completed bool `gorm:"column:is_done"` // กรณีที่ชื่อ field กับ column ไม่เหมือนกัน
}
```

```go
func main() {
	connectDb()

	todos, err := GetTodos()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(todos)
}

func GetTodos() ([]Todo, error) {
	todos := []Todo{}
	// return มาเป็น tx *gorm.DB
	tx := db.Find(&todos) // "SELECT * FROM todos"
	// ดึง error จาก tx.Error -> returns error or nil
	if tx.Error != nil {
		return nil, tx.Error
	}
  // tx.RowsAffected -> returns found records count, equals `len(todos)`
	if tx.RowsAffected == 0 {
		return nil, errors.New("no todos")
	}
	return todos, nil
}
```

### **การค้นหาข้อมูลแค่ 1 Row**

ถ้าต้องการค้นหาข้อมูลเพียงแค่ 1 row เท่านั้น สามารถใช้ `First`, `Take`, `Last` ซึ่งคำสั่ง SQL จะใส่ `LIMIT 1` ไว้ให้

```go
// Get the first record ordered by primary key
db.First(&todo)
// SELECT * FROM users ORDER BY id LIMIT 1;

// Get one record, no specified order
db.Take(&todo)
// SELECT * FROM users LIMIT 1;

// Get last record, ordered by primary key desc
db.Last(&todo)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

result := db.First(&todo)
result.RowsAffected // returns count of records found
result.Error        // returns error or nil

// กรณีที่ไม่มีข้อมูลจะได้ ErrRecordNotFound สามารถตรวจสอบได้จาก
errors.Is(result.Error, gorm.ErrRecordNotFound)
```

<aside>
💡 ถ้าไม่ต้องการให้เกิด ErrRecordNotFound ให้เปลี่ยนมาใช้ `db.Limit(1).Find(&todo)` แทน

</aside>

**การค้นหาข้อมูลจาก Id**

สามารถใช้ `First`แล้วส่ง `id` ที่ต้องการหาเข้าไปได้เลย

```go
func main() {
	connectDb()

	todo, err := GetTodo(1)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(todo)
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
```

<aside>
💡 ถ้า id เป็น string เช่น uuid ต้องใช้ `db.First(&todo, "id = ?", "1b74413f-f3b8-409f-ac47-e8c062e3472a")`

</aside>

## **การค้นหาแบบมีเงื่อนไข**

ถ้ามีเงื่อนไข ให้ใช้ `db.Where()` แล้วตามด้วย `First`, `Take`, `Last` หรือ `Find` ซึ่งจะรองรับทั้ง string, struct และ map

### **String Conditions**

```go
todo := Todo{}

// Get first matched record
db.Where("title = ?", "do somethig").First(&todo)
// SELECT * FROM todos WHERE title = 'do somethig' ORDER BY id LIMIT 1;

// Get all matched records
db.Where("title <> ?", "do somethig").Find(&todos)
// SELECT * FROM todos WHERE title <> 'do somethig';

// IN
db.Where("title IN ?", []string{"do somethig", "do somethig 2"}).Find(&todos)
// SELECT * FROM todos WHERE title IN ('do somethig','do somethig 2');

// LIKE
db.Where("title LIKE ?", "%something%").Find(&todos)
// SELECT * FROM todos WHERE title LIKE '%something%';

// AND
db.Where("title = ? AND is_done >= ?", "do somethig", true).Find(&todos)
// SELECT * FROM todos WHERE title = 'do somethig' AND is_done = true;

// Time
db.Where("updated_at > ?", lastWeek).Find(&todos)
// SELECT * FROM todos WHERE updated_at > '2000-01-01 00:00:00';

// BETWEEN
db.Where("created_at BETWEEN ? AND ?", lastWeek, today).Find(&todos)
// SELECT * FROM todos WHERE created_at BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';
```

### **Struct & Map Conditions**

```go
// Struct
db.Where(&Todo{Title: "do something", Completed: true}).First(&todo)
// SELECT * FROM todos WHERE title = "do something" AND is_done = true ORDER BY id LIMIT 1;

// Map
db.Where(map[string]interface{}{"title": "do something", "is_done": true}).Find(&todos)
// SELECT * FROM todos WHERE title = "do something" AND is_done = true;

// Slice of primary keys
db.Where([]int64{20, 21, 22}).Find(&todos)
// SELECT * FROM users WHERE id IN (20, 21, 22);
```

ข้อแตกต่างระหว่างการใช้ struct กับ map คือ ถ้าค่าใน struct เป็น zero value จะไม่ถูกนำมาใส่เป็นเงื่อนไข ต้องใช้ map แทน เช่น

```go
db.Where(&Todo{Title: ""}).Find(&todos)
// SELECT * FROM todos";

db.Where(map[string]interface{}{"Title": ""}).Find(&todos)
// SELECT * FROM todos WHERE title = "";
```

<aside>
💡 หรือจะใช้ Inline condition ใน method พวก First(), Find() ก็ได้ ใช้เหมือน Where() ได้เลย

</aside>

### **ตัวอย่างการค้นหาจากสถานะ**

```go
func main() {
	connectDb()

	wheres := map[string]interface{}{"is_done": true}
	todos, err := GetTodosWithStatus(wheres)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(todos)
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
```

## การเพิ่มข้อมูลลง Database

ถ้าใช้ GORM การสร้างข้อมูลใหม่ลง Database นั้นง่ายมาก เพียงแค่สร้างข้อมูลของ Model ขึ้นมา แล้วส่งเป็น pointer ไปยัง `db.Create()` และ GORM จะคืนค่า primary key กลับมาเลยด้วย

```go
func main() {
	connectDb()
	defer db.Close()

	todo := Todo{Title: "do something"}
	AddTodo(&todo)

  // returns inserted data's primary key
	fmt.Println(todo.ID)
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
```

## การแก้ไขข้อมูล

การแก้ไขข้อมูลทำได้ 2 แบบ คือ อัพเดททุกคอลัมน์ และอัพเดทแค่บางคอลัมน์

1. **อัพเดทค่าทุกคอลัมน์** โดยใช้ `Save` ซึ่งวิธีนี้จะต้องไปดึงข้อมูลมาก่อน 1 รอบ แล้วมาแก้ไขค่าที่ต้อง แต่ตอนอัพเดทจะอัพเดททุกคอลัมน์ เช่น

```go
func UpdateTodoStatus(id uint, completed bool) (*Todo, error) {
	todo := Todo{}
	tx := db.First(&todo, id)
  // SELECT * FROM "todos" WHERE "todos"."id" = 40 ORDER something',false) RETURNING "id" BY "todos"."id" LIMIT 1
	if tx.Error != nil {
		return nil, tx.Error
	}

	todo.Completed = completed
	tx = db.Save(&todo)
  // UPDATE "todos" SET "title"='do something',"is_done"=true WHERE "id" = 40
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &todo, nil
}
```

1. **อัพเดทแค่คอลัมน์เดียว**

```go
func UpdateTodoStatus(id uint, completed bool) (*Todo, error) {
	todo := Todo{ID: id}
  // Update with id
	tx := db.Model(&todo).Update("is_done", completed)
  // UPDATE "todos" SET "is_done"=true WHERE "id" = 44
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected <= 0 {
		return nil, errors.New("cannot update")
	}

	db.First(&todo)
  // SELECT * FROM "todos" WHERE "todos"."id" = 44 ORDER BY "todos"."id" LIMIT 1

	return &todo, nil
}

// Update with conditions
db.Model(&Todo{}).Where("title= ?", "do something").Update("is_done", completed)
// UPDATE todos SET "is_done"=true WHERE title="do something";

// Update with conditions and model value
db.Model(&todo).Where("title= ?", "do something").Update("is_done", completed)
// UPDATE todos SET "is_done"=true WHERE "id" = 44 and title="do something"
```

1. **อัพเดทหลายคอลัมน์** ด้วย `struct` หรือ `map[string]interface{}` แต่ถ้าใช้ `struct` จะไม่อัพเดทค่าที่เป็น zero-value ให้

```go
// Update attributes with `struct`, will only update non-zero fields
db.Model(&todo).Updates(Todo{Title: "change task", Completed: true})
// UPDATE todos SET title="change task", "is_done"=true WHERE "id" = 44;

// Update attributes with `map`
db.Model(&user).Updates(map[string]interface{}{"title": "change task" "is_done": true})
// UPDATE todos SET title="change task", "is_done"=true WHERE "id" = 44;
```

<aside>
💡 ถ้าต้องการใช้ `struct` และให้อัพเดท field ที่เป็น zero-value ด้วยต้องใช้คู่กับ `Select`

</aside>

```go
// Update attributes with `struct`, will only update non-zero fields
db.Model(&todo).Updates(Todo{Title: "change task", Completed: false})
// UPDATE todos SET title="change task", "is_done"=false WHERE "id" = 44;

// Select with Struct (select zero value fields)
db.Model(&todo).Select("Title", "Completed").Updates(Todo{Title: "change task", Completed: false})
// UPDATE todos SET title="change task", "is_done"=false WHERE "id" = 44;

// Select all fields (select all fields include zero value fields)
db.Model(&todo).Select("*").Updates(Todo{Title: "change task", Completed: false})
// UPDATE "todos" SET "id"=0,"title"='change task',"is_done"=false WHERE "id" = 44

// จะโดนอัพเดท id=0 ไปด้วย ต้องใส่ค่า id ให้ struct หรือไม่ก็ Omit("ID") ออกไป
db.Model(&todo).Select("*").Omit("ID").Updates(Todo{Title: "change task", Completed: false})
// UPDATE "todos" SET "title"='change task',"is_done"=false WHERE "id" = 44
```

## การลบข้อมูล

สามารถสั่งลบข้อมูลจาก `id` โดยใช้ `Delete` หรือจะใส่เงื่อนไขเพิ่มใช้คู่กับ `Where` ก็ได้

```go
func main() {
	connectDb()

	todo := Todo{Title: "do something"}
	err := AddTodo(&todo)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("New", todo)

	err = DeleteTodo(todo.ID)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = GetTodo(todo.ID)
	if err != nil {
		log.Println("Not found ID:", todo.ID)
		return
	}
}

func DeleteTodo(id uint) error {
	tx := db.Delete(&Todo{}, id)
  // DELETE FROM "todos" WHERE "todos"."id" = 50

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected <= 0 {
		return errors.New("cannot delete")
	}

	return nil
}

// Todo's ID is `10`
todo := Todo{ID: 10}
db.Delete(&todo)
// DELETE from todos where id = 10;

// Delete with additional conditions
db.Where("is_done = ?", false).Delete(&todo)
// DELETE from todos where id = 10 AND is_done = false;
```

## การใช้งาน Transaction

การใช้ Transaction ใน GORM แบบปกติใช้ `db.Transaction` โดยจะ `rollback` ให้เมื่อเกิด `error` และ `commit` ให้เมื่อจบฟังก์ชัน

```go
db.Transaction(func(tx *gorm.DB) error {
  // do some database operations in the transaction (use 'tx' from this point, not 'db')
  if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
    // return any error will rollback
    return err
  }

  if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
    return err
  }

  // return nil will commit the whole transaction
  return nil
})
```

ถ้าหาต้องการจัดการเองทำได้แบบนี้

```go
// begin a transaction
tx := db.Begin()

// do some database operations in the transaction (use 'tx' from this point, not 'db')
tx.Create(...)

// ...

// rollback the transaction in case of error
tx.Rollback()

// Or commit the transaction
tx.Commit()
```

---

ก็จบแล้วสำหรับการทำ CRUD ลงฐานข้อมูลโดยใช้ [GORM](https://github.com/go-gorm/gorm) จะเห็นว่าใช้งานง่าย มีฟีเจอร์ให้ใช้งานครบ ช่วยลดการเขียนโค้ดของเราได้เยอะเลย และสามารถศึกษาเพิ่มเติมได้จาก [https://gorm.io/docs/](https://gorm.io/docs/)
