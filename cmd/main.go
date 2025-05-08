package main

import (
	"github.com/enigmasterr/final_project/internal/application"
)

func main() {
	app := application.New()
	//app.Run()
	app.RunServer()
}
