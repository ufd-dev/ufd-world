package watermarker

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/ufd-dev/ufd-world/internal/watermarker/internal/videoproc"
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

func ProcFile(collection string, in multipart.File, tags string) (uuid.UUID, error) {
	ct, err := detectContentType(in)
	if err != nil {
		return uuid.UUID{}, err
	}

	srcPath, err := saveTmpFile(collection, in, ct)
	if err != nil {
		return uuid.UUID{}, err
	}

	w, h, err := getResolution(ct, srcPath)

	id := uuid.New()

	overlayPath, err := createOverlay(id, w, h, tags)
	if err != nil {
		return uuid.UUID{}, err
	}

	processor := videoproc.GetSingleton()
	processor.AddJob(srcPath, overlayPath, "out1.mp4")
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
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

func saveTmpFile(collection string, in multipart.File, ct ContentType) (string, error) {
	var ext string
	switch ct {
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	case "video/mp4":
		ext = ".mp4"
	}

	path := os.TempDir() + string(os.PathSeparator) + collection
	tempFile, err := os.CreateTemp(path, "upload-"+"*"+ext)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	if _, err = io.Copy(tempFile, in); err != nil {
		return "", err
	}

	return filepath.Base(tempFile.Name()), nil
}
