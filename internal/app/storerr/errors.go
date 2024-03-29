package storerr

import "errors"

var ErrIdAlreadyExists = errors.New("item with such id already exists")
var ErrNoItemFound = errors.New("no such item found")
