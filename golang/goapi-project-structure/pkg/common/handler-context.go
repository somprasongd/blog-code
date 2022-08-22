package common

type HandleFunc func(ctx HContext) error

type HContext interface {
	BodyParser(interface{}) error
	QueryParser(interface{}) error
	Query(string) (string, bool)
	DefaultQuery(string, string) string
	Param(string) string
	Header(string) string
	Authorization() string
	RequestId() string
	SendStatus(int) error
	SendJSON(int, interface{}) error
}
