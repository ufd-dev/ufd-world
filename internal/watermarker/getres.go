package watermarker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type ffProbeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
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

get
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run yourfile.go <mp4_file_path>")
		os.Exit(1)
	}
	filePath := os.Args[1]

	width, height, err := getMP4Resolution(filePath)
	if err != nil {
		log.Fatalf("Error getting resolution: %v", err)
	}

	fmt.Printf("Video Resolution: %d x %d pixels\n", width, height)
}