package event

import "fmt"

type AlreadyStartedError struct {
}

func (e *AlreadyStartedError) Error() string {
	return "Additional events cannot be subscribed to once the Subscription has started"
}

type MonitoringError struct {
	err error
}

func (e *MonitoringError) Error() string {
	return fmt.Sprintf("Error from subscription: %v", e.err)
}

func (e *MonitoringError) Unwrap() error {
	return e.err
}
