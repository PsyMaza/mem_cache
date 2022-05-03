package main

import (
	"fmt"
	storage "memCache/storage"
)

func main() {
	s := storage.NewMemCache()

	s.Set("key", "fds")

	v, err := s.Get("key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(v)
}
