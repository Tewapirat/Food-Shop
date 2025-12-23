package exception

import "fmt"

type InvalidItemCodeError struct {
	Raw string
}

func (e *InvalidItemCodeError) Error() string {
	if e.Raw == "" {
		return "Error: invalid item code. Please provide a non-empty item code."
	}
	return fmt.Sprintf("Error: invalid item code %q. Please provide a valid item code.", e.Raw)
}
