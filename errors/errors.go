package errors

import "fmt"

type UnexpectedToken struct {
	Token string
}

func (e *UnexpectedToken) Error() string {
	return fmt.Sprintf("Unexpected token: %q", e.Token)
}
