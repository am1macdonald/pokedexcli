package apiLink

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const baseUrl string = "https://pokeapi.co/api/v2/location-area/"

func FetchMap(idx int) ([]byte, error) {
	fmt.Println("hmm... quite fetching...")
	res, err := doFetch(baseUrl + strconv.Itoa(idx))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func doFetch(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
