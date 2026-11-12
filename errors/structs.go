package errors

import (
	"fmt"
)

type Err struct {
	Code    uint16 `json:"code"`
	Message string `json:"message"`
}

func (err *Err) Error() string {
	return fmt.Sprintf("[%d]: %s", err.Code, err.Message)
}

func CreateError(code uint16, message string) *Err {
	return &Err{Code: code, Message: message}
}
