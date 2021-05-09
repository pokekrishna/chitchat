package data

import "fmt"

type InvalidUser struct {
	Reason string
}

func (i *InvalidUser) Error() string {
	return fmt.Sprintf("Invalid User: %s", i.Reason)
}