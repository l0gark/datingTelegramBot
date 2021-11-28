package main

import (
	"github.com/xlab/closer"
	"log"
)

type application struct {
	config *config
}

func main() {
	_, cleanup, err := initApp()
	if err != nil {
		log.Fatal("could not init application", err)
	}
	closer.Bind(func() {
		log.Print("stopping server")
		cleanup()
	})

}