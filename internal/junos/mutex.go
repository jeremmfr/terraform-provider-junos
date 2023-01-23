package junos

import "sync"

var Mutex = &sync.Mutex{} //nolint: gochecknoglobals
