package ipc

import "fmt"

type MonitoringError struct {
	err error
}

func (e *MonitoringError) Error() string {
	return fmt.Sprintf("Error from subscription: %v", e.err)
}

func (e *MonitoringError) Unwrap() error {
	return e.err
}
