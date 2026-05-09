package imageutil

import (
	"io"
	"os"

	"github.com/disintegration/imaging"
)

const (
	MaxWidth  = 1920
	MaxHeight = 1920
	// JPEG quality: 75 keeps phone photos visually clean while shaving ~25% off
	// file size compared to 80 — invisible loss, easier on storage and bandwidth.
	Quality = 75
)

// Process reads an image from src, applies EXIF orientation, strips other EXIF
// metadata by re-encoding, resizes if larger than MaxWidth/MaxHeight, compresses
// as JPEG, and writes to dstPath.
func Process(src io.Reader, dstPath string) error {
	// imaging.Decode honours the EXIF Orientation tag so portrait shots from
	// phone cameras stay portrait after re-encode. Plain image.Decode ignores
	// EXIF and leaves images sideways.
	img, err := imaging.Decode(src, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w > MaxWidth || h > MaxHeight {
		img = imaging.Fit(img, MaxWidth, MaxHeight, imaging.Lanczos)
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return imaging.Encode(out, img, imaging.JPEG, imaging.JPEGQuality(Quality))
}
