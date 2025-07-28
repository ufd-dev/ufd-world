package media

import (
	"encoding/json"
	"io"
	"os"
)

type FileType string

const (
	FileTypeStillImg = "img"
	FileTypeGIF      = "anim"
	FileTypeVideo    = "video"
)

type Item struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Tags     string `json:"tags"`
}

func GetList() ([]Item, error) {
	file, err := os.Open("media/db.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var media []Item
	err = json.Unmarshal(jsonData, &media)
	if err != nil {
		return nil, err
	}

	return media, nil
}
