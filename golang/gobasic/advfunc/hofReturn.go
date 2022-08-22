package advfunc

import "fmt"

func LearnHOFReturn() {
	s := hofGreeting(newNameFunc("Ball"))

	fmt.Println(s)
}

func newNameFunc(name string) nameFunc {
	// สร้างตัวแปรแบบ first class function
	nameFn := func() string {
		return name // เปลี่ยนมาใช้ค่าจากที่ส่งมาแทน
	}

	return nameFn
}

// สร้าง type ใหม่ขึ้นมาเป็น function ที่ return string
type nameFunc func() string

// เปลี่ยนรับ parameter มาเป็น type ใหม่ที่สร้างมา
func hofGreeting(nameFn nameFunc) string {
	return fmt.Sprintf("Hello %s", nameFn())
}
