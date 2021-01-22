package comm

import "fmt"

type BlockNotFoundError struct {
	Name string
}

func (e *BlockNotFoundError) Error() string {
	return fmt.Sprintf("Block with the given name not registered: '%s'", e.Name)
}

type BlockInitializationError struct {
	err  error
	Name string
}

func (e *BlockInitializationError) Error() string {
	return fmt.Sprintf("Error while initializing block: '%s'", e.Name)
}

func (e *BlockInitializationError) Unwrap() error {
	return e.err
}

type TagNotFoundError struct {
	Tag string
}

func (e *TagNotFoundError) Error() string {
	return fmt.Sprintf("The tag was not found: '%s'", e.Tag)
}

type TagCannotReceiveError struct {
	Tag string
}

func (e *TagCannotReceiveError) Error() string {
	return fmt.Sprintf("The tag is not a receiver: '%s'", e.Tag)
}

type TagReceiveError struct {
	err error
	Tag string
}

func (e *TagReceiveError) Error() string {
	return fmt.Sprintf("Tag returned an error: '%s'", e.Tag)
}

func (e *TagReceiveError) Unwrap() error {
	return e.err
}
