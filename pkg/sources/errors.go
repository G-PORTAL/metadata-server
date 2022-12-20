package sources

import (
	"errors"
)

var ErrFailedGetRemoteAddress = errors.New("failed to parse remote address")
var ErrNoMatchingMetadata = errors.New("no matching metadata found")
var ErrNoDatasourceFound = errors.New("no datasource found")
