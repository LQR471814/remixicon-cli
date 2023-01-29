package common

import (
	"encoding/gob"
	"os"
	"sync"
)

type Store[T any] struct {
	Data T

	lock sync.Mutex
	path string
}

func NewStore[T any](path string, initial T) *Store[T] {
	return &Store[T]{
		Data: initial,
		lock: sync.Mutex{},
		path: path,
	}
}

func (w *Store[T]) Load() error {
	defer w.lock.Unlock()
	w.lock.Lock()

	f, err := os.Open(w.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	value := new(T)
	decoder := gob.NewDecoder(f)
	err = decoder.Decode(value)
	if err != nil {
		return err
	}
	w.Data = *value
	return nil
}

func (w *Store[T]) Write() error {
	defer w.lock.Unlock()
	w.lock.Lock()

	f, err := os.Create(w.path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	return encoder.Encode(w.Data)
}
