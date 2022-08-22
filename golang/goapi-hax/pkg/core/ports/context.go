package ports

type Context interface {
	Bind(interface{}) error
	Query(string) (string, bool)
	Param(string) string
	// TransactionID() string
	Error(error)
	JSON(int, interface{})
}
