package mirror

import (
	"github.com/TorchedSammy/circleload/beatmap"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Chimu struct {
	Options Options
}

func (c Chimu) GetMapset(id int) (beatmap.Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/set/%d", id))
	if err != nil {
		return beatmap.Mapset{}, err
	}

	if resp.StatusCode == 404 {
		return beatmap.Mapset{}, ErrMapsetNotFound
	}

	jsondata := struct{
		SetID int `json:"SetId"`
		Artist string `json:"Artist"`
		Title string `json:"Title"`
		Creator string `json:"Creator"`
		Status int `json:"RankedStatus"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsondata)

	set := beatmap.Mapset{
		SetID: jsondata.SetID,
		Artist: jsondata.Artist,
		Title: jsondata.Title,
		Mapper: jsondata.Creator,
		Status: beatmap.Status(jsondata.Status),
	}

	return set, nil
}

func (c Chimu) GetMapsetFromMap(id int) (beatmap.Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.chimu.moe/v1/map/%d", id))
	if err != nil {
		return beatmap.Mapset{}, err
	}

	if resp.StatusCode == 404 {
		return beatmap.Mapset{}, ErrMapsetNotFound
	}

	jsondata := struct{
		SetID int `json:"ParentSetId"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsondata)

	return c.GetMapset(jsondata.SetID)
}

func (c Chimu) Search(query string, options Options) ([]beatmap.Mapset, error) {
	reqUrl, _ := url.Parse(fmt.Sprintf("https://api.chimu.moe/v1/search?query=%s&amount=%d", query, options.MaxResults))
	// we should only add the mode query if we dont want all gamemodes
	if options.Mode != beatmap.ModeAny {
		q := reqUrl.Query()
		q.Add("mode", fmt.Sprintf("%d", options.Mode))
		reqUrl.RawQuery = q.Encode()
	}

	if options.Status != beatmap.StatusAny {
		q := reqUrl.Query()
		q.Add("status", fmt.Sprintf("%d", options.Status))
		reqUrl.RawQuery = q.Encode()
	}

	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return nil, err
	}

	jsondata := struct{
		Data []struct{
			SetID int `json:"SetId"`
			Artist string `json:"Artist"`
			Title string `json:"Title"`
			Creator string `json:"Creator"`
			Status int `json:"RankedStatus"`
		}
	}{}
	json.NewDecoder(resp.Body).Decode(&jsondata)

	sets := []beatmap.Mapset{}
	for _, set := range jsondata.Data {
		sets = append(sets, beatmap.Mapset{
			SetID: set.SetID,
			Artist: set.Artist,
			Title: set.Title,
			Mapper: set.Creator,
			Status: beatmap.Status(set.Status),
		})
	}

	return sets, nil
}

func (c Chimu) GetMapsetData(id int, options Options) (*http.Response, error) {
	noVideo := 0
	if options.NoVideo {
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

func (c Chimu) MirrorOptions() Options {
	return c.Options
}

func (c Chimu) SetOptions(options Options) {
	c.Options = options
}
