package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "Anki-test",
		Usage: "Make creating Anki cards more easier",
		Action: func(c *cli.Context) error {
			fmt.Printf("Hello %q", c.Args().Get(0))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
