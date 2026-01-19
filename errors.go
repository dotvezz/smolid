package smolid

import (
	"fmt"
)

var ErrUntyped = fmt.Errorf("no type id is embedded")
var ErrInvalidType = fmt.Errorf("invalid type id: must be less than or equal to " + fmt.Sprint(v1TypeSize))
