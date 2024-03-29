package apiLink

import (
	"io"
	"net/http"
)

const (
	baseMapUrl     string = "https://pokeapi.co/api/v2/location-area/"
	basePokemonUrl string = "https://pokeapi.co/api/v2/pokemon/"
)

func FetchMap(idx string) ([]byte, error) {
	res, err := doFetch(baseMapUrl + idx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FetchPokemon(name string) ([]byte, error) {
	res, err := doFetch(basePokemonUrl + name)
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
