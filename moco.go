//go:generate goversioninfo -64

package main

import (
	"C"
	"moco/config"
	"moco/server"
	"moco/utils/console"
	"moco/utils/logger"

	"os"
	"os/signal"
)
import "fmt"

func main() {
	console.InitWindows()
	log := logger.New("main")

	config, _, err := config.Load()

	if err != nil {
		// show error and ask the user to press any key to exit
		var input string

		// show the header without contextual data
		console.Header(nil)

		log.ErrorErr(err)
		fmt.Println("Press any key to exit...")
		fmt.Scanln(&input)
		os.Exit(1)
	}

	console.Header(config)

	// handle Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Info("bye 👋")
			os.Exit(0)
		}
	}()

	// Start the moco server
	err = server.Start(config.Port, config.Cors, config.Endpoints)

	if err != nil {
		log.PanicErr(err)
		os.Exit(1)
	}
}
