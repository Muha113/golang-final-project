package main

import (
	"fmt"
	"log"

	"github.com/Muha113/golang-final-project/internal/app/server"
)

func main() {
	//cfgPath := "golang-final-project/config/serverconfig.json"
	// config, err := server.NewConfig(cfgPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	config := &server.Config{
		Host: "localhost",
		Port: "8080",
	}

	srv := server.NewServer(config)
	fmt.Println("Starting server ...")
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
