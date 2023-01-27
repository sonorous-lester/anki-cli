package oxford

type Response struct {
	Metadata struct {
		Operation string `json:"operation"`
		Provider  string `json:"provider"`
		Schema    string `json:"schema"`
	} `json:"metadata"`
	Query   string `json:"query"`
	Results []struct {
		Id             string `json:"id"`
		Language       string `json:"language"`
		LexicalEntries []struct {
			Entries []struct {
				Pronunciations []struct {
					AudioFile        string   `json:"audioFile"`
					Dialects         []string `json:"dialects"`
					PhoneticNotation string   `json:"phoneticNotation"`
					PhoneticSpelling string   `json:"phoneticSpelling"`
				} `json:"pronunciations"`
				Senses []struct {
					Definitions []string `json:"definitions"`
					Examples    []struct {
						Text string `json:"text"`
					} `json:"examples"`
					Id        string `json:"id"`
					Subsenses []struct {
						Definitions []string `json:"definitions"`
						Examples    []struct {
							Text string `json:"text"`
						} `json:"examples"`
						Id string `json:"id"`
					} `json:"subsenses,omitempty"`
				} `json:"senses"`
			} `json:"entries"`
			Language        string `json:"language"`
			LexicalCategory struct {
				Id   string `json:"id"`
				Text string `json:"text"`
			} `json:"lexicalCategory"`
			Text string `json:"text"`
		} `json:"lexicalEntries"`
		Type string `json:"type"`
		Word string `json:"word"`
	} `json:"results"`
}
