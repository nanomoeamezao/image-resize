package download

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
)

func DownloadImage(url string) ([]byte, error) {
	log.Println("downloading image")
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New("could not get image: " + err.Error())
	}
	var targetImg bytes.Buffer
	io.Copy(&targetImg, res.Body)
	res.Body.Close()
	return targetImg.Bytes(), nil
}
