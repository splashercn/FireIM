package main

import (
	"github.com/splashercn/fireim/pkg/server"
)

func main() {
	s := server.NewServer()
	s.Run()
}
