package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"test/test/limiter"
)

func LimitMaxRequests(f http.HandlerFunc, limit int) http.HandlerFunc {
	lim := limiter.NewLimit(limit)
	return func(w http.ResponseWriter, r *http.Request) {
		if lim.Current >= lim.Max {
			HandleError("too many requests", w, http.StatusTooManyRequests)
			return
		}
		lim.Inc()
		f(w, r)
		lim.Dec()
	}
}

func ParseParams(urlVal string, heightStr string, widthStr string) (urlParsed string, height int, width int, err error) {
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

func prepareErrorForResponse(errStr string) string {
	return fmt.Sprintf(`{"error": "%s"}`, errStr)
}

func HandleError(errStr string, w http.ResponseWriter, status int) {
	log.Println(errStr)
	sendErrResponse(errStr, status, w)
}

func sendErrResponse(errStr string, status int, w http.ResponseWriter) {
	jsonErr := prepareErrorForResponse(errStr)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(jsonErr)
}

func GetQueryParams(url *url.URL) (urlVal string, height string, width string) {
	log.Println("parsing query")
	query := url.Query()
	urlVal = query.Get("url")
	height = query.Get("height")
	width = query.Get("width")
	return urlVal, height, width
}
