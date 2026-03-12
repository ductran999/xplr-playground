package main

import (
	"log"
	"play-ground/software_acrh/master_worker/internal/master"
)

func main() {
	if err := master.Run(); err != nil {
		log.Fatalln("Server crashed:", err)
	}
}
