package data

import "fmt"

type InvalidUser struct {
	Reason string
}

func (i *InvalidUser) Error() string {
	return fmt.Sprintf("Invalid user: %s", i.Reason)
}

type InvalidDBConn struct{ Reason string }

func (i *InvalidDBConn) Error() string {
	return fmt.Sprintf("Invalid db object: %s", i.Reason)
}
