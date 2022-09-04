package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"

	"test/test/download"
	"test/test/request"
	"test/test/resizer"

	"github.com/stretchr/testify/require"
)

func TestGetImage(t *testing.T) {
	img, err := download.DownloadImage("https://www.thearmorylife.com/wp-content/uploads/2021/03/article-springfield-xd-4-service-model-9mm-review-1.jpg")
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

	resized, _, err := resizer.ResizeImage(intImgBuff.Bytes(), 100, 100)
	require.NoError(t, err)

	imgBuff := bytes.NewBuffer(resized)
	imageDec, _, err := image.Decode(imgBuff)
	require.NoError(t, err)

	size := imageDec.Bounds()
	require.Equal(t, 100, size.Dx())
	require.Equal(t, 100, size.Dy())
}

func TestServer(t *testing.T) {
	limit := 2
	http.HandleFunc("/resize", request.LimitMaxRequests(handleResize, limit))
	srv := http.Server{Addr: "0.0.0.0:3300"}
	go srv.ListenAndServe()
	responses := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func(responses chan int) {
		for v := range responses {
			fmt.Println("resp: ", v)
		}
		wg.Done()
	}(responses)
	go func(responses chan int) {
		for i := 0; i < 21; i++ {
			go func() {
				res, _ := http.Get("http://0.0.0.0:3300/resize")
				if res != nil {
					res.Body.Close()
					responses <- res.StatusCode
				}
				return
			}()
		}
		wg.Done()
		close(responses)
	}(responses)
	wg.Wait()
	srv.Shutdown(context.Background())
}
