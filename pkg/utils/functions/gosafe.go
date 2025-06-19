package transactions

import "log"

// GoSafe go safe
func GoSafe(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("panic: ", r)
			}
		}()
		fn()
	}()
}
