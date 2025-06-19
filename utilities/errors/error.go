package errors

type CustomError struct {
	Code          uint
	Message       string
	MessageToSend string
}

func (e *CustomError) Error() string {
	return e.Message
}
func (e *CustomError) Compile() {
	if e.MessageToSend == "" {
		e.MessageToSend = e.Message
	}
}
