package anki

import "fmt"

type Card struct {
	Text         string
	PartOfSpeech string
	IPA          string
	Sound        string
	SoundAddr    string
	SoundName    string
	Definition   string
	Example      string
}

func (c Card) AnkiString() string {
	return fmt.Sprintf("%s;%s;%s;%s;%s;%s", c.Text, c.PartOfSpeech, c.IPA, c.Sound, c.Definition, c.Example)
}
