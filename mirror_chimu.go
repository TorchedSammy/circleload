package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type chimuMirror struct{
	opts mirrorOptions
}

type chimuResponse struct {
	Code int
	Message string
	Data osuMapset
}

type chimuSearchResponse struct {
	chimuResponse
	Data []osuMapset
}

func (k chimuMirror) GetMapset(id int) (osuMapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/set/%d", id))
	if err != nil {
		return osuMapset{}, err
	}

	if resp.StatusCode == 404 {
		return osuMapset{}, ErrMapsetNotFound
	}

	var apiResp chimuResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	set := apiResp.Data

	return set, nil
}

func (k chimuMirror) Search(query string) ([]osuMapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/search?query=%s&amount=5", query))
	if err != nil {
		return nil, err
	}

	var apiResp chimuSearchResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	sets := apiResp.Data

	return sets, nil
}

func (k chimuMirror) GetMapsetData(id int) (*http.Response, error) {
	noVideo := 0
	if k.opts.noVideo {
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
