package c_api

import (
	"runtime"
	"sync/atomic"
)

var mainThreadCh chan func() any
var resultCh chan any
var runtimeRunning = atomic.Bool{}
var asyncTasks = atomic.Uint32{}

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
	if mainThreadCh == nil {
		return
	}
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

	if r == nil {
		var t T
		return t
	}
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
