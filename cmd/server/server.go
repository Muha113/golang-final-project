package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Muha113/golang-final-project/internal/app/server"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "cfgPath", "config/serverconfig.json", "cfg path")
	flag.Parse()

	config, err := server.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(config)
	fmt.Println("Starting server ...")
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
