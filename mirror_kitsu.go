package main

// This defines the kitsu.moe mirror.
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type kitsuMirror struct {
	opts mirrorOptions
}

// kitsu is NOT a proper restful api
// on success, it doesnt send code
// but on failure it does (and the actual response code is still 200)
// so to check for errors we just check is Code isnt 0
type kitsuResponse struct {
	osuMapset
	Code int
	Message string
}

func (k kitsuMirror) GetMapset(id int) (osuMapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/api/s/%d", id))
	if err != nil {
		return osuMapset{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return osuMapset{}, err
	}

	var apiResp kitsuResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return osuMapset{}, err
	}

	if apiResp.Code != 0 {
		switch apiResp.Code {
		case 404:
			return osuMapset{}, ErrMapsetNotFound
		}
	}

	set := osuMapset{
		SetID: apiResp.SetID,
		Title: apiResp.Title,
		Artist: apiResp.Artist,
	}

	return set, nil
}

// thanks copilot
func (k kitsuMirror) Search(query string) ([]osuMapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/api/search?query=%s&amount=5", query))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sets []osuMapset
	err = json.Unmarshal(body, &sets)
	if err != nil {
		return nil, err
	}

	return sets, nil
}

// get beatmap from kitsu
func (k kitsuMirror) GetMapsetData(id int) (*http.Response, error) {
	// kitsu doesnt have a noVideo option
	// log that it doesnt
	if k.opts.noVideo {
		warn("kitsu mirror doesnt support noVideo")
	}

	resp, err := http.Get(fmt.Sprintf("https://kitsu.moe/d/%d", id))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrMapsetNotFound
	}

	return resp, nil
}

