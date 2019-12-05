package main

import (
	"log"
	"rank-task/routers"
)

func main() {
	router := routers.Load()

	err := router.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
