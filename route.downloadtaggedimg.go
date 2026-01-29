package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/tidbyt/go-libwebp/webp"
)

func handleDownloadTaggedImg(w http.ResponseWriter, r *http.Request) {
	const maxUpload = 20 << 20 // 20 MiB since 2^20 = 1MiB
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)
	if err := r.ParseMultipartForm(maxUpload); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("File exceeds %v bytes", maxUpload), http.StatusRequestEntityTooLarge)
		return
	}

	tags := cleanTagInput(r.FormValue("tags"))
	if tags == "" {
		log.Println("no tags")
		http.Error(w, "No tags to add to the image", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		log.Println(err)
		http.Error(w, "Error reading uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	resultBuf, contentType, err := ProcessUpload(file, tags)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	replacer := strings.NewReplacer(" ", "-")
	filename := replacer.Replace("ufd "+tags) + contentTypeToExt(contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, filename))
	w.Header().Set("Content-Length", strconv.Itoa(resultBuf.Len()))

	if _, err = w.Write(resultBuf.Bytes()); err != nil {
		log.Println(err)
	}
}

func cleanTagInput(input string) string {
	replacer := strings.NewReplacer(",", " ", ";", " ")
	s := replacer.Replace(input)

	reg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	s = reg.ReplaceAllString(s, "")

	// handle multiple spaces and spaces on either edge
	fields := strings.Fields(s)
	return strings.ToLower(strings.Join(fields, " "))
}

// addLabelOverlay draws the text over a white box on the image
func addOverlay(src image.Image, text string) (image.Image, error) {
	const padding = 8.0
	const minFont float64 = 16
	const maxFont float64 = 32
	const topLine = "UnicornFartDust.com"

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	dc := gg.NewContext(width, height)
	dc.DrawImage(src, 0, 0)

	fontSize := min(max(float64(height)*.05, minFont), maxFont)
	if err := dc.LoadFontFace("arial.ttf", fontSize); err != nil {
		return nil, fmt.Errorf("could not load font: %v", err)
	}

	// draw the topLine text and white box behind it
	textW, textH := dc.MeasureString(topLine)
	boxW := textW + (padding * 2)
	boxH := textH + (padding * 2)
	dc.SetRGBA(0, 0, 0, .6)
	dc.DrawRectangle(0, 0, boxW, boxH)
	dc.Fill()
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(topLine, padding, padding, 0, 1)

	// draw the tag text and white box behind it
	text = "tags: ufd " + text
	textW, textH = dc.MeasureString(text)
	boxW = textW + (padding * 2)
	boxH = textH + (padding * 2)
	y := float64(height) - boxH
	dc.SetRGBA(0, 0, 0, .6)
	dc.DrawRectangle(0, y, boxW, boxH)
	dc.Fill()
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(text, padding, y+padding, 0, 1)

	return dc.Image(), nil
}

// ProcessUpload handles the detection and processing logic
func ProcessUpload(file multipart.File, tags string) (*bytes.Buffer, string, error) {
	// 1. Reset file pointer to start (just in case)
	if _, err := file.Seek(0, 0); err != nil {
		return nil, "", err
	}

	// 2. Peek at the format without decoding the whole image
	// image.DecodeConfig only reads the header
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, "", fmt.Errorf("unknown format: %v", err)
	}

	// 3. Rewind the file pointer so the actual decoders can read from the start
	if _, err := file.Seek(0, 0); err != nil {
		return nil, "", err
	}

	// 4. Route based on format
	outputBuffer := new(bytes.Buffer)

	switch format {
	case "gif":
		// Handle GIFs (potentially animated)
		err = processAnimatedGIF(file, outputBuffer, tags)
		// webp is not a mistake; animations are converted to this format
		return outputBuffer, "image/webp", err

	case "jpeg":
		// Handle Static JPEG
		err = processStaticImage(file, outputBuffer, func(w io.Writer, i image.Image) error {
			return jpeg.Encode(w, i, nil)
		}, tags)
		return outputBuffer, "image/jpeg", err

	case "png":
		// Handle Static PNG
		err = processStaticImage(file, outputBuffer, png.Encode, tags)
		return outputBuffer, "image/png", err

	default:
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}
}

// processStaticImage handles PNGs and JPEGs
func processStaticImage(r io.Reader, w io.Writer, encoder func(io.Writer, image.Image) error, tags string) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	finalImg, err := addOverlay(img, tags)
	if err != nil {
		return err
	}

	return encoder(w, finalImg)
}

// processAnimatedGIF handles the complex logic of frame iteration and quantization
func processAnimatedGIF(r io.Reader, w io.Writer, tags string) error {
	// Decode all frames
	g, err := gif.DecodeAll(r)
	if err != nil {
		return err
	}

	imageRect := image.Rectangle{
		Max: image.Point{X: g.Config.Width, Y: g.Config.Height},
	}

	var imgFrames []image.Image
	for i, frame := range g.Image {
		bounds := frame.Bounds()

		// 1. Draw frame onto RGBA canvas (handling disposal/transparency roughly)
		// Note: For perfect disposal handling, you need a virtual canvas that persists
		// across loops, but for simple overlays, drawing the current frame usually works.
		img := image.NewRGBA(imageRect)
		if i > 0 {
			draw.Draw(img, imageRect, imgFrames[i-1], image.Point{}, draw.Src)
		}
		draw.Draw(img, bounds, frame, bounds.Min, draw.Over)

		imgFrames = append(imgFrames, img)
	}

	// do this after to avoid adding the background's alpha channel across frames
	for i, frame := range imgFrames {
		processedImg, err := addOverlay(frame, tags)
		if err != nil {
			return err
		}
		imgFrames[i] = processedImg
	}

	// TODO: preserve variable delay
	frameDuration := time.Duration(g.Delay[0]*10) * time.Millisecond

	animatedWebPBytes, err := encodeAnimatedWebP(imgFrames, frameDuration)
	if err != nil {
		return err
	}

	if _, err = w.Write(animatedWebPBytes); err != nil {
		return err
	}

	return nil
}

func encodeAnimatedWebP(frames []image.Image, duration time.Duration) ([]byte, error) {
	if len(frames) == 0 {
		return nil, fmt.Errorf("no frames provided")
	}

	firstFrame := frames[0]
	bounds := firstFrame.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Initialize the encoder with width, height, and loop settings (kmin/kmax for keyframes)
	// Kmin and Kmax can be set to 0 to use library defaults.
	enc, err := webp.NewAnimationEncoder(width, height, 0, 0)
	if err != nil {
		return nil, err
	}
	defer enc.Close()

	// Add each frame with the specified duration
	for _, frame := range frames {
		if err := enc.AddFrame(frame, duration); err != nil {
			return nil, err
		}
	}

	// Assemble all frames into the final animated WebP byte slice
	return enc.Assemble()
}

func contentTypeToExt(ct string) string {
	exts, err := mime.ExtensionsByType(ct)
	if err != nil || len(exts) == 0 {
		return ""
	}
	return exts[0]
}
