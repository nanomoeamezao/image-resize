package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/disintegration/imaging"
)

func main() {
	log.Println("listening")
	limit := readFlags()
	log.Println("limit: ", limit)
	http.HandleFunc("/resize", limitMaxRequests(handleResize, limit))
	http.ListenAndServe("0.0.0.0:3300", nil)
}

type limiter struct {
	current int
	max     int
	mu      *sync.Mutex
}

func NewLimit(limit int) limiter {
	mu := new(sync.Mutex)
	return limiter{
		current: 0,
		max:     limit,
		mu:      mu,
	}
}

func (l *limiter) Inc() {
	l.mu.Lock()
	l.current++
	l.mu.Unlock()
}

func (l *limiter) Dec() {
	l.mu.Lock()
	l.current--
	l.mu.Unlock()
}

func limitMaxRequests(f http.HandlerFunc, limit int) http.HandlerFunc {
	lim := NewLimit(limit)
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(lim.current, " ", lim.max)

		if lim.current >= lim.max {
			handleError("too many requests", w, http.StatusTooManyRequests)
			return
		}
		lim.Inc()
		f(w, r)
		lim.Dec()
	}
}

type ImageParams struct {
	URL    string
	Height int
	Width  int
}

func handleResize(w http.ResponseWriter, r *http.Request) {
	log.Println("got request")

	urlVal, height, width := getQueryParams(r.URL)
	if urlVal == "" || height == "" || width == "" {
		handleError("missing required params", w, http.StatusBadRequest)
		return
	}
	Img := ImageParams{}
	url, heightInt, widthInt, err := parseParams(urlVal, height, width)
	if err != nil {
		handleError("could not parse params: "+err.Error(), w, http.StatusBadRequest)
		return
	}
	Img.Width = widthInt
	Img.Height = heightInt
	Img.URL = url

	targetImg, err := downloadImage(Img.URL)
	if err != nil {
		handleError("could not get image: "+err.Error(), w, http.StatusInternalServerError)
		return
	}
	resImg, format, err := resizeImage(targetImg, Img.Width, Img.Height)
	if err != nil {
		handleError("could not resize image: "+err.Error(), w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/"+format)
	_, err = w.Write(resImg)
	if err != nil {
		log.Println("could not write response: ", err.Error())
		panic(err)
	}
}

func parseParams(urlVal string, heightStr string, widthStr string) (urlParsed string, height int, width int, err error) {
	decUrl, err := url.PathUnescape(urlVal)
	if err != nil {
		return "", 0, 0, errors.New("could not parse url " + err.Error())
	}
	heightInt, err := strconv.Atoi(heightStr)
	if err != nil {
		return "", 0, 0, errors.New("could not parse height " + err.Error())
	}
	widthInt, err := strconv.Atoi(widthStr)
	if err != nil {
		return "", 0, 0, errors.New("could not parse width " + err.Error())
	}
	return decUrl, heightInt, widthInt, nil
}

func resizeImage(targetImg []byte, width, height int) ([]byte, string, error) {
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

func downloadImage(url string) ([]byte, error) {
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

func prepareErrorForResponse(errStr string) string {
	return fmt.Sprintf(`{"error": "%s"}`, errStr)
}

func handleError(errStr string, w http.ResponseWriter, status int) {
	log.Println(errStr)
	sendErrResponse(errStr, status, w)
}

func sendErrResponse(errStr string, status int, w http.ResponseWriter) {
	jsonErr := prepareErrorForResponse(errStr)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(jsonErr)
}

func getQueryParams(url *url.URL) (urlVal string, height string, width string) {
	log.Println("parsing query")
	query := url.Query()
	urlVal = query.Get("url")
	height = query.Get("height")
	width = query.Get("width")
	return urlVal, height, width
}

func readFlags() int {
	// Read flags
	flagLimit := flag.Int("limit", 2, "amount of parallel requests allowed")
	flag.Parse()
	return *flagLimit
}
