package error

import "fmt"

type genericError struct {
	Msg string `json:"error"`
}

func (ge genericError) Error() string {
	return ge.Msg
}

// Errors definition

type KeyNotFound struct {
	genericError
}

type CouldNotCast struct {
	genericError
}

//Errors implementation

func NewStoreKeyNotFound(key string) KeyNotFound {
	return KeyNotFound{
		genericError: genericError{
			Msg: fmt.Sprintf("key <%s> not found", key),
		},
	}
}

func NewCouldNotCast(val interface{}) CouldNotCast {
	return CouldNotCast{
		genericError: genericError{
			Msg: fmt.Sprintf("could not cast value <%v> to string", val),
		},
	}
}
