package ipc

import "fmt"

type MonitoringError struct {
	err error
}

func (e *MonitoringError) Error() string {
	return fmt.Sprintf("subscription: %v", e.err)
}

func (e *MonitoringError) Unwrap() error {
	return e.err
}
