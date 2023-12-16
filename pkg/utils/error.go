package utils

import "sync"

type errorTrap struct {
	err error
	mu  sync.RWMutex
}

// Create a new atomic error trap which can be used to capture the first error accross multiple goroutines.
func NewErrorTrap() *errorTrap {
	return &errorTrap{
		err: nil,
		mu:  sync.RWMutex{},
	}
}

// Set the error value. If an error has been already set the new value will be discarded.
func (et *errorTrap) Set(err error) {
	et.mu.Lock()
	defer et.mu.Unlock()

	if err != nil && et.err == nil {
		et.err = err
	}
}

// Check whether an error has been set.
func (et *errorTrap) IsSet() bool {
	et.mu.RLock()
	defer et.mu.RUnlock()

	return et.err != nil
}

// Access the error value
func (et *errorTrap) Err() error {
	et.mu.RLock()
	defer et.mu.RUnlock()

	return et.err
}
