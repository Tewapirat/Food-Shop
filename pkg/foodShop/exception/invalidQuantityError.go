package exception

import "fmt"

type InvalidQuantityError struct {
	Qty int
}

func (e InvalidQuantityError) Error() string {
	return fmt.Sprintf("Error: invalid quantity: %d (must be >= 1)", e.Qty)
}