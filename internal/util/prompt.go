package util

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomFrom returns a random index from the provided list length.
func RandomFrom(list []string) (int, error) {
	if len(list) == 0 {
		return 0, errors.New("list is empty")
	}
	return rand.Intn(len(list)), nil
}

// RandomFrom returns a random index from the provided list length.
func RandomFromLength(length int) (int, error) {
	if length == 0 {
		return 0, errors.New("list is empty")
	}
	return rand.Intn(length), nil
}
