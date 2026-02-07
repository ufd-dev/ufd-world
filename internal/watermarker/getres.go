package watermarker

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type ffProbeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func getGIFResolution(filePath string) (int, int, error) {
	return 0, 0, fmt.Errorf("GIF res detection not implemented")
}

func getWebPResolution(filePath string) (int, int, error) {
	return 0, 0, fmt.Errorf("WebP res detection not implemented")
}

func getMP4Resolution(filePath string) (int, int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v",
		"-show_entries", "stream=width,height", "-of", "json", filePath)

	stdout, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe command failed: %w", err)
	}

	var result ffProbeOutput
	if err := json.Unmarshal(stdout, &result); err != nil {
		return 0, 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	if len(result.Streams) > 0 {
		return result.Streams[0].Width, result.Streams[0].Height, nil
	}

	return 0, 0, fmt.Errorf("no video stream found in file")
}

func getResolution(ct ContentType, filePath string) (int, int, error) {
	switch ct {
	case ContentTypeGIF:
		return getGIFResolution(filePath)
	case ContentTypeWebP:
		return getWebPResolution(filePath)
	case ContentTypeMP4:
		return getMP4Resolution(filePath)
	default:
		return 0, 0, fmt.Errorf("Unsupported format: %v", ct)
	}
}
