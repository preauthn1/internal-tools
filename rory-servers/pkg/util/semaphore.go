package util

// Semaphore ?
type Semaphore chan struct{}

// NewSemaphore ?
func NewSemaphore(n int) Semaphore {
	return make(Semaphore, n)
}

// Acquire ?
func (s Semaphore) Acquire() {
	s <- struct{}{}
}

// Release ?
func (s Semaphore) Release() {
	<-s
}
