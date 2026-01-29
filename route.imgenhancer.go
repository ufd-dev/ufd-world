package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

type EnhancedImg struct {
	Filename string
	Error    error
}

func handleImgEnhancer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		renderImgEnhancer(w, nil)
		return
	}

	const maxUpload = 20 << 20 // 20 MiB since 2^20 = 1MiB
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)
	if err := r.ParseMultipartForm(maxUpload); err != nil {
		renderImgEnhancer(w, &EnhancedImg{Error: fmt.Errorf("File exceeds %v bytes", maxUpload)})
		return
	}

	tags := cleanTagInput(r.FormValue("tags"))
	if tags == "" {
		renderImgEnhancer(w, &EnhancedImg{Error: fmt.Errorf("No tags to add to the image")})
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		log.Printf("FormFile error: %v", err)
		renderImgEnhancer(w, &EnhancedImg{Error: fmt.Errorf("Error reading uploaded file")})
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		renderImgEnhancer(w, &EnhancedImg{Error: fmt.Errorf("Invalid image format")})
		return
	}

	finalImage, err := addLabelOverlay(img, tags)
	if err != nil {
		log.Printf("Image processing error: %v", err)
		renderImgEnhancer(w, &EnhancedImg{Error: fmt.Errorf("Could not update image")})
		return
	}

	tmpFile, err := os.CreateTemp("", "processed-*.jpg")
	if err != nil {
		renderImgEnhancer(w, &EnhancedImg{Error: fmt.Errorf("Could not write updated image")})
		return
	}
	defer tmpFile.Close()

	if err := jpeg.Encode(tmpFile, finalImage, &jpeg.Options{Quality: 90}); err != nil {
		renderImgEnhancer(w, &EnhancedImg{Error: fmt.Errorf("Could not save updated image")})
		return
	}

	tmpPath := tmpFile.Name()
	go func(path string) {
		time.Sleep(2 * time.Minute)
		os.Remove(path)
	}(tmpPath)

	renderImgEnhancer(w, &EnhancedImg{Filename: filepath.Base(tmpPath)})
}

func cleanTagInput(input string) string {
	replacer := strings.NewReplacer(",", " ", ";", " ")
	s := replacer.Replace(input)

	reg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	s = reg.ReplaceAllString(s, "")

	// handle multiple spaces and spaces on either edge
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

// addLabelOverlay draws the text over a white box on the image
func addLabelOverlay(src image.Image, text string) (image.Image, error) {
	const padding = 10.0
	const fontSize float64 = 32
	const topLine = "UnicornFartDust.com"

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dc := gg.NewContext(width, height)
	dc.DrawImage(src, 0, 0)

	if err := dc.LoadFontFace("arial.ttf", fontSize); err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	// want the topLine width first to set a min width for the tags box
	topLineW, topLineH := dc.MeasureString(topLine)

	// draw the tag text and white box behind it
	text = "tags: ufd " + text
	textW, textH := dc.MeasureString(text)
	boxW := max(topLineW, textW) + (padding * 2)
	boxH := textH + (padding * 2)
	y := float64(height) - boxH
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(0, y, boxW, boxH)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(text, padding, y+padding, 0, 1)

	// draw the topLine text and white box behind it
	boxW = topLineW + (padding * 2)
	boxH = topLineH + (padding * 2)
	y = y - boxH
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(0, y, boxW, boxH)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(topLine, padding, y+padding, 0, 1)

	return dc.Image(), nil
}

func renderImgEnhancer(w http.ResponseWriter, img *EnhancedImg) {
	if img != nil {
		log.Println("EnhancedImg: ", img.Filename, img.Error)
	}
	renderTemplate(w, "img-enhancer.tpl.html", nil)
}
