package junos

import "sync"

var globalMutex = new(sync.Mutex) //nolint:gochecknoglobals

func MutexLock() {
	globalMutex.Lock()
}

func MutexUnlock() {
	globalMutex.Unlock()
}
