package utils

import "sync"

type errorTrap struct {
	err error
	mu  sync.Mutex
}

func NewErrorTrap() *errorTrap {
	return &errorTrap{
		err: nil,
		mu:  sync.Mutex{},
	}
}

func (et *errorTrap) Set(err error) {
	et.mu.Lock()
	defer et.mu.Unlock()

	if err != nil && et.err == nil {
		et.err = err
	}
}

func (et *errorTrap) IsSet() bool {
	et.mu.Lock()
	defer et.mu.Unlock()

	return et.err != nil
}

func (et *errorTrap) Err() error {
	et.mu.Lock()
	defer et.mu.Unlock()

	return et.err
}
