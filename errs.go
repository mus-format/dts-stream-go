package dtss

import "errors"

// ErrWrongDTM happens when DTS tries to unmarshal data with wrong DTM.
var ErrWrongDTM = errors.New("wrong data type metadata")
