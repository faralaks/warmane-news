package db

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrParsingFaled = errors.New("parsing failed")
