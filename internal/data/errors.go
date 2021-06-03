package data

import "fmt"

type InvalidUser struct {
	Reason string
}

func (i *InvalidUser) Error() string {
	return fmt.Sprintf("Invalid User: %s", i.Reason)
}

type InvalidDBConn struct{Reason string}

func (i *InvalidDBConn) Error() string{
	return fmt.Sprintf("Invalid DB object: %s", i.Reason)
}
