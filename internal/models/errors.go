package models

import "errors"

// errors.New equi to new Error() in JS
var ErrNoRecord = errors.New("models: no matching record found")
