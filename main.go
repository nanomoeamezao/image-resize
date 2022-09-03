package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/SKF/go-image-resizer/resizer"
)

func main() {
	http.HandleFunc("/resize", handleResize)
	log.Println("listening")
	err := http.ListenAndServe("localhost:3300", nil)
	if err != nil {
		panic(err)
	}
}

type ImageParams struct {
	URL    string
	Height int
	Width  int
}

func handleResize(w http.ResponseWriter, r *http.Request) {
	urlVal, height, width := getQueryParams(r.URL)
	Img := ImageParams{}
	decUrl, err := url.PathUnescape(urlVal)
	if err != nil {
		log.Println("could not parse url: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
	Img.URL = decUrl
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		log.Println("could not parse height: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		log.Println("could not parse width: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
	Img.Width = widthInt
	Img.Height = heightInt

	res, err := http.Get(Img.URL)
	if err != nil {
		log.Println("could not get image: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	// targetImg := make([]byte, res.ContentLength)
	// _, err = res.Body.Read(targetImg)
	// if err != nil {
	// 	log.Println("could not read image: ", err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// }

	tfile, _ := os.Create("test.jpg")
	io.Copy(tfile, res.Body)
	res.Body.Close()
	stat, _ := tfile.Stat()
	targetImg := make([]byte, stat.Size())
	tfile.Seek(0, 0)
	tfile.Read(targetImg)
	tfile.Close()
	os.Remove("test.jpg")
	resImg, err := resizer.ResizeImage(targetImg, resizer.JpegEncoder, Img.Width, Img.Height)
	if err != nil {
		log.Println("could not resize image: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "image/jpeg")
	_, err = w.Write(resImg)
	if err != nil {
		log.Println("could not write response: ", err.Error())
		panic(err)
	}
}

func getQueryParams(url *url.URL) (urlVal string, height string, width string) {
	query := url.Query()
	urlVal = query.Get("url")
	height = query.Get("height")
	width = query.Get("width")
	return urlVal, height, width
}
