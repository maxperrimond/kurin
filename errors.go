package kurin

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
		Obj    interface{}
		Errors *kensho.ValidationError
	}

	Unauthorized struct {
		Reason string
	}
)

func NewNotFound(id string, typeName string) *NotFound {
	return &NotFound{id, typeName}
}

func (err *NotFound) Error() string {
	return fmt.Sprintf(`not found "%s" with given id "%s"`, err.TypeName, err.ID)
}

func NewInvalid(obj interface{}, errors *kensho.ValidationError) *Invalid {
	return &Invalid{obj, errors}
}

func (err *Invalid) Error() string {
	return "data in object is invalid"
}

func NewUnauthorized(reason string) *Unauthorized {
	return &Unauthorized{reason}
}

func (err *Unauthorized) Error() string {
	return err.Reason
}
