package imageutil

import (
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/disintegration/imaging"
)

const (
	MaxWidth  = 1920
	MaxHeight = 1920
	Quality   = 80
)

// Process reads an image from src, strips EXIF metadata by decoding/re-encoding,
// resizes if larger than MaxWidth/MaxHeight, compresses as JPEG, and writes to dstPath.
// Returns the output file path (always .jpg).
func Process(src io.Reader, dstPath string) error {
	// Decode image (strips EXIF automatically since we re-encode)
	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Resize if needed (maintain aspect ratio, fit within MaxWidth x MaxHeight)
	if w > MaxWidth || h > MaxHeight {
		img = imaging.Fit(img, MaxWidth, MaxHeight, imaging.Lanczos)
	}

	// Write as JPEG with compression
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: Quality})
}
