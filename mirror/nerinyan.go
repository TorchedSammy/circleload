package mirror

import (
	"github.com/TorchedSammy/circleload/beatmap"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var nerinyanBaseUrl = "https://api.nerinyan.moe"
type Nerinyan struct {
	Options Options
}

func (c Nerinyan) GetMapset(id int) (beatmap.Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("%s/search?q=%d&option=setId", nerinyanBaseUrl, id))
	if err != nil {
		return beatmap.Mapset{}, err
	}

	if resp.StatusCode == 404 {
		return beatmap.Mapset{}, ErrMapsetNotFound
	}

	jsondata := []struct{
		SetID int `json:"id"`
		Artist string `json:"artist"`
		Title string `json:"title"`
		Creator string `json:"creator"`
		Status string `json:"status"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsondata)

	set := beatmap.Mapset{
		SetID: jsondata[0].SetID,
		Artist: jsondata[0].Artist,
		Title: jsondata[0].Title,
		Mapper: jsondata[0].Creator,
		Status: beatmap.StatusFromString(jsondata[0].Status),
	}

	return set, nil
}

func (c Nerinyan) GetMapsetFromMap(id int) (beatmap.Mapset, error) {
	resp, err := http.Get(fmt.Sprintf("%s/search?q=%d&option=mapId", nerinyanBaseUrl, id))
	if err != nil {
		return beatmap.Mapset{}, err
	}

	if resp.StatusCode == 404 {
		return beatmap.Mapset{}, ErrMapsetNotFound
	}

	jsondata := []struct{
		SetID int `json:"id"`
		Artist string `json:"artist"`
		Title string `json:"title"`
		Creator string `json:"creator"`
		Status string `json:"status"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsondata)

	set := beatmap.Mapset{
		SetID: jsondata[0].SetID,
		Artist: jsondata[0].Artist,
		Title: jsondata[0].Title,
		Mapper: jsondata[0].Creator,
		Status: beatmap.StatusFromString(jsondata[0].Status),
	}

	return set, nil
}

func (c Nerinyan) Search(query string, options Options) ([]beatmap.Mapset, error) {
	reqUrl, _ := url.Parse(fmt.Sprintf("%s/search?q=%s&option=title,artist,creator", nerinyanBaseUrl, query))
	// we should only add the mode query if we dont want all gamemodes
	if options.Mode != beatmap.ModeAny {
		q := reqUrl.Query()
		q.Add("m", fmt.Sprintf("%s", options.Mode))
		reqUrl.RawQuery = q.Encode()
	}

	if options.Status != beatmap.StatusAny {
		q := reqUrl.Query()
		q.Add("s", fmt.Sprintf("%s", options.Status))
		reqUrl.RawQuery = q.Encode()
	}

	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return nil, err
	}

	jsondata := []struct{
		SetID int `json:"id"`
		Artist string `json:"artist"`
		Title string `json:"title"`
		Creator string `json:"creator"`
		Status string `json:"status"`
	}{}
	json.NewDecoder(resp.Body).Decode(&jsondata)

	sets := []beatmap.Mapset{}
	for _, set := range jsondata {
		sets = append(sets, beatmap.Mapset{
			SetID: set.SetID,
			Artist: set.Artist,
			Title: set.Title,
			Mapper: set.Creator,
			Status: beatmap.StatusFromString(set.Status),
		})
	}

	return sets, nil
}

func (c Nerinyan) GetMapsetData(id int, options Options) (*http.Response, error) {
	noVideo := "false"
	if options.NoVideo {
		noVideo = "true"
	}

	resp, err := http.Get(fmt.Sprintf("%s/d/%d?noVideo=%s", nerinyanBaseUrl, id, noVideo))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrMapsetNotFound
	}

	return resp, nil
}

func (c Nerinyan) MirrorOptions() Options {
	return c.Options
}

func (c Nerinyan) SetOptions(options Options) {
	c.Options = options
}
