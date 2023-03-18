package junos

import "sync"

var mutex = &sync.Mutex{} //nolint:gochecknoglobals

func MutexLock() {
	mutex.Lock()
}

func MutexUnlock() {
	mutex.Unlock()
}
