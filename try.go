//
// This file simulate try-catch structure.
//
package gox

// Try simulate try catch
func Try(f func(), catcher func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			catcher(err)
		}
	}()
	f()
}
