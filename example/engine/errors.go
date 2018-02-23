package engine

import (
	"fmt"

	"github.com/maxperrimond/kensho"
)

type (
	NotFound struct {
		ID       string
		TypeName string
	}

	Invalid struct {
		Obj     interface{}
		Errors  *kensho.ValidationError
		Message string
	}
)

func NewNotFound(id string, typeName string) *NotFound {
	return &NotFound{id, typeName}
}

func (err *NotFound) Error() string {
	return fmt.Sprintf(`not found "%s" with given id "%s"`, err.TypeName, err.ID)
}

func NewInvalid(obj interface{}, errors *kensho.ValidationError) *Invalid {
	return &Invalid{obj, errors, ""}
}

func (err *Invalid) Error() string {
	if err.Message != "" {
		return err.Message
	}

	return "request is invalid"
}
