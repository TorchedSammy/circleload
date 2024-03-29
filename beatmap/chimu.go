package beatmap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	reqUrl, _ := url.Parse(fmt.Sprintf("https://api.chimu.moe/v1/search?query=%s&amount=%d", query, c.Options.MaxResults))
	// all gamemodes, golang ""enums"" am i right
	// basically we should only add the mode query if we dont want all gamemodes
	if c.Options.Mode != ModeAny {
		q := reqUrl.Query()
		q.Add("mode", fmt.Sprintf("%d", c.Options.Mode))
		reqUrl.RawQuery = q.Encode()
	}
	if c.Options.Status != StatusAny {
		q := reqUrl.Query()
		q.Add("status", fmt.Sprintf("%d", c.Options.Status))
		reqUrl.RawQuery = q.Encode()
	}

	resp, err := http.Get(reqUrl.String())
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

func (c *Chimu) SetMode(mode Mode) {
	c.Options.Mode = mode
}

func (c *Chimu) SetStatus(status Status) {
	c.Options.Status = status
}
