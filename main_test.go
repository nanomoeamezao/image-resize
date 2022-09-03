package main

import (
	"bytes"
	"image"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetImage(t *testing.T) {
	img, err := downloadImage("https://www.thearmorylife.com/wp-content/uploads/2021/03/article-springfield-xd-4-service-model-9mm-review-1.jpg")
	imgBuff := bytes.NewBuffer(img)
	require.NoError(t, err)
	_, _, err = image.Decode(imgBuff)
	require.NoError(t, err)
}

func TestResize(t *testing.T) {
	img, err := os.Open("testimage.jpg")
	defer img.Close()

	var intImgBuff bytes.Buffer
	io.Copy(&intImgBuff, img)

	resized, _, err := resizeImage(intImgBuff.Bytes(), 100, 100)
	require.NoError(t, err)

	imgBuff := bytes.NewBuffer(resized)
	imageDec, _, err := image.Decode(imgBuff)
	require.NoError(t, err)

	size := imageDec.Bounds()
	require.Equal(t, 100, size.Dx())
	require.Equal(t, 100, size.Dy())
}
