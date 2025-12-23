package exception

import "fmt"


type UnknownMenuItemError struct {
	Code string
}

func (e UnknownMenuItemError) Error() string {
	return fmt.Sprintf("Error: unknown menu item code: %s", e.Code)
}
