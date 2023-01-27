package main

import (
	"anki-cli/oxford"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	appID := os.Getenv("OXFORD_APP_ID")
	appKey := os.Getenv("OXFORD_APP_KEY")
	ankiMedia := os.Getenv("ANKI_MEDIA")
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
					resp := queryWordToOxford(appID, appKey, "swimming")
					downloadAudio(resp.Results[0].LexicalEntries[0].Entries[0].Pronunciations[0].AudioFile, ankiMedia, "swimming")
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
	file := "2023-01-27.txt"
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

func queryWordToOxford(id, key, word string) oxford.Response {
	req, _ := http.NewRequest("GET", "https://od-api.oxforddictionaries.com/api/v2/words/en-gb", nil)
	q := req.URL.Query()
	q.Add("q", word)
	q.Add("fields", "definitions,examples,pronunciations")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app_id", id)
	req.Header.Add("app_key", key)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var rsp oxford.Response
	_ = json.NewDecoder(res.Body).Decode(&rsp)
	fmt.Printf("rsp: %+v", rsp)
	return rsp
}

func downloadAudio(url string, filepath, filename string) {
	fullPath := fmt.Sprintf("%s/%s.mp3", filepath, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Create file error\n")
		return
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Download file error\n")
		return
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Download file error\n")
	}
}

// Get the word
// Done!
// Check file is existing if not create a new file.
// Done!
// Send request to oxford
// Done!
// Download audio to specific folder
// Done!
// Mapping the response to anki struct
// Writing anki struct to file
// Response Message
