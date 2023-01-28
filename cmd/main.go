package main

import (
	"anki-cli/anki"
	"anki-cli/oxford"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
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
					resp, err := queryWordToOxford(appID, appKey, q)
					if err != nil {
						return err
					}
					cards := mappingToCard(resp)
					wg := new(sync.WaitGroup)
					wg.Add(len(cards) * 2)
					for _, c := range cards {
						go downloadAudio(c.SoundAddr, ankiMedia, c.SoundName, wg)
						go writeToFile(ankiFile, c.AnkiString(), wg)
					}
					wg.Wait()
					fmt.Printf("Create a new \"%s\" info into file.\n", q)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func queryWordToOxford(id, key, word string) (oxford.Response, error) {
	req, _ := http.NewRequest("GET", "https://od-api.oxforddictionaries.com/api/v2/words/en-gb", nil)
	q := req.URL.Query()
	q.Add("q", word)
	q.Add("fields", "definitions,examples,pronunciations")
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app_id", id)
	req.Header.Add("app_key", key)

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	var rsp oxford.Response

	if err != nil {
		return rsp, err
	}

	if res.StatusCode != http.StatusOK {
		return rsp, errors.New("no matched word")
	}

	err = json.NewDecoder(res.Body).Decode(&rsp)
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}

func downloadAudio(u string, filepath, filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, err := url.Parse(u); err != nil {
		return
	}

	fullPath := fmt.Sprintf("%s/%s.mp3", filepath, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Create file error\n")
		return
	}
	defer file.Close()

	resp, err := http.Get(u)
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
	if len(resp.Results) == 0 {
		return out
	}
	for _, v := range resp.Results[0].LexicalEntries {
		if len(v.Entries) == 0 {
			continue
		}
		entry := v.Entries[0]

		ipa := "N/A"
		soundAddr := ""
		if len(entry.Pronunciations) != 0 {
			pronuc := entry.Pronunciations[0]
			ipa = pronuc.PhoneticSpelling
			soundAddr = pronuc.AudioFile
		}

		definition := "N/A"
		example := "N/A"

		if len(entry.Senses) != 0 {
			sense := entry.Senses[0]

			if len(sense.Definitions) != 0 {
				definition = sense.Definitions[0]
			}

			if len(sense.Examples) != 0 {
				example = sense.Examples[0].Text
			}
		}

		c := anki.Card{
			Text:         v.Text,
			PartOfSpeech: v.LexicalCategory.Id,
			IPA:          ipa,
			Sound:        fmt.Sprintf("[sound:%s_%s.mp3]", v.Text, v.LexicalCategory.Id),
			SoundName:    fmt.Sprintf("%s_%s", v.Text, v.LexicalCategory.Id),
			SoundAddr:    soundAddr,
			Definition:   definition,
			Example:      example,
		}
		out = append(out, c)
	}
	return out
}

func writeToFile(filePath, s string, wg *sync.WaitGroup) {
	defer wg.Done()
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
