package model

type Definition struct {
	Definition string `json:"definition" db:"definition"`
	ThumbsUp   int    `json:"thumbs_up" db:"thumbs_up"`
	Author     string `json:"author" db:"author"`
	Word       string `json:"word" db:"word"`
	WrittenOn  string `json:"written_on" db:"written_on"`
}
