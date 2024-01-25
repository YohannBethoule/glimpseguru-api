package main

import (
	"fmt"
	"glimpseguru-api/router"
)

func main() {
	r := router.New()
	errRouter := r.Run()
	if errRouter != nil {
		panic(fmt.Sprintf("Unable to run router, %e", errRouter))
	}
}
