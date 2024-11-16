package libraryDereferer

import (
	libraryLogging "github.com/s-r-engineer/library/logging"
	librarySync "github.com/s-r-engineer/library/sync"
)

const defaultStackSize = 100

func GetDefaultDereferer() (func(func()), func()) {
	return GetDereferer(defaultStackSize)
}

func GetDereferer(stackSize int) (func(func()), func()) {
	derefererStack := make(chan func(), stackSize) // Use a buffered channel.
	add, done, wait := librarySync.GetWait()
	lock, unlock := librarySync.GetMutex()
	once := librarySync.GetOnce()
	return func(f func()) {
			if f == nil {
				libraryLogging.Error("empty function")
				return
			}
			select {
			case derefererStack <- f:
			default:
				libraryLogging.Error("dereferer stack is full, function rejected")
			}
		}, func() {
			lock()
			defer unlock()
			once(func() {
				close(derefererStack)
			})
			for funcFromStack := range derefererStack {
				add()
				go func(f func()) {
					defer done()
					f()
				}(funcFromStack)
			}
			wait()
		}
}
