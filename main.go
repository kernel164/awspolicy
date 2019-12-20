package main

import (
	"github.com/kernel164/awspolicy/internal/engine"
)

func main() {
	err := engine.New().Run()
	if err != nil {
		panic(err)
	}
}
