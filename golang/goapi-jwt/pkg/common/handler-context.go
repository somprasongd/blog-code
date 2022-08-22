package common

type HandleFunc func(ctx HContext) error

type HContext interface {
	Method() string
	Path() string
	BodyParser(interface{}) error
	QueryParser(interface{}) error
	Query(string) (string, bool)
	DefaultQuery(string, string) string
	Param(string) string
	Header(string) string
	Authorization() string
	RequestId() string
	Locals(key string, value ...interface{}) interface{}
	Next() error
	SendStatus(int) error
	SendJSON(int, interface{}) error
}
