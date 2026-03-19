package c_api

import "runtime"

var mainThreadCh chan func() any
var resultCh chan any

func runMainThread() {
	mainThreadCh = make(chan func() any)
	resultCh = make(chan any)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		for task := range mainThreadCh {
			r := task()
			resultCh <- r
		}
	}()
}

func closeMainThread() {
	close(mainThreadCh)
	close(resultCh)
}

// SyncCall handles calls that return a value
func SyncCall[T any](task func() T) T {
	// Avoid using this function without calling CreateRuntime first
	if mainThreadCh == nil {
		return task()
	}

	mainThreadCh <- func() any {
		return task()
	}
	r := <-resultCh
	return r.(T)
}

// SyncCallVoid handles calls with no return value
func SyncCallVoid(task func()) {
	// Avoid using this function without calling CreateRuntime first
	if mainThreadCh == nil {
		task()
		return
	}
	mainThreadCh <- func() any {
		task()
		return nil
	}
	<-resultCh
}
