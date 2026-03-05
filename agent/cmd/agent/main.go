package main

import (
	"agent/internal/snapshot"
	"fmt"
	"log"
)

func main() {
	total, cores, err := snapshot.ReadCPUSnapshot()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TOTAL:", total)
	fmt.Println("cores:", len(cores))
	if c4, ok := cores["cpu4"]; ok {
		fmt.Println("cpu4:", c4)
	}
}

// тест
/*
docker run --rm -it \
  -v "$PWD":/usr/src/app \
  -w /usr/src/app \
  golang:1.25.5 \
  go run ./cmd/agent
TOTAL: {3721 0 2048 611165 425 0 867 0 618226}
cores: 8
cpu4: {429 0 259 76532 47 0 25 0 77292}
*/
