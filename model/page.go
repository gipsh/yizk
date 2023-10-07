package model

import (
	"encoding/json"
	"os"
)

type YizkPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type YizkBlock struct {
	//Points           []YizkPoint `json:"points"`
	UpperLeftPoint   YizkPoint `json:"upper_left_point"`
	BottomRightPoint YizkPoint `json:"bottom_right_point"`
	Text             string    `json:"text"`
	OriginLanguage   string    `json:"origin_language"`
	TargetLanguage   string    `json:"target_language"`
	TranslatedText   string    `json:"translated_text"`
}

type YizkPage struct {
	Id         string      `json:"id"` // page id in the original document
	Blocks     []YizkBlock `json:"blocks"`
	Filename   string      `json:"filename"`
	Order      int         `json:"order"`       // relative order in the document
	PageNumber string      `json:"page_number"` // page number in the original document
}

func ReadMetadata(filename string) (*YizkPage, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var page YizkPage
	err = json.Unmarshal([]byte(file), &page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}
