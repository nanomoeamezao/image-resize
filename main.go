package main

import (
	"log"
	"net/http"

	"test/test/download"
	"test/test/flags"
	"test/test/request"
	"test/test/resizer"
)

func main() {
	log.Println("listening")
	limit := flags.ReadFlags()
	log.Println("limit: ", limit)
	http.HandleFunc("/resize", request.LimitMaxRequests(handleResize, limit))
	err := http.ListenAndServe("0.0.0.0:3300", nil)
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
	log.Println("got request")

	urlVal, height, width := request.GetQueryParams(r.URL)
	if urlVal == "" || height == "" || width == "" {
		request.HandleError("missing required params", w, http.StatusBadRequest)
		return
	}
	Img := ImageParams{}
	url, heightInt, widthInt, err := request.ParseParams(urlVal, height, width)
	if err != nil {
		request.HandleError("could not parse params: "+err.Error(), w, http.StatusBadRequest)
		return
	}
	Img.Width = widthInt
	Img.Height = heightInt
	Img.URL = url

	targetImg, err := download.DownloadImage(Img.URL)
	if err != nil {
		request.HandleError("could not get image: "+err.Error(), w, http.StatusInternalServerError)
		return
	}
	resImg, format, err := resizer.ResizeImage(targetImg, Img.Width, Img.Height)
	if err != nil {
		request.HandleError("could not resize image: "+err.Error(), w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/"+format)
	_, err = w.Write(resImg)
	if err != nil {
		log.Println("could not write response: ", err.Error())
		panic(err)
	}
}
