package mirror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Chimu struct {
	Options Options
}

type chimuResponse struct {
	Code int
	Message string
	Data Mapset
}

type chimuMap struct {
	BeatmapID int `json:"BeatmapId"`
	SetID int `json:"ParentSetId"`
}

type chimuMapResponse struct {
	Code int
	Message string
	Data chimuMap
}

type chimuSearchResponse struct {
	chimuResponse
	Data []Mapset
}

func (c Chimu) GetMapset(id int) (Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/set/%d", id))
	if err != nil {
		return Mapset{}, err
	}

	if resp.StatusCode == 404 {
		return Mapset{}, ErrMapsetNotFound
	}

	var apiResp chimuResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	set := apiResp.Data

	return set, nil
}

func (c Chimu) GetMapsetFromMap(id int) (Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/map/%d", id))
	if err != nil {
		return Mapset{}, err
	}

	if resp.StatusCode == 404 {
		return Mapset{}, ErrMapsetNotFound
	}

	var apiResp chimuMapResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	beatmap := apiResp.Data

	return c.GetMapset(beatmap.SetID)
}

func (c Chimu) Search(query string) ([]Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/search?query=%s&amount=%d", query, c.Options.MaxResults))
	if err != nil {
		return nil, err
	}

	var apiResp chimuSearchResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	sets := apiResp.Data

	return sets, nil
}

func (c Chimu) GetMapsetData(id int) (*http.Response, error) {
	noVideo := 0
	if c.Options.NoVideo {
		noVideo = 1
	}

	resp, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/download/%d?n=%d", id, noVideo))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrMapsetNotFound
	}

	return resp, nil
}
