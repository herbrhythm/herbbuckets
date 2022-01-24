package bucket

import "errors"

var ErrNotFound = errors.New("file path not found")

var ErrExists = errors.New("file exists")
