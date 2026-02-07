package watermarker

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type ContentType string

const (
	ContentTypeGIF  ContentType = "image/gif"
	ContentTypeWebP ContentType = "image/webp"
	ContentTypeMP4  ContentType = "video/mp4"
)

type File struct {
	inPath  string
	inType  ContentType
	outPath string
	tags    string
}

func NewFile(in multipart.File, tags string) (*File, error) {
	ct, err := detectContentType(in)
	if err != nil {
		return nil, err
	}


	return &File{inType: ct, tags: tags}, nil
}

func detectContentType(f multipart.File) (ContentType, error) {
	buffer := make([]byte, 512)
	if _, err := f.Read(buffer); err != nil {
		return "", err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return "", err
	}
	contentType := ContentType(http.DetectContentType(buffer))

	switch contentType {
	case ContentTypeGIF:
		fallthrough
	case ContentTypeWebP:
		fallthrough
	case ContentTypeMP4:
		return contentType, nil
	default:
		return "", fmt.Errorf("Unsupported format: %v", contentType)
	}
}

func saveTmpFile(in multipart.File, ct ContentType) (string, error) {
	var ext string
	switch ct {
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	case "video/mp4":
		ext = ".mp4"
	}

	tempFile, err := os.CreateTemp("", "upload-"+"*"+ext)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	if _, err = io.Copy(tempFile, in); err != nil {
		return "", err
	}

	return tempFile.
}
