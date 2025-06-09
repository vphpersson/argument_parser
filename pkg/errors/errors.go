package errors

import (
	"errors"
)

var (
	ErrNilValue                    = errors.New("nil value")
	ErrMultipleOptionsWithSameName = errors.New("multiple options with same name")
	ErrNameNotFound                = errors.New("name not found")
	ErrUnsetOption = errors.New("unset option")
	ErrNoCurrentOption = errors.New("no current option")
	ErrNilOption = errors.New("nil option")
	ErrNonZeroArgCombinedOption = errors.New("non-zero arg combined option")
)
