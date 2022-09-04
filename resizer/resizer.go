package resizer

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/disintegration/imaging"
)

func ResizeImage(targetImg []byte, width, height int) ([]byte, string, error) {
	log.Println("resizing image")
	i, format, err := image.Decode(bytes.NewReader(targetImg))
	if err != nil {
		return nil, "", err
	}

	resizedImage := imaging.Resize(i, width, height, imaging.Lanczos)
	buf := new(bytes.Buffer)
	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, resizedImage, nil)
	case "png":
		err = png.Encode(buf, resizedImage)
	default:
		return nil, "", errors.New("bad image type: " + format)
	}
	return buf.Bytes(), format, nil
}
