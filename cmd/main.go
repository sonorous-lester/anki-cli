package main

import (
	"anki-cli/anki"
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
	ankiFile := os.Getenv("ANKI_FILE")
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
					q := c.Args().First()
					resp := queryWordToOxford(appID, appKey, q)
					cards := mappingToCard(resp)
					for _, c := range cards {
						downloadAudio(c.SoundAddr, ankiMedia, c.SoundName)
						writeToFile(ankiFile, c.AnkiString())
					}
					fmt.Printf("Create a new \"%s\" info in to file.\n", q)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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

func mappingToCard(resp oxford.Response) []anki.Card {
	var out []anki.Card
	for _, v := range resp.Results[0].LexicalEntries {
		entry := v.Entries[0]
		pronuc := entry.Pronunciations[0]
		sense := entry.Senses[0]
		c := anki.Card{
			Text:         v.Text,
			PartOfSpeech: v.LexicalCategory.Id,
			IPA:          pronuc.PhoneticSpelling,
			Sound:        fmt.Sprintf("[sound:%s_%s.mp3]", v.Text, v.LexicalCategory.Id),
			SoundName:    fmt.Sprintf("%s_%s", v.Text, v.LexicalCategory.Id),
			SoundAddr:    pronuc.AudioFile,
			Definition:   sense.Definitions[0],
			Example:      sense.Examples[0].Text,
		}
		out = append(out, c)
	}
	return out
}

func writeToFile(filePath, s string) {
	fileName := time.Now().Format("2006-01-02") + ".txt"
	fullPath := filepath.Join(filePath, fileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(s + "\n")
	if err != nil {
		panic(err)
	}
}
