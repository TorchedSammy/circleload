package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type chimuMirror struct{}

type chimuResponse struct {
	Code int
	Message string
	Data osuMapset
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

func (k chimuMirror) GetMapsetData(id int, opts mirrorOptions) (io.ReadCloser, error) {
	noVideo := 0
	if opts.noVideo {
		noVideo = 1
	}

	resp, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/download/%d?n=%d", id, noVideo))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrMapsetNotFound
	}

	return resp.Body, nil
}
