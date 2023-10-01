package model

type YizkPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type YizkBlock struct {
	//Points           []YizkPoint `json:"points"`
	UpperLeftPoint   YizkPoint `json:"upper_left_point"`
	BottomRightPoint YizkPoint `json:"bottom_right"`
	Text             string    `json:"text"`
	OriginLanguage   string    `json:"origin_language"`
	TargetLanguage   string    `json:"target_language"`
	TranslatedText   string    `json:"translated_text"`
}

type YizkPage struct {
	Id       string      `json:"id"`
	Blocks   []YizkBlock `json:"blocks"`
	Filename string      `json:"filename"`
	Order    int         `json:"order"`
}
