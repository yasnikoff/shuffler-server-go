package md

import (
	"fmt"
)

var (
// derived Path errors
	ErrPathNotInitialized = fmt.Errorf("path is not initialized")
	ErrPathInitialized = fmt.Errorf("path is already initialized")
	ErrCorruptedMetadata = fmt.Errorf("path has corrupted metadata")
)
