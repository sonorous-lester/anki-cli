package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	app := &cli.App{
		Name:  "ankictl",
		Usage: "Make creating Anki cards more easier",
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "add a new card info to file",
				Action: func(c *cli.Context) error {
					fileExisting := checkFile()
					if !fileExisting {
						createNewFile()
					}
					fmt.Printf("Create a new \"%s\" info in to file.\n", c.Args().First())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func checkFile() bool {
	file := "20220127.txt"
	path := os.Getenv("HOME") + "/Desktop"
	fullPath := filepath.Join(path, file)
	fmt.Printf("fullpath is %s\n", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("file %s does not exist in %s\n", file, path)
		return false
	} else {
		fmt.Printf("file %s exists in %s\n", file, path)
		return true
	}
}

func createNewFile() {
	desktop := os.Getenv("HOME") + "/Desktop/"
	currentTime := time.Now()
	fileName := currentTime.Format("2006-01-02") + ".txt"
	fullPath := filepath.Join(desktop, fileName)

	f, err := os.Create(fullPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	fmt.Println("File created: ", fullPath)
}

// Get the word
// Done!
// Check file is existing if not create a new file.
// Done!
// Send request to oxford
// Download audio to specific folder
// Mapping the response to anki struct
// Writing anki struct to file
// Response Message
