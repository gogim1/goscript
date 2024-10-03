package file

import (
	"fmt"
)

type Error struct {
	Location SourceLocation
	Message  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[Error %s] %s", e.Location, e.Message)
}
