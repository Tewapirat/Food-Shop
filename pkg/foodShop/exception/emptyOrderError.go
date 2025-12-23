package exception

type EmptyOrderError struct{}

func (e *EmptyOrderError) Error() string {
	return "Error: empty order. Please add at least 1 item."
}

