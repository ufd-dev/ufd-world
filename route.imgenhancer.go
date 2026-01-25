package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	tags := r.FormValue("tags")
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

// addLabelOverlay draws the text over a white box on the image
func addLabelOverlay(src image.Image, text string) (image.Image, error) {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dc := gg.NewContext(width, height)
	dc.DrawImage(src, 0, 0)

	// Load a font - REQUIRED.
	// Ensure "arial.ttf" or a similar font file exists in your project root.
	// You can adjust the size (48) based on your needs.
	if err := dc.LoadFontFace("arial.ttf", 48); err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	// Calculate text width to size the white box
	textW, textH := dc.MeasureString(text)
	padding := 20.0
	boxW := textW + (padding * 2)
	boxH := textH + (padding * 2)

	// Position: Bottom center of the image
	x := float64(width) / 2
	y := float64(height) - 100 // 100px from bottom

	// Draw the White Box (centered on x, y)
	dc.SetRGB(1, 1, 1) // White
	dc.DrawRectangle(x-(boxW/2), y-(boxH/2), boxW, boxH)
	dc.Fill()

	// Draw the Text (Black)
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(text, x, y, 0.5, 0.35) // 0.5, 0.35 centers text visually

	return dc.Image(), nil
}

func renderImgEnhancer(w http.ResponseWriter, img *EnhancedImg) {
	if img != nil {
		log.Println("EnhancedImg: ", img.Filename, img.Error)
	}
	renderTemplate(w, "img-enhancer.tpl.html", nil)
}
