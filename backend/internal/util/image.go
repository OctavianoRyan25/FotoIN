package util

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png" // support png decoding
	"io"
	"strings"
)

const (
	// Target compression threshold: 3 MB
	CompressThresholdBytes = 3 * 1024 * 1024
	// Target max dimension
	MaxImageDimension = 3000
	// Target JPEG quality
	JpegQuality = 82
)

// CompressImage optimizes image bytes if larger than threshold or excessive resolution.
// Returns compressed bytes, updated content-type, and size.
func CompressImage(r io.Reader, filename string, originalSize int64, contentType string) (io.ReadSeeker, int64, string, error) {
	// If not an image or already under 3MB, skip compression
	isImage := strings.HasPrefix(contentType, "image/") ||
		strings.HasSuffix(strings.ToLower(filename), ".jpg") ||
		strings.HasSuffix(strings.ToLower(filename), ".jpeg") ||
		strings.HasSuffix(strings.ToLower(filename), ".png")

	if !isImage || originalSize <= CompressThresholdBytes {
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, 0, contentType, err
		}
		return bytes.NewReader(data), int64(len(data)), contentType, nil
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, 0, contentType, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		// If decoding fails, fall back to original
		return bytes.NewReader(data), int64(len(data)), contentType, nil
	}

	var buf bytes.Buffer
	opt := &jpeg.Options{Quality: JpegQuality}
	if err := jpeg.Encode(&buf, img, opt); err != nil {
		return bytes.NewReader(data), int64(len(data)), contentType, nil
	}

	compressedData := buf.Bytes()
	// Only use compressed if it actually saved space
	if len(compressedData) < len(data) {
		return bytes.NewReader(compressedData), int64(len(compressedData)), "image/jpeg", nil
	}

	return bytes.NewReader(data), int64(len(data)), contentType, nil
}
