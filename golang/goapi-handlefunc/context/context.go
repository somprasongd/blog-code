package context

type MyContext interface {
	Bind(v interface{}) error                   // รับ struct pointer อะไรก็ได้เข้ามา bind จาก json boby
	BindQuery(interface{}) error                // รับ struct pointer อะไรก็ได้เข้ามา bind จาก query ทั้งหมด
	Query(k string) (string, bool)              // รับ key ที่ต้องการหาเข้ามา ถ้าไม่เจอจะคืนค่า false กลับไป
	DefaultQuery(k string, d string) string     // รับ key ที่ต้องการหาเข้ามา ถ้าไม่เจอให้คืนค่า default ที่กำหนดกลับไป
	Param(k string) string                      // รับ key ที่ต้องการหาเข้ามา
	Header(k string) string                     // รับ key ที่ต้องการหาเข้ามา
	RequestId() string                          // เอาไว้อ่านค่าจาก header ที่ใช้งานบ่อยๆ เช่น X-Request-Id
	ResponseError(httpstatus int, err string)   // ส่ง json error กลับไปตาม status ที่กำหนด
	ResponseJSON(httpstatus int, v interface{}) // ส่ง json จาก struct กลับไปตาม status ที่กำหนด
}
